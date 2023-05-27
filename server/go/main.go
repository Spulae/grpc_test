package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "server/pb"
)

type server struct {
	pb.UnimplementedServerServer
	db *pgxpool.Pool
}

func (s *server) AddNumber(ctx context.Context, in *pb.AddNumberReq) (*pb.AddNumberRes, error) {
	log.Printf("AddNumber")
	return &pb.AddNumberRes{A: in.GetA() + in.GetB()}, nil
}

func (s *server) ReadFromDB(ctx context.Context, in *emptypb.Empty) (*pb.ReadFromDBRes, error) {
	log.Printf("ReadFromDB")
	var resp int32
	err := s.db.QueryRow(ctx, "select 1").Scan(&resp)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, fmt.Errorf("failed to read from db: %w", err)
	}

	return &pb.ReadFromDBRes{A: resp}, nil
}

func main() {
	// Parsing configuration and setting up connection pool
	config, _ := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

	// Database connection
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Creating TCP listener for the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Creating gRPC server
	s := grpc.NewServer()
	pb.RegisterServerServer(s, &server{db: conn})

	log.Printf("server listening at %v", lis.Addr())

	// Starting the gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
