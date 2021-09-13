package filmpb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilmServer struct {
	Collection *mongo.Collection
	DbCtx      context.Context
}

func (s *FilmServer) mustEmbedUnimplementedFilmServiceServer() {}

func (s *FilmServer) CreateFilm(_ context.Context, message *CreateFilmMessage) (*CreateFilmResponse, error) {

	film := message.GetFilm()

	data := Film{
		Name:      film.GetName(),
		Upvotes:   film.GetUpvotes(),
		Downvotes: film.GetDownvotes(),
		Score:     film.GetScore(),
	}

	result, err := s.Collection.InsertOne(s.DbCtx, data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	oid := result.InsertedID.(primitive.ObjectID)

	film.Id = oid.Hex()

	return &CreateFilmResponse{Film: film}, nil
}
