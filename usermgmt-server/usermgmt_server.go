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

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	user_list *pb.UserList
}

func (server *UserManagementServer) Run() error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", ln.Addr())
	return s.Serve(ln)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	user_id := int32(rand.Intn(1000))
	new_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}
	s.user_list.Users = append(s.user_list.Users, new_user)
	return new_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
}

func main() {
	user_mgmt_server := NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatal(err)
	}
}
