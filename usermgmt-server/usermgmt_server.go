package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
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
	user_list := &pb.UserList{}
	user_id := int32(rand.Intn(1000))
	new_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}
	readBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			user_list.Users = append(user_list.Users, new_user)
			jsonBytes, err := protojson.Marshal(user_list)
			if err != nil {
				log.Fatal(err)
			}
			if err := ioutil.WriteFile("users.json", jsonBytes, 0644); err != nil {
				log.Fatal(err)
			}
			return new_user, nil
		} else {
			log.Fatal(err)
		}
	}
	if err := protojson.Unmarshal(readBytes, user_list); err != nil {
		log.Fatal(err)
	}
	user_list.Users = append(user_list.Users, new_user)
	jsonBytes, err := protojson.Marshal(user_list)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("users.json", jsonBytes, 0644); err != nil {
		log.Fatal(err)
	}
	return new_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}
	user_list := &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, user_list); err != nil {
		log.Fatal(err)
	}
	return user_list, err
}

func main() {
	user_mgmt_server := NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatal(err)
	}
}
