package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"onboarding/common/data/dbbackends"
	"onboarding/common/data/flatupdate"
)

type GenericMongoBackend struct{}

/*
This function returns flat bson.E list for a safe mongo update function for fields that are arrays of embedded docs
fieldName is the embedded document array field name in the main document
entity is the entity model of the embedded document type (that the array consists of)
*/
func (*GenericMongoBackend) EmbeddedEntityToArrUpdateFields(fieldName string, entity interface{}, disallowedFields *[]string) (*[]bson.E, error) {
	updateMap, err := flatupdate.Flatten(entity, false)
	if err != nil {
		return nil, err
	}

	updateFields := []bson.E{}
	for k, v := range updateMap {
		err = VerifyAllowedFieldsOnly(k, disallowedFields)
		if err != nil {
			return nil, err
		}
		if k != "_id" {
			updateFields = append(updateFields, bson.E{fmt.Sprintf("%s.$.%s", fieldName, k), v})
		}
	}
	return &updateFields, nil
}

func VerifyAllowedFieldsOnly(updateField string, disallowedFields *[]string) error {
	// disallow updating unique fields with this function, can be incorporated in the flatbson functionality in future
	if disallowedFields != nil {
		for _, uniqueField := range *disallowedFields {
			if updateField == uniqueField && updateField != "_id" {
				return errors.New(fmt.Sprintf("cannot update user with unique field: %s", uniqueField))
			}
		}
	}
	return nil
}

func (*GenericMongoBackend) ArrEmbeddedEntityProjection(ctx context.Context, fieldName string,
	dbCollection *mongo.Collection, mainID primitive.ObjectID, embeddedID primitive.ObjectID) *mongo.SingleResult {
	projection := bson.D{
		{fieldName, bson.D{
			{"$elemMatch", bson.D{{"_id", embeddedID}}},
		}}}
	res := dbCollection.FindOne(ctx, bson.D{{"_id", mainID}},
		options.FindOne().SetProjection(projection))
	return res
}

func (*GenericMongoBackend) ArrEmbeddedEntityProjectionByName(ctx context.Context, fieldName string, dbCollection *mongo.Collection, mainName string, embeddedName string) *mongo.SingleResult {
	projection := bson.D{
		{fieldName, bson.D{
			{"$elemMatch", bson.D{{"name", embeddedName}}},
		}}}
	res := dbCollection.FindOne(ctx, bson.D{{"name", mainName}},
		options.FindOne().SetProjection(projection))
	return res
}

func (*GenericMongoBackend) RunInTransaction(ctx context.Context, sess mongo.Session, f func(sessCtx mongo.SessionContext) error) error {
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(false)
	sessCtx := mongo.NewSessionContext(ctx, sess)
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	err := sess.StartTransaction(txnOpts)
	if err != nil {
		return err
	}
	err = f(sessCtx)
	if err != nil {
		if abortErr := sess.AbortTransaction(sessCtx); abortErr != nil {
			log.Printf("error trying to abort failed transaction: %s", abortErr)
		}
		return err
	}
	err = sess.CommitTransaction(sessCtx)
	if err != nil {
		if abortErr := sess.AbortTransaction(sessCtx); abortErr != nil {
			log.Printf("error trying to abort failed transaction: %s", abortErr)
		}
	}
	return err
}

func (*GenericMongoBackend) ResolveError(mongoErr error, resourceName string) error {
	if errors.Is(mongoErr, mongo.ErrNoDocuments) {
		return dbbackends.NewErrNotFound(resourceName)
	}
	if mongo.IsDuplicateKeyError(mongoErr) {
		return dbbackends.NewErrAlreadyExists(resourceName)
	}

	return dbbackends.NewErrUnexpected()
}
