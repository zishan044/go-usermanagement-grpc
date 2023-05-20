package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	users := make(map[string]int32)
	users["Alice"] = 32
	users["Bob"] = 22

	for name, age := range users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(`User Details
Name: %s,
Age: %d
Id: %d`, r.GetName(), r.GetAge(), r.GetId())
	}

	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\nUSER LIST:\n")
	fmt.Printf("r.UserList: %v", r.GetUsers())
}
