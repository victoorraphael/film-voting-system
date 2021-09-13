package filmpb

import (
	"context"
	"fmt"
	"github.com/victoorraphael/film-voting-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilmServer struct {
	Collection *mongo.Collection
	DbCtx      context.Context
}

func (f *FilmServer) mustEmbedUnimplementedFilmServiceServer() {}

func (f *FilmServer) CreateFilm(_ context.Context, message *CreateFilmMessage) (*CreateFilmResponse, error) {

	film := message.GetFilm()

	data := Film{
		Name:      film.GetName(),
		Upvotes:   film.GetUpvotes(),
		Downvotes: film.GetDownvotes(),
		Score:     film.GetScore(),
	}

	result, err := f.Collection.InsertOne(f.DbCtx, &data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	oid := result.InsertedID.(primitive.ObjectID)

	film.Id = oid.Hex()

	return &CreateFilmResponse{Film: film}, nil
}

func (f *FilmServer) GetFilm(ctx context.Context, message *GetFilmMessage) (*GetFilmResponse, error) {
	oid, err := primitive.ObjectIDFromHex(message.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error to convert ObjectId: %v", err))
	}

	result := f.Collection.FindOne(ctx, bson.M{"_id": oid})
	data := models.Film{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find film with ObjectId %s: %v", message.Id, err))
	}

	film := Film{
		Id:        oid.Hex(),
		Name:      data.Name,
		Upvotes:   data.Upvotes,
		Downvotes: data.Downvotes,
		Score:     data.Score,
	}

	response := GetFilmResponse{Film: &film}

	return &response, nil
}

func (f *FilmServer) UpvoteFilm(_ context.Context, message *UpvoteFilmMessage) (*UpvoteFilmResponse, error) {
	panic("implement me")
}

func (f *FilmServer) DownvoteFilm(_ context.Context, message *DownvoteFilmMessage) (*DownvoteFilmResponse, error) {
	panic("implement me")
}

func (f *FilmServer) DeleteFilm(_ context.Context, message *DeleteFilmMessage) (*DeleteFilmResponse, error) {
	panic("implement me")
}

func (f *FilmServer) ListFilm(message *ListFilmMessage, server FilmService_ListFilmServer) error {
	panic("implement me")
}
