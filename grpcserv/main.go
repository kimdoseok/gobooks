package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "grpcserv/pb"
	"grpcserv/sy"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getPool(ctx context.Context) *pgxpool.Pool {
	dburl := "postgres://doseok:kim7795004@localhost:5432/gobooks"
	config, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		panic(err)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
			log.Printf("%s: %v", msg, data)
		}),
		LogLevel: tracelog.LogLevelInfo,
		//Logger:   tracelog.LoggerFunc(log.Printf),
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	//pool, err := pgxpool.New(ctx, dburl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}
	return pool
}

func runGrpc(ctx context.Context, pool *pgxpool.Pool) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	repo_user := sy.NewUserRepo(pool, ctx)
	pb.RegisterUserServiceServer(grpcServer, sy.NewUserService(repo_user))

	log.Println("gRPC server listening on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func runGateway(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register the proxy handler pointing to our background gRPC server port
	//err := pb.RegisterUserServiceHandlerServer(ctx, mux, opts)
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	log.Println("Serving gRPC-Gateway REST API on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}

}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool := getPool(ctx)
	go runGrpc(ctx, pool)
	runGateway(ctx)

}
