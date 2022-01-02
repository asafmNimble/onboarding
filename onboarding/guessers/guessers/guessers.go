package guessers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"onboarding/common/grpc/api"
	"onboarding/common/grpc/guessers"
	"sync"
	"time"
)

var guessersMap = make(map[int64]map[string]int64)
var chanMap = make(map[int64]chan api.NumGuessResponse)
var IDs int64 = 1

type GuessServer struct {
	guessers.UnimplementedGuessersServer
}

func (*GuessServer) AddGuesser(_ context.Context, guesserRequest *guessers.AddGuesserRequest) (*guessers.AddGuesserResponse, error) {
	guesserID := IDs
	beginAt := guesserRequest.BeginAt
	incrementBy := guesserRequest.IncrementBy
	sleep := guesserRequest.Sleep
	IDs++
	m := make(map[string]int64)
	m["beginAt"] = beginAt
	m["incrementBy"] = incrementBy
	m["sleep"] = sleep
	guessersMap[guesserID] = m
	// start guessing
	inC := make(chan api.NumGuessResponse)
	chanMap[guesserID] = inC
	go newGuesser(guesserID, beginAt, incrementBy, sleep, outGuessC, inC)
	return &guessers.AddGuesserResponse{GuesserID: guesserID}, nil
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

func (*GuessServer) RemoveGuesser(_ context.Context, guesserRequest *guessers.RemoveGuesserRequest) (*guessers.RemoveGuesserResponse, error) {
	id := guesserRequest.GuesserID
	_, found := guessersMap[id]
	if !found {
		return nil, errors.New("guesser doesn't exist in database")
	}
	delete(guessersMap, id)
	return &guessers.RemoveGuesserResponse{
		Ok:        true,
		GuesserID: id,
	}, nil
}

func (*GuessServer) QueryGuesser(_ context.Context, guesserRequest *guessers.QueryGuesserRequest) (*guessers.QueryGuesserResponse, error) {
	id := guesserRequest.GuesserID
	_, found := guessersMap[id]
	if !found {
		return nil, errors.New("guesser doesn't exist in database")
	}
	return &guessers.QueryGuesserResponse{}, nil
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
	guessers.RegisterGuessersServer(s, &gs)
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
