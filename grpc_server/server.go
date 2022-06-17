package grpc_server

import (
	"context"
	"example.com/http_kv/cache"
	pb "example.com/http_kv/grpc_server/protos"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

// server is used to implement gRPC service kv.Cache
type server struct {
	c cache.Cache
	pb.UnimplementedCacheServer
}

// Get implements kv.Cache.Get
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	vRaw := s.c.Get(in.GetKey())

	if vRaw == nil {
		err := status.Error(codes.NotFound, "key was not found")
		return nil, err
	}

	v := fmt.Sprintf("%v", vRaw)

	return &pb.GetReply{Value: v}, nil
}

// Set implements kv.Cache.Set
func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	if in.GetTtlSeconds() > 0 {
		s.c.SetWithTTl(in.GetKey(), in.GetValue(), uint(in.GetTtlSeconds()))
	} else {
		s.c.Set(in.GetKey(), in.GetValue())
	}
	return &pb.SetReply{}, nil
}

func NewServer(c cache.Cache) *server {
	return &server{c, pb.UnimplementedCacheServer{}}
}

func (s *server) Start(port *int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcS := grpc.NewServer()
	pb.RegisterCacheServer(grpcS, s)
	log.Printf("gRPC server with %v listening at %v", s.c.Name(), lis.Addr())
	reflection.Register(grpcS)
	if err := grpcS.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
