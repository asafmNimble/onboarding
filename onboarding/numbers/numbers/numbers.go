package numbers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"onboarding/common/data/managers/numbers"
	numbers2 "onboarding/common/grpc/numbers"
)

var numMap = make(map[int64]int64)

type NumsServer struct { //Defined a server that executes the far side functions
	numbers2.UnimplementedNumbersServer
	MongoManage numbers.Manager
}

// connects the function to the server
//       v
func (*NumsServer) AddNum(_ context.Context, numReq *numbers2.AddNumRequest) (*numbers2.AddNumResponse, error) {

	i := numReq.Num
	_, found := numMap[i]
	if found {
		return nil, errors.New("number already in database")
	}
	numMap[i] = i
	return &numbers2.AddNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (*NumsServer) GetNums(_ context.Context, numReq *numbers2.GetNumsRequest) (*numbers2.GetNumsResponse, error) {
	return &numbers2.GetNumsResponse{
		Ok:      true,
		NumsMap: numMap,
	}, nil
}

func (*NumsServer) RemoveNum(_ context.Context, numReq *numbers2.RemoveNumRequest) (*numbers2.RemoveNumResponse, error) {
	i := numReq.Num
	_, found := numMap[i]
	if !found {
		return nil, errors.New("number doesn't exist in database")
	}
	delete(numMap, i)
	return &numbers2.RemoveNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (*NumsServer) QueryNumber(_ context.Context, numReq *numbers2.QueryNumberRequest) (*numbers2.QueryNumberResponse, error) {
	i := numReq.Num
	_, found := numMap[i]
	if !found {
		return nil, errors.New("number doesn't exist in database")
	}
	return &numbers2.QueryNumberResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func RealNumbers() int {
	s := grpc.NewServer()
	ns := NumsServer{}
	numbers2.RegisterNumbersServer(s, &ns) //connects the far server with the grpc server
	lis, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("Recieved the following error : %v", err)
	}
	if err := s.Serve(lis); err != nil { // assigns lis' port to s
		log.Fatalf("Recieved the following error : %v", err)
	}
	return 0
}
