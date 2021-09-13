package main

import (
	"context"
	"github.com/labstack/echo/v4"
	_ "github.com/victoorraphael/film-voting-system/models"
	filmpb "github.com/victoorraphael/film-voting-system/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"time"
)

var filmClient filmpb.FilmServiceClient

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server on port 8000: %v", err)
	}

	defer conn.Close()

	filmClient = filmpb.NewFilmServiceClient(conn)

	e := echo.New()

	r := e.Group("/film")

	r.POST("/", createFilm)
	r.GET("/", listFilm)
	r.GET("/:id/", getFilmById)
	r.DELETE("/:id/", deleteFilm)
	r.POST("/upvote/:id/", upvoteFilm)
	r.POST("/downvote/:id/", downvoteFilm)

	e.Logger.Fatal(e.Start(":4000"))
}

func downvoteFilm(c echo.Context) error {
	uid := c.Param("id")

	message := filmpb.DownvoteFilmMessage{Id: uid}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	res, err := filmClient.DownvoteFilm(ctx, &message)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func upvoteFilm(c echo.Context) error {
	uid := c.Param("id")

	message := filmpb.UpvoteFilmMessage{Id: uid}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	res, err := filmClient.UpvoteFilm(ctx, &message)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func deleteFilm(c echo.Context) error {
	uid := c.Param("id")

	message := filmpb.DeleteFilmMessage{Id: uid}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	res, err := filmClient.DeleteFilm(ctx, &message)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func getFilmById(c echo.Context) error {
	uid := c.Param("id")

	message := filmpb.GetFilmMessage{Id: uid}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	res, err := filmClient.GetFilm(ctx, &message)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func listFilm(c echo.Context) error {
	film := filmpb.ListFilmMessage{}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	stream, err := filmClient.ListFilm(ctx, &film)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var response []*filmpb.Film

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		response = append(
			response,
			res.GetFilm(),
		)
	}

	return c.JSON(http.StatusOK, response)
}

func createFilm(c echo.Context) error {
	film := filmpb.Film{}
	if err := c.Bind(&film); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	defer cancel()

	fm := filmpb.CreateFilmMessage{Film: &film}

	res, err := filmClient.CreateFilm(ctx, &fm)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, &res)
}
