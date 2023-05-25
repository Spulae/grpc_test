package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"server/pb"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	db *pgx.Conn
	pb.UnimplementedServerServer
}

func (s *server) AddNumber(ctx context.Context, in *pb.AddNumberReq) (*pb.AddNumberRes, error) {
	log.Printf("AddNumber")
	return &pb.AddNumberRes{A: in.GetA() + in.GetB()}, nil
}

func (s *server) ReadFromDB(ctx context.Context, in *emptypb.Empty) (*pb.ReadFromDBRes, error) {
	log.Printf("ReadFromDB")
	var resp int32
	err := s.db.QueryRow(context.Background(), "select 1").Scan(&resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return &pb.ReadFromDBRes{A: resp}, nil
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServerServer(s, &server{db: conn})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
