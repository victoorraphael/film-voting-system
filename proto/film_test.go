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
		Name:      "teste1",
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
