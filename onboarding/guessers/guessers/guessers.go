package guessers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"net"
	"onboarding/common/data/dbbackends/mongo"
	"onboarding/common/data/entities"
	"onboarding/common/data/managers/guessers"
	"onboarding/common/data/managers/numbers"
	"onboarding/common/grpc/api"
	guesserspb "onboarding/common/grpc/guessers"
	"strconv"
	"sync"
	"time"
)

var chanMap = make(map[int64]chan api.NumGuessResponse)
var killGuessers = make(map[int64]bool)
var IDs int64 = 1

type GuessServer struct {
	guesserspb.UnimplementedGuessersServer
	MongoManage guessers.Manager
}

// Local machinery
type PrimeDetails struct {
	guesserID   int64
	timeGuessed time.Time
	originNum   int64
}

var primesMap = make(map[int64][]PrimeDetails)

func findPrime(originNum int64) int64 {
	if originNum == 1 {
		return 2
	} else if originNum == 2 {
		return 3
	}
	base := 2
	checkIfPrime := originNum + 1
	if checkIfPrime%2 == 0 {
		checkIfPrime++
	}
	for {
		basePowCheckIfPrimes := math.Pow(float64(base), float64(checkIfPrime-1))
		fermatLittletheorm := int64(basePowCheckIfPrimes) % checkIfPrime
		if fermatLittletheorm == 1 {
			return checkIfPrime
		}
		checkIfPrime = checkIfPrime + 2
	}
}

func createPrimeDetails(guesserID int64, timeGuessed time.Time, originNum int64) PrimeDetails {
	return PrimeDetails{
		guesserID:   guesserID,
		timeGuessed: timeGuessed,
		originNum:   originNum,
	}
}

// end of local machinery

func (gs *GuessServer) AddGuesser(_ context.Context, guesserRequest *guesserspb.AddGuesserRequest) (*guesserspb.AddGuesserResponse, error) {
	guesserID := IDs
	beginAt := guesserRequest.BeginAt
	incrementBy := guesserRequest.IncrementBy
	sleep := guesserRequest.Sleep
	IDs++
	_, err := gs.MongoManage.AddGuesser(guesserID, beginAt, incrementBy, sleep)
	if err != nil {
		return nil, err
	}
	// start guessing
	inC := make(chan api.NumGuessResponse)
	chanMap[guesserID] = inC
	killGuessers[guesserID] = false
	go newGuesser(guesserID, beginAt, incrementBy, sleep, outGuessC, inC)
	return &guesserspb.AddGuesserResponse{GuesserID: guesserID}, nil
}

func newGuesser(guesserID int64, beginAt int64, incrementBy int64, sleep int64, outC chan api.NumGuessRequest, inC chan api.NumGuessResponse) {
	fmt.Println("started a new goroutine guesser")
	numToGuess := beginAt
	for !killGuessers[guesserID] {
		// guess beginAt
		fmt.Printf("Guessers Client, newGuesser, Guessing %v\n", numToGuess)
		outC <- api.NumGuessRequest{Num: numToGuess, GuesserID: guesserID}
		fmt.Printf("Guessers Client, newGuesser, Sent %v on channel\n", numToGuess)
		resp := <-inC
		fmt.Printf("Guessers Client, newGuesser, received %v on channel\n", resp.Num)
		if !resp.Ok {
			fmt.Printf("Recieved the following error : %v", resp.Err)
		}
		if resp.Found {
			primeNum := findPrime(resp.Num)
			timeFound := time.Now()
			currPrimeDetails := createPrimeDetails(resp.GuesserID, timeFound, resp.Num)
			currPrimeList := primesMap[primeNum]
			currPrimeList = append(currPrimeList, currPrimeDetails)
			primesMap[primeNum] = currPrimeList

			// Update Number guessers list
			guessForNumber := entities.GuessType{
				FoundBy: strconv.Itoa(int(resp.GuesserID)),
				FoundAt: timeFound,
			}
			mn := numbers.NewManager(mongo.NewMongoNumber(mongo.NewMongoConnector()))
			_, err := mn.UpdateGuessForNumber(numToGuess, &guessForNumber)
			if err != nil {
				log.Fatalf("GuesserID: %v encountered the following error : %v", guesserID, err)
			}

			correctGuess := entities.Guess{
				GuessNum:  numToGuess,
				GuessedAt: timeFound,
			}
			mg := guessers.NewManager(mongo.NewMongoGuesser(mongo.NewMongoConnector()))
			_, err = mg.UpdateGuessedNumForGuesser(guesserID, &correctGuess)
			if err != nil {
				log.Fatalf("GuesserID: %v encountered the following error : %v", guesserID, err)
			}

		}
		// sleep sleep
		time.Sleep(time.Duration(sleep) * time.Second)
		// incrementBy
		numToGuess += incrementBy
	}
	return
}

