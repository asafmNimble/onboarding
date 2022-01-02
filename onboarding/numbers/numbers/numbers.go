package numbers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"onboarding/common/data/managers/numbers"
	numberspb "onboarding/common/grpc/numbers"
)

var numMap = make(map[int64]int64)

type NumsServer struct { //Defined a server that executes the far side functions
	numberspb.UnimplementedNumbersServer
	MongoManage numbers.Manager
}

// connects the function to the server
//       v
func (*NumsServer) AddNum(_ context.Context, numReq *numberspb.AddNumRequest) (*numberspb.AddNumResponse, error) {
	//TODO: use MongoManager to update mongo collection
	i := numReq.Num
	_, found := numMap[i]
	if found {
		return nil, errors.New("number already in database")
	}
	numMap[i] = i
	return &numberspb.AddNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (*NumsServer) GetNums(_ context.Context, numReq *numberspb.GetNumsRequest) (*numberspb.GetNumsResponse, error) {
	return &numberspb.GetNumsResponse{
		Ok:      true,
		NumsMap: numMap,
	}, nil
}

func (*NumsServer) RemoveNum(_ context.Context, numReq *numberspb.RemoveNumRequest) (*numberspb.RemoveNumResponse, error) {
	i := numReq.Num
	_, found := numMap[i]
	if !found {
		return nil, errors.New("number doesn't exist in database")
	}
	delete(numMap, i)
	return &numberspb.RemoveNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (*NumsServer) QueryNumber(_ context.Context, numReq *numberspb.QueryNumberRequest) (*numberspb.QueryNumberResponse, error) {
	i := numReq.Num
	_, found := numMap[i]
	if !found {
		return nil, errors.New("number doesn't exist in database")
	}
	return &numberspb.QueryNumberResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func RealNumbers() int {
	s := grpc.NewServer()
	ns := NumsServer{}
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
