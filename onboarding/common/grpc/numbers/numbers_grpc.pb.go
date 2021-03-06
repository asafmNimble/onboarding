// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package numbers

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NumbersClient is the client API for Numbers service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NumbersClient interface {
	AddNum(ctx context.Context, in *AddNumRequest, opts ...grpc.CallOption) (*AddNumResponse, error)
	// rpc getNums (getNumsRequest) returns (getNumsResponse) {}
	RemoveNum(ctx context.Context, in *RemoveNumRequest, opts ...grpc.CallOption) (*RemoveNumResponse, error)
	QueryNumber(ctx context.Context, in *QueryNumberRequest, opts ...grpc.CallOption) (*QueryNumberResponse, error)
}

type numbersClient struct {
	cc grpc.ClientConnInterface
}

func NewNumbersClient(cc grpc.ClientConnInterface) NumbersClient {
	return &numbersClient{cc}
}

func (c *numbersClient) AddNum(ctx context.Context, in *AddNumRequest, opts ...grpc.CallOption) (*AddNumResponse, error) {
	out := new(AddNumResponse)
	err := c.cc.Invoke(ctx, "/numberspb.numbers/addNum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *numbersClient) RemoveNum(ctx context.Context, in *RemoveNumRequest, opts ...grpc.CallOption) (*RemoveNumResponse, error) {
	out := new(RemoveNumResponse)
	err := c.cc.Invoke(ctx, "/numberspb.numbers/removeNum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *numbersClient) QueryNumber(ctx context.Context, in *QueryNumberRequest, opts ...grpc.CallOption) (*QueryNumberResponse, error) {
	out := new(QueryNumberResponse)
	err := c.cc.Invoke(ctx, "/numberspb.numbers/query_number", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NumbersServer is the server API for Numbers service.
// All implementations must embed UnimplementedNumbersServer
// for forward compatibility
type NumbersServer interface {
	AddNum(context.Context, *AddNumRequest) (*AddNumResponse, error)
	// rpc getNums (getNumsRequest) returns (getNumsResponse) {}
	RemoveNum(context.Context, *RemoveNumRequest) (*RemoveNumResponse, error)
	QueryNumber(context.Context, *QueryNumberRequest) (*QueryNumberResponse, error)
	mustEmbedUnimplementedNumbersServer()
}

// UnimplementedNumbersServer must be embedded to have forward compatible implementations.
type UnimplementedNumbersServer struct {
}

func (UnimplementedNumbersServer) AddNum(context.Context, *AddNumRequest) (*AddNumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNum not implemented")
}
func (UnimplementedNumbersServer) RemoveNum(context.Context, *RemoveNumRequest) (*RemoveNumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveNum not implemented")
}
func (UnimplementedNumbersServer) QueryNumber(context.Context, *QueryNumberRequest) (*QueryNumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryNumber not implemented")
}
func (UnimplementedNumbersServer) mustEmbedUnimplementedNumbersServer() {}

// UnsafeNumbersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NumbersServer will
// result in compilation errors.
type UnsafeNumbersServer interface {
	mustEmbedUnimplementedNumbersServer()
}

func RegisterNumbersServer(s grpc.ServiceRegistrar, srv NumbersServer) {
	s.RegisterService(&Numbers_ServiceDesc, srv)
}

func _Numbers_AddNum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NumbersServer).AddNum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/numberspb.numbers/addNum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NumbersServer).AddNum(ctx, req.(*AddNumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Numbers_RemoveNum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveNumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NumbersServer).RemoveNum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/numberspb.numbers/removeNum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NumbersServer).RemoveNum(ctx, req.(*RemoveNumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Numbers_QueryNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NumbersServer).QueryNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/numberspb.numbers/query_number",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NumbersServer).QueryNumber(ctx, req.(*QueryNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Numbers_ServiceDesc is the grpc.ServiceDesc for Numbers service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Numbers_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "numberspb.numbers",
	HandlerType: (*NumbersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "addNum",
			Handler:    _Numbers_AddNum_Handler,
		},
		{
			MethodName: "removeNum",
			Handler:    _Numbers_RemoveNum_Handler,
		},
		{
			MethodName: "query_number",
			Handler:    _Numbers_QueryNumber_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "numbers.proto",
}
