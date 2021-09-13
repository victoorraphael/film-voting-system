// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.0
// source: film.proto

package filmpb

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

// FilmServiceClient is the client API for FilmService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilmServiceClient interface {
	CreateFilm(ctx context.Context, in *CreateFilmMessage, opts ...grpc.CallOption) (*CreateFilmResponse, error)
}

type filmServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFilmServiceClient(cc grpc.ClientConnInterface) FilmServiceClient {
	return &filmServiceClient{cc}
}

func (c *filmServiceClient) CreateFilm(ctx context.Context, in *CreateFilmMessage, opts ...grpc.CallOption) (*CreateFilmResponse, error) {
	out := new(CreateFilmResponse)
	err := c.cc.Invoke(ctx, "/film.FilmService/CreateFilm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilmServiceServer is the server API for FilmService service.
// All implementations must embed UnimplementedFilmServiceServer
// for forward compatibility
type FilmServiceServer interface {
	CreateFilm(context.Context, *CreateFilmMessage) (*CreateFilmResponse, error)
	mustEmbedUnimplementedFilmServiceServer()
}

// UnimplementedFilmServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFilmServiceServer struct {
}

func (UnimplementedFilmServiceServer) CreateFilm(context.Context, *CreateFilmMessage) (*CreateFilmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFilm not implemented")
}
func (UnimplementedFilmServiceServer) mustEmbedUnimplementedFilmServiceServer() {}

// UnsafeFilmServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilmServiceServer will
// result in compilation errors.
type UnsafeFilmServiceServer interface {
	mustEmbedUnimplementedFilmServiceServer()
}

func RegisterFilmServiceServer(s grpc.ServiceRegistrar, srv FilmServiceServer) {
	s.RegisterService(&FilmService_ServiceDesc, srv)
}

func _FilmService_CreateFilm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFilmMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilmServiceServer).CreateFilm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/film.FilmService/CreateFilm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilmServiceServer).CreateFilm(ctx, req.(*CreateFilmMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// FilmService_ServiceDesc is the grpc.ServiceDesc for FilmService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FilmService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "film.FilmService",
	HandlerType: (*FilmServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFilm",
			Handler:    _FilmService_CreateFilm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "film.proto",
}