func (gs *GuessServer) RemoveGuesser(_ context.Context, guesserRequest *guesserspb.RemoveGuesserRequest) (*guesserspb.RemoveGuesserResponse, error) {
	id := guesserRequest.GuesserID
	_, err := gs.MongoManage.GetGuesser(id)
	if err != nil {
		return nil, err
	}
	_, err = gs.MongoManage.RemoveGuesser(id)
	if err != nil {
		return nil, err
	}
	killGuessers[id] = true
	return &guesserspb.RemoveGuesserResponse{
		Ok:        true,
		GuesserID: id,
	}, nil
}

func (gs *GuessServer) QueryGuesser(_ context.Context, guesserRequest *guesserspb.QueryGuesserRequest) (*guesserspb.QueryGuesserResponse, error) {
	id := guesserRequest.GuesserID
	guesser, err := gs.MongoManage.GetGuesser(id)
	if err != nil {
		return nil, err
	}
	var guesses []*guesserspb.Guess
	for _, g := range guesser.GuessesMade {
		guesses = append(guesses, &guesserspb.Guess{
			Num:  g.GuessNum,
			Time: g.GuessedAt.Unix(),
		})
	}
	return &guesserspb.QueryGuesserResponse{
		GuesserID: id,
		GuessList: guesses,
		Active:    guesser.Active,
	}, nil
}

func (gs *GuessServer) QueryPrimes(_ context.Context, primesRequest *guesserspb.QueryPrimesRequest) (*guesserspb.QueryPrimesResponse, error) {
	// TODO: move to mongo
	var primes []*guesserspb.Prime
	for p, _ := range primesMap {
		var primeList []*guesserspb.Primedets
		for _, dets := range primesMap[p] {
			primeList = append(primeList, &guesserspb.Primedets{
				GuesserID: dets.guesserID,
				Time:      dets.timeGuessed.Unix(),
				OriginNum: dets.originNum,
			})
		}
		primes = append(primes, &guesserspb.Prime{
			Prime:     p,
			PrimeList: primeList,
		})
	}
	return &guesserspb.QueryPrimesResponse{
		Ok:     true,
		Primes: primes,
	}, nil
}

// API Server
var apiServerAddr = "localhost:5000"
var apiClient api.GuessNumsClient

var apiOnce sync.Once

func getApiClientInternal(server string) (api.GuessNumsClient, error) {
	if apiClient != nil {
		return apiClient, nil
	}
	var conn *grpc.ClientConn
	var err error
	for i := 1; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if conn == nil {
			conn, err = grpc.Dial(server, grpc.WithInsecure())
		}
	}
	if err != nil {
		return nil, err
	}

	apiClient = api.NewGuessNumsClient(conn)
	return apiClient, nil
}

func getApiClient() (api.GuessNumsClient, error) {
	var err error
	apiOnce.Do(func() {
		apiClient, err = getApiClientInternal(apiServerAddr)
	})
	if err != nil {
		return nil, err
	}
	return apiClient, nil
}

var outGuessC = make(chan api.NumGuessRequest)

func sendGuesses(outC chan api.NumGuessRequest, stream api.GuessNums_GuessNumClient) {
	for {
		req := <-outC
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Failed to receive a response : %v", err)
		}
	}
}

func receiveGuesses(stream api.GuessNums_GuessNumClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		guesserId := resp.GuesserID
		fmt.Printf("Guessers Client, receiveGuesses, GuesserID is: %v\n", resp.GuesserID)
		inC := chanMap[guesserId]
		if err == io.EOF {
			// read done.
			close(inC)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a response : %v", err)
		}
		inC <- *resp
	}
}

func RealGuessers() int {
	go GuesserServer()
	go ManageGuesses()
	return 0
}

func ManageGuesses() {
	apiClient, err := getApiClient()
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	stream, err := apiClient.GuessNum(context.Background())
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	go sendGuesses(outGuessC, stream)
	go receiveGuesses(stream)
	return
}

func GuesserServer() {
	s := grpc.NewServer()
	mg := guessers.NewManager(mongo.NewMongoGuesser(mongo.NewMongoConnector()))
	gs := GuessServer{MongoManage: *mg}
	guesserspb.RegisterGuessersServer(s, &gs)
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	if err := s.Serve(lis); err != nil { // assigns lis' port to s
		log.Fatalf("Recieved the following error : %v", err)
	}
	return
}
