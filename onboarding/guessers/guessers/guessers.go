package guessers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"onboarding/common/data/managers/guessers"
	"onboarding/common/grpc/api"
	guesserspb "onboarding/common/grpc/guessers"
	"sync"
	"time"
)

//var guessersMap = make(map[int64]map[string]int64)
var chanMap = make(map[int64]chan api.NumGuessResponse)
var IDs int64 = 1

type GuessServer struct {
	guesserspb.UnimplementedGuessersServer
	MongoManage guessers.Manager
}

func notFound(s string) bool {
	if s == "mongo: no documents in result" {
		return true
	}
	return false
}

func (gs *GuessServer) AddGuesser(_ context.Context, guesserRequest *guesserspb.AddGuesserRequest) (*guesserspb.AddGuesserResponse, error) {
	guesserID := IDs
	beginAt := guesserRequest.BeginAt
	incrementBy := guesserRequest.IncrementBy
	sleep := guesserRequest.Sleep
	IDs++
	_, err := gs.MongoManage.AddGuesser(guesserID, beginAt, incrementBy,sleep)
	/*
	m := make(map[string]int64)
	m["beginAt"] = beginAt
	m["incrementBy"] = incrementBy
	m["sleep"] = sleep
	guessersMap[guesserID] = m
	*/
	if err != nil {
		return nil, err
	}
	// start guessing
	inC := make(chan api.NumGuessResponse)
	chanMap[guesserID] = inC
	go newGuesser(guesserID, beginAt, incrementBy, sleep, outGuessC, inC)
	return &guesserspb.AddGuesserResponse{GuesserID: guesserID}, nil
}

func newGuesser(guesserID int64, beginAt int64, incrementBy int64, sleep int64, outC chan api.NumGuessRequest, inC chan api.NumGuessResponse) {
	numToGuess := beginAt
	for {
		// guess beginAt
		outC <- api.NumGuessRequest{Num: numToGuess, GuesserID: guesserID}
		resp := <-inC
		if !resp.Ok {
			fmt.Printf("Recieved the following error : %v", resp.Err)
		}
		if resp.Found {
			// TODO: Read in Onboarding what needs to be done - initiate find of closest bigger prime
		}
		// sleep sleep
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		// incrementBy
		numToGuess += incrementBy
	}

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

	//_, found := guessersMap[id]
	//if !found {
	//	return nil, errors.New("guesser doesn't exist in database")
	//}
	//delete(guessersMap, id)
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
	//_, found := guessersMap[id]
	//if !found {
	//	return nil, errors.New("guesser doesn't exist in database")
	//}
	var guesses []*guesserspb.Guess
	// TODO: add if list is empty
	for _, g := range guesser.GuessesMade {
		guesses = append(guesses, &guesserspb.Guess{
			Num: g.GuessNum,
			Time:    g.GuessedAt.Unix(),
		})
	}
	return &guesserspb.QueryGuesserResponse{
		GuesserID: id,
		GuessList: guesses,
		Active:    guesser.Active,
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
	conn, err := grpc.Dial(server, grpc.WithInsecure())
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
var stream api.GuessNums_GuessNumClient

func sendGuesses(outC chan api.NumGuessRequest) {
	for {
		req := <-outC
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Failed to receive a response : %v", err)
		}
	}
}

func receiveGuesses() {
	for {
		resp, err := stream.Recv()
		guesserId := resp.GuesserID
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
	s := grpc.NewServer()
	gs := GuessServer{}
	guesserspb.RegisterGuessersServer(s, &gs)
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	if err := s.Serve(lis); err != nil { // assigns lis' port to s
		log.Fatalf("Recieved the following error : %v", err)
	}
	apiClient, err := getApiClient()
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	stream, err = apiClient.GuessNum(context.Background())
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	go sendGuesses(outGuessC)
	go receiveGuesses()
	return 0
}
