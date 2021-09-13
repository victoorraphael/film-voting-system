package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/victoorraphael/film-voting-system/db"
	filmpb "github.com/victoorraphael/film-voting-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	collection := db.GetCollection("rankTest")
	dbCtx := db.GetContext()

	fs := FilmServer{
		Collection: collection,
		DbCtx:      dbCtx,
	}

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	filmpb.RegisterFilmServiceServer(s, &fs)

	log.Println("serving grpc service ...")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestInputFilm(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := filmpb.NewFilmServiceClient(conn)

	film := filmpb.Film{
		Name:      "teste" + time.Now().String(),
		Upvotes:   0,
		Downvotes: 0,
		Score:     0,
	}

	req := filmpb.CreateFilmMessage{Film: &film}

	resp, err := client.CreateFilm(ctx, &req)
	if err != nil {
		t.Fatalf("CreateFilm failed: %v", err)
	}

	assert.Equal(t, resp.Film.Name, film.Name)
}

func TestGetFilm(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := filmpb.NewFilmServiceClient(conn)

	filmMock := filmpb.Film{
		Id:        "613f75f24c15bc7b7d4f6404",
		Name:      "teste1",
		Upvotes:   2,
		Downvotes: 1,
		Score:     1,
	}

	req := filmpb.GetFilmMessage{Id: "613f75f24c15bc7b7d4f6404"}
	res, err := client.GetFilm(ctx, &req)
	if err != nil {
		t.Fatalf("Failed to GetFilm: %v", err)
	}

	assert.Equal(t, res.Film.Id, filmMock.Id)
	assert.Equal(t, res.Film.Name, filmMock.Name)
	assert.Equal(t, res.Film.Upvotes, filmMock.Upvotes)
	assert.Equal(t, res.Film.Downvotes, filmMock.Downvotes)
	assert.Equal(t, res.Film.Score, filmMock.Score)
}

func TestUpvoteFilm(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := filmpb.NewFilmServiceClient(conn)

	filmMock := filmpb.UpvoteFilmMessage{Id: "613fa8c79cd243857a44a5bb"}

	res, err := client.UpvoteFilm(ctx, &filmMock)
	if err != nil {
		t.Fatalf("Failed to Upvote film: %v", err)
	}

	assert.Equal(t, res.Success, true)
}

func TestDownvoteFilm(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := filmpb.NewFilmServiceClient(conn)

	filmMock := filmpb.DownvoteFilmMessage{Id: "613fa8c79cd243857a44a5bb"}

	res, err := client.DownvoteFilm(ctx, &filmMock)
	if err != nil {
		t.Fatalf("Failed to Upvote film: %v", err)
	}

	assert.Equal(t, res.Success, true)
}
