package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"onboarding/common/data/dbbackends/redis"
	"onboarding/common/data/entities"
	"onboarding/common/data/managers/guessers_counters"
	"onboarding/common/grpc/api"
	guesserspb "onboarding/common/grpc/guessers"
	numberspb "onboarding/common/grpc/numbers"
	"sync"
	"time"
)

// Numbers Server

var numServerAddr = "localhost:7000"
var numbersClient numberspb.NumbersClient
var numOnce sync.Once

func getNumbersClientInternal(server string) (numberspb.NumbersClient, error) {
	if numbersClient != nil {
		return numbersClient, nil
	}
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	numbersClient = numberspb.NewNumbersClient(conn)
	return numbersClient, nil
}

func getNumClient() (numberspb.NumbersClient, error) {
	var err error
	numOnce.Do(func() {
		numbersClient, err = getNumbersClientInternal(numServerAddr)
	})
	if err != nil {
		return nil, err
	}
	return numbersClient, nil
}

type addNumInput struct {
	Num int64 `json:"num"`
}

func addNum(c *gin.Context) {
	var numToAdd addNumInput
	if err := c.ShouldBindJSON(&numToAdd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	numClient, err := getNumClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := numClient.AddNum(context.Background(), &numberspb.AddNumRequest{Num: numToAdd.Num})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    "Accepted",
		"num_added": response.Num,
	})
}

type removeNumInput struct {
	Num int64 `json:"num"`
}

func removeNum(c *gin.Context) {
	var numToRemove removeNumInput
	if err := c.ShouldBindJSON(&numToRemove); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	numClient, err := getNumClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := numClient.RemoveNum(context.Background(), &numberspb.RemoveNumRequest{Num: numToRemove.Num})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "Accepted",
		"num_removed": response.Num,
	})
}

type queryNumInput struct {
	Num int64 `json:"num"`
}

