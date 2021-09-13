package main

import (
	"github.com/labstack/gommon/log"
	"github.com/victoorraphael/film-voting-system/db"
	filmpb "github.com/victoorraphael/film-voting-system/proto"
	"google.golang.org/grpc"
	"net"
)

func main() {
	conn, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen on port 8000: %v", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	collection := db.GetCollection()
	dbCtx := db.GetContext()

	s := filmpb.FilmServer{Collection: collection, DbCtx: dbCtx}

	grpcServer := grpc.NewServer()

	filmpb.RegisterFilmServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 8000: %v", err)
	}
}
