package mongo

/*
import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"common/data/dbbackends"
	"common/data/entities"
	"common/data/flatupdate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var peerIndexOnce sync.Once

type MongoPeer struct {
	GenericMongoBackend
	dbCollection *mongo.Collection
	resourceName string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewMongoPeer(dbc DBConnector) *MongoPeer {
	db := dbc.GetDB()
	peer := &MongoPeer{
		dbCollection: db.Collection("peers"),
		resourceName: "Peer",
	}
	peerIndexOnce.Do(func() { go CreatePeerIndexes(peer) })
	return peer

}

func CreatePeerIndexes(peer *MongoPeer) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	_, err := peer.dbCollection.Indexes().CreateMany(ctx,
		[]mongo.IndexModel{{Keys: bson.D{{"updated_at", -1}}},
			{Keys: bson.D{{"is_disabled", 1}, {"is_up", 1}, {"updated_at", 1},
				{"country_code", 1}}},
			{Keys: bson.D{{"is_disabled", 1}, {"is_up", 1}, {"updated_at", 1}}}})
	if err != nil {
		log.Printf("error creating indexes for peer doc, err: %s", err)
	}
}

func (p *MongoPeer) Create(e *entities.Peer) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	res, err := p.dbCollection.InsertOne(ctx, e)
	if err != nil {
		return e.ID, p.ResolveError(err, p.resourceName)
	}
	if res.InsertedID != e.ID {
		// TODO: Rollback?
		return e.ID, entities.ErrInvalidEntity
	}

	return e.ID, nil
}

func (p *MongoPeer) CreateOrUpdate(e *entities.Peer) (string, error) {
	// TODO: add config to common module (timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(true)
	_, err := p.dbCollection.UpdateOne(ctx, bson.D{{"_id", e.ID}}, bson.D{{"$set", e}}, opts)
	if err != nil {
		return e.ID, p.ResolveError(err, p.resourceName)
	}

	return e.ID, nil
}

func (p *MongoPeer) Get(id string) (*entities.Peer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var res entities.Peer
	err := p.dbCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&res)
	if err != nil {
		return nil, p.ResolveError(err, p.resourceName)
	}
	return &res, nil
}

func (p *MongoPeer) Update(e *entities.Peer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(false)
	updateMap, err := flatupdate.Flatten(e, false)
	res, err := p.dbCollection.UpdateOne(ctx, bson.D{{"_id", e.ID}}, bson.D{{"$set", updateMap}}, opts)
	if err != nil {
		return p.ResolveError(err, p.resourceName)
	}
	if res.ModifiedCount != 1 {
		return dbbackends.NewErrNotFound(p.resourceName)
	}
	return nil
}

func (p *MongoPeer) List() (*[]*entities.Peer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var res []*entities.Peer
	findOptions := options.Find()
	findOptions.SetLimit(100)
	cur, err := p.dbCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return &res, p.ResolveError(err, p.resourceName)
	}
	err = cur.All(ctx, &res)
	if err != nil {
		return &res, p.ResolveError(err, p.resourceName)
	}
	return &res, nil
}

func (p *MongoPeer) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	res, err := p.dbCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return p.ResolveError(err, p.resourceName)
	}
	if res.DeletedCount != 1 {
		return dbbackends.NewErrNotFound(p.resourceName)
	}
	return nil
}

func (p *MongoPeer) Search(query *entities.Peer, limit int) (*[]*entities.Peer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	var res []*entities.Peer
	var res1 []*entities.Peer
	var res2 []*entities.Peer
	heartBeatWindow := 3600
	queryMap, err := flatupdate.Flatten(query, true)
	if err != nil {
		return nil, err
	}
	randomWindowSecond := time.Duration(rand.Intn(heartBeatWindow))
	randomWindowSecondTime := time.Now().Add(time.Second * -randomWindowSecond)
	queryMap["updated_at"] = map[string]time.Time{"$gte": randomWindowSecondTime}
	opts := options.Find()
	opts.SetSort(bson.D{{"updated_at", 1}})
	opts.SetLimit(int64(limit))
	cur1, err1 := p.dbCollection.Find(ctx, queryMap, opts)
	opts.SetSort(bson.D{{"updated_at", -1}})
	queryMap["updated_at"] = map[string]time.Time{"$lt": randomWindowSecondTime,
		"$gte": time.Now().Add(time.Duration(-heartBeatWindow) * time.Second)}
	cur2, err2 := p.dbCollection.Find(ctx, queryMap, opts)
	// Both queries failed
	if err1 != nil && err2 != nil {
		return nil, p.ResolveError(err1, p.resourceName)
	}
	if err1 == nil {
		err1 = cur1.All(ctx, &res1)
	}
	if err2 == nil {
		err2 = cur2.All(ctx, &res2)
	}

	if err1 != nil && err2 != nil {
		// both cursor failed to load entities
		return &res1, p.ResolveError(err1, p.resourceName)
	}

	if err1 != nil {
		// the first cursor failed to load return the second result
		return &res2, nil
	}
	if err2 != nil {
		// the second cursor failed to load return the first result
		return &res1, nil
	}
	// TODO both cursor load successfully find the peers with updated_at closest to the random window timestamp from the result
	res = append(res1, res2...)
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return &res, nil
}

// TODO The following is temporary way to log peers activity

type PeerReport struct {
	Id struct {
		CountryCode string `bson:"country_code"`
		IsUp        bool   `bson:"is_up"`
		IsDisabled  bool   `bson:"is_disabled"`
	} `bson:"_id"`
	Sum int `bson:"count"`
}

func (p *MongoPeer) GetPeerReport() ([]*PeerReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	peersInfo := []*PeerReport{}
	specificCountryRes := []*PeerReport{}
	allCountriesRes := []*PeerReport{}

	specificCountry := bson.D{{"$group", bson.D{
		{"_id",
			bson.D{{"country_code", "$country_code"},
				{"is_up", "$is_up"},
				{"is_disabled", "$is_disabled"}}},
		{"count", bson.D{{"$sum", 1}}}}}}
	allCountries := bson.D{{"$group", bson.D{
		{"_id",
			bson.D{{"country_code", "ALL"},
				{"is_up", "$is_up"},
				{"is_disabled", "$is_disabled"}}},
		{"count", bson.D{{"$sum", 1}}}}}}
	specificCountryCur, specificCountryErr := p.dbCollection.Aggregate(ctx, mongo.Pipeline{specificCountry})
	if specificCountryErr != nil {
		return peersInfo, p.ResolveError(specificCountryErr, p.resourceName)
	}
	allCountriesCur, allCountriesErr := p.dbCollection.Aggregate(ctx, mongo.Pipeline{allCountries})
	if allCountriesErr != nil {
		return peersInfo, p.ResolveError(allCountriesErr, p.resourceName)
	}
	err := specificCountryCur.All(ctx, &specificCountryRes)
	if err != nil {
		return peersInfo, p.ResolveError(err, p.resourceName)
	}
	err = allCountriesCur.All(ctx, &allCountriesRes)
	if err != nil {
		return peersInfo, p.ResolveError(err, p.resourceName)
	}
	peersInfo = append(specificCountryRes, allCountriesRes...)
	return peersInfo, nil
}

func (p *MongoPeer) IncrementCounter(peerId string, success bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()
	counter := "success_requests"
	if !success {
		counter = "failed_requests"
	}
	res, err := p.dbCollection.UpdateOne(ctx, bson.D{{"_id", peerId}},
		bson.D{{"$inc", bson.D{{counter, 1}}}})
	if err != nil {
		return p.ResolveError(err, p.resourceName)
	}
	if res.ModifiedCount != 1 {
		return dbbackends.NewErrNotFound(p.resourceName)
	}
	return nil
}
*/
