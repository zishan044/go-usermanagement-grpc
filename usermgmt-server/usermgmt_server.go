package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	user_id := int32(rand.Intn(1000))
	return &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}, nil
}

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("server listening at %v", ln.Addr())
	if err := s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
