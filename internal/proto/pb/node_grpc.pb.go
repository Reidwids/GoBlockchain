// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: node.proto

package pb

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

// NodeServiceClient is the client API for NodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeServiceClient interface {
	GetNodes(ctx context.Context, in *GetNodesReq, opts ...grpc.CallOption) (*NodeList, error)
}

type nodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeServiceClient(cc grpc.ClientConnInterface) NodeServiceClient {
	return &nodeServiceClient{cc}
}

func (c *nodeServiceClient) GetNodes(ctx context.Context, in *GetNodesReq, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/pb.NodeService/GetNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServiceServer is the server API for NodeService service.
// All implementations must embed UnimplementedNodeServiceServer
// for forward compatibility
type NodeServiceServer interface {
	GetNodes(context.Context, *GetNodesReq) (*NodeList, error)
	mustEmbedUnimplementedNodeServiceServer()
}

// UnimplementedNodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServiceServer struct {
}

func (UnimplementedNodeServiceServer) GetNodes(context.Context, *GetNodesReq) (*NodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodes not implemented")
}
func (UnimplementedNodeServiceServer) mustEmbedUnimplementedNodeServiceServer() {}

// UnsafeNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServiceServer will
// result in compilation errors.
type UnsafeNodeServiceServer interface {
	mustEmbedUnimplementedNodeServiceServer()
}

func RegisterNodeServiceServer(s grpc.ServiceRegistrar, srv NodeServiceServer) {
	s.RegisterService(&NodeService_ServiceDesc, srv)
}

func _NodeService_GetNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNodesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).GetNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.NodeService/GetNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).GetNodes(ctx, req.(*GetNodesReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeService_ServiceDesc is the grpc.ServiceDesc for NodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.NodeService",
	HandlerType: (*NodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNodes",
			Handler:    _NodeService_GetNodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
