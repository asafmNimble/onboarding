package numbers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"onboarding/common/data/dbbackends/mongo"
	"onboarding/common/data/managers/numbers"
	numberspb "onboarding/common/grpc/numbers"
)

var numMap = make(map[int64]int64)

type NumsServer struct { //Defined a server that executes the far side functions
	numberspb.UnimplementedNumbersServer
	MongoManage numbers.Manager
}

func notFound(s string) bool{
	if s == "mongo: no documents in result" {
		return true
	}
	return false
}

// connects the function to the server
//       v
func (ns *NumsServer) AddNum(_ context.Context, numReq *numberspb.AddNumRequest) (*numberspb.AddNumResponse, error) {
	i := numReq.Num
	//_, found := numMap[i]
	//if found {
	//	return nil, errors.New("number already in database")
	//}
	//numMap[i] = i
	_, err := ns.MongoManage.Get(i)
	if err != nil && !notFound(err.Error()){
		return nil, err
	}
	if err == nil {
		return nil, errors.New("number already in database")
	}
	ns.MongoManage.AddNum(i)
	return &numberspb.AddNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (ns *NumsServer) RemoveNum(_ context.Context, numReq *numberspb.RemoveNumRequest) (*numberspb.RemoveNumResponse, error) {
	i := numReq.Num
	//_, found := numMap[i]
	//if !found {
	//	return nil, errors.New("number doesn't exist in database")
	//}
	_, err := ns.MongoManage.Get(i)
	if err != nil {
		return nil, err
	}
	//delete(numMap, i)
	ns.MongoManage.RemoveNum(i)
	return &numberspb.RemoveNumResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (ns *NumsServer) QueryNumber(_ context.Context, numReq *numberspb.QueryNumberRequest) (*numberspb.QueryNumberResponse, error) {
	i := numReq.Num
	//_, found := numMap[i]
	_, found := numMap[i]
	if !found {
		return nil, errors.New("number doesn't exist in database")
	}
	return &numberspb.QueryNumberResponse{
		Ok:  true,
		Num: i,
	}, nil
}

func (ns *NumsServer) GetNums(_ context.Context, numReq *numberspb.GetNumsRequest) (*numberspb.GetNumsResponse, error) {
	return &numberspb.GetNumsResponse{
		Ok:      true,
		NumsMap: numMap,
	}, nil
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