func queryNumber(c *gin.Context) {
	var numToGet queryNumInput
	if err := c.ShouldBindJSON(&numToGet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	numToQuery := numToGet.Num
	response, err := internalQueryNumber(numToQuery)
	if err != nil {
		c.Error(err)
		return
	}

	var guesses []entities.GuessType

	for _, g := range response.GuessList {
		guesses = append(guesses, entities.GuessType{
			FoundBy: g.Guesser,
			FoundAt: time.Unix(g.Time, 0),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Accepted",
		"num_got": response.Num,
		"guesses": guesses,
	})
}

func internalQueryNumber(numToQuery int64) (*numberspb.QueryNumberResponse, error) {
	numClient, err := getNumClient()
	if err != nil {
		return nil, err
	}
	response, err := numClient.QueryNumber(context.Background(), &numberspb.QueryNumberRequest{Num: numToQuery})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Guessers Server
var guessServerAddr = "localhost:6000"
var guessersClient guesserspb.GuessersClient
var guessOnce sync.Once

func getGuessersClientInternal(server string) (guesserspb.GuessersClient, error) {
	if guessersClient != nil {
		return guessersClient, nil
	}
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	guessersClient = guesserspb.NewGuessersClient(conn)
	return guessersClient, nil
}

func getGuessersClient() (guesserspb.GuessersClient, error) {
	var err error
	guessOnce.Do(func() {
		guessersClient, err = getGuessersClientInternal(guessServerAddr)
	})
	if err != nil {
		return nil, err
	}
	return guessersClient, nil
}

type addGuesserInput struct {
	BeginAt     int64 `json:"begin_at"`
	IncrementBy int64 `json:"increment_by"`
	Sleep       int64 `json:"sleep"`
}

func addGuesser(c *gin.Context) {
	var guesserToAdd addGuesserInput
	if err := c.ShouldBindJSON(&guesserToAdd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	guessClient, err := getGuessersClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := guessClient.AddGuesser(context.Background(), &guesserspb.AddGuesserRequest{
		BeginAt:     guesserToAdd.BeginAt,
		IncrementBy: guesserToAdd.IncrementBy,
		Sleep:       guesserToAdd.Sleep,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":           "Accepted",
		"guesser_id_added": response.GuesserID,
	})
}

type removeGuesserInput struct {
	GuesserID int64 `json:"guesser_id"`
}

func removeGuesser(c *gin.Context) {
	var guesserToRemove removeGuesserInput
	if err := c.ShouldBindJSON(&guesserToRemove); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	guessClient, err := getGuessersClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := guessClient.RemoveGuesser(context.Background(), &guesserspb.RemoveGuesserRequest{
		GuesserID: guesserToRemove.GuesserID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "Accepted",
		"guesser_removed": response.GuesserID,
	})
}

type queryGuesserInput struct {
	GuesserID int64 `json:"guesser_id"`
}

func queryGuesser(c *gin.Context) {
	var guesserToQuery queryGuesserInput
	if err := c.ShouldBindJSON(&guesserToQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	guessClient, err := getGuessersClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := guessClient.QueryGuesser(context.Background(), &guesserspb.QueryGuesserRequest{
		GuesserID: guesserToQuery.GuesserID,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":          "Accepted",
		"guesser_queried": response.GuesserID,
		"guesses_made":    response.GuessList,
		"active":          response.Active,
	})
}

func queryPrimes(c *gin.Context) {
	guessClient, err := getGuessersClient()
	if err != nil {
		c.Error(err)
		return
	}
	response, err := guessClient.QueryPrimes(context.Background(), &guesserspb.QueryPrimesRequest{})
	if err != nil {
		c.Error(err)
		return
	}

	var primes []entities.Prime

	for _, p := range response.Primes {
		var primeDets []entities.PrimeDetails
		for _, dets := range p.PrimeList {
			primeDets = append(primeDets, entities.PrimeDetails{
				GuesserID: dets.GuesserID,
				Time:      time.Unix(dets.Time, 0),
				OriginNum: dets.OriginNum,
			})
		}
		primes = append(primes, entities.Prime{
			Prime:     p.Prime,
			PrimeDets: primeDets,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      response.Ok,
		"primes_list": response.Primes,
	})
}

type ApiServer struct {
	api.UnimplementedGuessNumsServer
	RedisManage guessers_counters.Manager
}

func (s *ApiServer) GuessNum(stream api.GuessNums_GuessNumServer) error {
	for {
		fmt.Println("API server, GuessNum, receiving guesses on stream")
		guessReq, err := stream.Recv()
		fmt.Printf("API server, GuessNum, Received %v on channel\n", guessReq.Num)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		id, i := guessReq.GuesserID, guessReq.Num
		response, err := internalQueryNumber(i)
		if err != nil {
			err = stream.Send(&api.NumGuessResponse{
				Ok:        false,
				Err:       err.Error(),
				Found:     false,
				Num:       0,
				GuesserID: id,
			})
			if err != nil {
				return err
			}
			continue
		}
		_, err = s.RedisManage.IncreaseGuesserCounter(id)
		if err != nil {
			err = stream.Send(&api.NumGuessResponse{
				Ok:        false,
				Err:       err.Error(),
				Found:     false,
				Num:       0,
				GuesserID: id,
			})
			if err != nil {
				return err
			}
			continue
		}
		fmt.Printf("API server, GuessNum, Will send %v on channel\n", response.Num)
		err = stream.Send(&api.NumGuessResponse{
			Ok:        true,
			Err:       "",
			Found:     response.Ok,
			Num:       response.Num,
			GuesserID: id,
		})
		if err != nil {
			return err
		}
		fmt.Printf("API server, GuessNum, Sent %v on channel\n", response.Num)

	}
}

func RealApi() int {
	// Guessing numbers server
	go serverForGuessRequests()

	//API client
	router := gin.Default()
	router.POST("/addNum", addNum)
	router.DELETE("/removeNum", removeNum)
	router.GET("/queryNumber", queryNumber)

	router.POST("/addGuesser", addGuesser)
	router.DELETE("/removeGuesser", removeGuesser)
	router.GET("/queryGuesser", queryGuesser)
	router.GET("/queryPrimes", queryPrimes)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("%V", err)
	}
	return 0
}

func serverForGuessRequests() {
	s := grpc.NewServer()
	rm := guessers_counters.NewManager(redis.NewRedisGuessersCounter())
	as := ApiServer{RedisManage: *rm}
	api.RegisterGuessNumsServer(s, &as)
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	if err := s.Serve(lis); err != nil { // assigns lis' port to s
		log.Fatalf("Recieved the following error : %v", err)
	}
}
