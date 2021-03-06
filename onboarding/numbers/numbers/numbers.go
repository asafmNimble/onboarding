package numbers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"onboarding/common/data/dbbackends/mongo"
	"onboarding/common/data/entities"
	"onboarding/common/data/managers/numbers"
	numberspb "onboarding/common/grpc/numbers"
)

type NumsServer struct { //Defined a server that executes the far side functions
	numberspb.UnimplementedNumbersServer
	MongoManage numbers.Manager
}

func notFound(s string) bool {
	if s == "mongo: no documents in result" {
		return true
	}
	return false
}

// connects the function to the server
//       v
func (ns *NumsServer) AddNum(_ context.Context, numReq *numberspb.AddNumRequest) (*numberspb.AddNumResponse, error) {
	i := numReq.Num
	_, err := ns.MongoManage.GetNumber(i)
	if err != nil && !notFound(err.Error()) {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("number already in database")
	}
	_, err = ns.MongoManage.AddNum(i)
	if err != nil {
		return nil, err
	}
	return &numberspb.AddNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (ns *NumsServer) RemoveNum(_ context.Context, numReq *numberspb.RemoveNumRequest) (*numberspb.RemoveNumResponse, error) {
	i := numReq.Num

	_, err := ns.MongoManage.GetNumber(i)
	if err != nil {
		return nil, err
	}
	_, err = ns.MongoManage.RemoveNum(i)
	if err != nil {
		return nil, err
	}
	return &numberspb.RemoveNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (ns *NumsServer) QueryNumber(_ context.Context, numReq *numberspb.QueryNumberRequest) (*numberspb.QueryNumberResponse, error) {
	i := numReq.Num
	number, err := ns.MongoManage.GetNumber(i)
	var guesses []*numberspb.Guess
	if err != nil && notFound(err.Error()) {
		return &numberspb.QueryNumberResponse{
			Ok:        false,
			Num:       i,
			GuessList: guesses,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	// Number found
	for _, g := range number.Guesses {
		guesses = append(guesses, &numberspb.Guess{
			Guesser: g.FoundBy,
			Time:    g.FoundAt.Unix(),
		})
	}
	return &numberspb.QueryNumberResponse{
		Ok:        true,
		Num:       i,
		GuessList: guesses,
	}, nil
}

// TODO: write this func or another updating func
func (ns *NumsServer) UpdateGuessForNumber(n int64, guess *entities.GuessType) (string, error) {
	return "", nil
}

func RealNumbers() int {
	s := grpc.NewServer()
	mn := numbers.NewManager(mongo.NewMongoNumber(mongo.NewMongoConnector()))
	ns := NumsServer{MongoManage: *mn}
	numberspb.RegisterNumbersServer(s, &ns) //connects the far server with the grpc server
	lis, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	if err := s.Serve(lis); err != nil { // assigns lis' port to s
		log.Fatalf("Recieved the following error : %v", err)
	}
	return 0
}
