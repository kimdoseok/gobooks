package main

import (
	"context"
	"log"
	"time"

	pb "grpclient/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the server using insecure credentials (HTTP/2 without TLS for local dev)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client_user := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client_user.GetUserByLogin(ctx, &pb.UserLoginRequest{YusCode: "ADMIN", YusPasswd: "admin"})
	if err != nil {
		log.Fatalf("could not find user: %v", err)
	}

	log.Printf("Server Response1: %+v", resp)

	resp, err = client_user.GetUserByCode(ctx, &pb.UserRequest{YusCode: "TEST"})
	if err != nil {
		log.Fatalf("could not find user: %v", err)
	}
	log.Printf("Server Response2: %+v", resp)
	newid := resp.Id
	if newid > 0 {
		log.Printf("New user ID: %d", newid)
		resp1, err := client_user.DeleteUserByCode(ctx, &pb.UserRequest{YusCode: "TEST"})
		if err != nil {
			log.Fatalf("could not delete user: %+v, %+v", err, resp1)
		}
		log.Printf("Deleted user with ID: %d", newid)
	}

	resp2, err := client_user.AddUser(ctx, &pb.UserAddRequest{
		YusCode:   "TEST",
		YusPasswd: "test",
		YusName:   "Test User",
		YusMemo:   "This is a test user.",
		YusActive: true,
	})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("Server Response2: %+v", resp2)

	resp, err = client_user.GetUserById(ctx, &pb.UserRequest{Id: resp2.Id})
	if err != nil {
		log.Fatalf("could not find user: %v", err)
	}
	log.Printf("Server Response3: %+v", resp)

	resp2, err = client_user.SetUser(ctx, &pb.UserSetRequest{
		Id:        resp.Id,
		YusCode:   "TEST",
		YusPasswd: "test",
		YusName:   "Test User",
		YusMemo:   "This is a test user.",
		YusActive: true,
	})
	if err != nil {
		log.Fatalf("could not set user: %v", err)
	}
	log.Printf("Server Response4: %+v", resp2)

	if resp2.Success {
		resp3, err := client_user.DeleteUserById(ctx, &pb.UserRequest{Id: resp2.Id})
		if err != nil {
			log.Fatalf("could not delete user: %+v, %+v", err, resp3)
		}
		log.Printf("Deleted user with ID: %d", resp2.Id)
	}
	log.Print("Client finished successfully.")
}
