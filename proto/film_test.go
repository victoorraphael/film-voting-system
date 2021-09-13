package filmpb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/victoorraphael/film-voting-system/db"
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

	RegisterFilmServiceServer(s, &fs)

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
	log.Println("initializing create film test ...")
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := NewFilmServiceClient(conn)

	film := Film{
		Name:      "teste" + time.Now().String(),
		Upvotes:   0,
		Downvotes: 0,
		Score:     0,
	}

	log.Println("calling CreateMessage service ...")
	req := CreateFilmMessage{Film: &film}

	resp, err := client.CreateFilm(ctx, &req)
	if err != nil {
		t.Fatalf("CreateFilm failed: %v", err)
	}

	log.Printf("Response: %v", resp)

	assert.Equal(t, resp.Film.Name, film.Name)
}

func TestGetFilm(t *testing.T) {
	log.Println("initializing get film test ...")
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := NewFilmServiceClient(conn)

	filmMock := Film{
		Id:        "613f75f24c15bc7b7d4f6404",
		Name:      "teste1",
		Upvotes:   0,
		Downvotes: 0,
		Score:     0,
	}

	req := GetFilmMessage{Id: "613f75f24c15bc7b7d4f6404"}
	res, err := client.GetFilm(ctx, &req)
	if err != nil {
		t.Fatalf("Failed to GetFilm: %v", err)
	}

	log.Printf("Response: %v", res)

	assert.Equal(t, res.Film.Id, filmMock.Id)
	assert.Equal(t, res.Film.Name, filmMock.Name)
	assert.Equal(t, res.Film.Upvotes, filmMock.Upvotes)
	assert.Equal(t, res.Film.Downvotes, filmMock.Downvotes)
	assert.Equal(t, res.Film.Score, filmMock.Score)
}
