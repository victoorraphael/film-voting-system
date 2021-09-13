package service

import (
	"context"
	"fmt"
	"github.com/victoorraphael/film-voting-system/models"
	filmpb "github.com/victoorraphael/film-voting-system/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilmServer struct {
	Collection *mongo.Collection
	DbCtx      context.Context
	filmpb.UnimplementedFilmServiceServer
}

func (f *FilmServer) mustEmbedUnimplementedFilmServiceServer() {}

func (f *FilmServer) CreateFilm(_ context.Context, message *filmpb.CreateFilmMessage) (*filmpb.CreateFilmResponse, error) {

	film := message.GetFilm()

	data := filmpb.Film{
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

	return &filmpb.CreateFilmResponse{Film: film}, nil
}

func (f *FilmServer) GetFilm(ctx context.Context, message *filmpb.GetFilmMessage) (*filmpb.GetFilmResponse, error) {
	oid, err := primitive.ObjectIDFromHex(message.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error to convert ObjectId: %v", err))
	}

	result := f.Collection.FindOne(ctx, bson.M{"_id": oid})
	data := models.Film{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find film with ObjectId %s: %v", message.Id, err))
	}

	film := filmpb.Film{
		Id:        oid.Hex(),
		Name:      data.Name,
		Upvotes:   data.Upvotes,
		Downvotes: data.Downvotes,
		Score:     data.Score,
	}

	response := filmpb.GetFilmResponse{Film: &film}

	return &response, nil
}

func (f *FilmServer) UpvoteFilm(ctx context.Context, message *filmpb.UpvoteFilmMessage) (*filmpb.UpvoteFilmResponse, error) {
	oid, err := primitive.ObjectIDFromHex(message.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error to convert ObjectId: %v", err))
	}

	filter := bson.M{"_id": oid}

	_, err = f.Collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"upvotes": 1, "score": 1}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find film with ObjectId %s: %v", message.GetId(), err))
	}

	return &filmpb.UpvoteFilmResponse{Success: true}, nil
}

func (f *FilmServer) DownvoteFilm(ctx context.Context, message *filmpb.DownvoteFilmMessage) (*filmpb.DownvoteFilmResponse, error) {
	oid, err := primitive.ObjectIDFromHex(message.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error to convert ObjectId: %v", err))
	}

	filter := bson.M{"_id": oid}

	_, err = f.Collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"downvotes": 1, "score": -1}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find film with ObjectId %s: %v", message.GetId(), err))
	}

	return &filmpb.DownvoteFilmResponse{Success: true}, nil
}

func (f *FilmServer) DeleteFilm(_ context.Context, message *filmpb.DeleteFilmMessage) (*filmpb.DeleteFilmResponse, error) {
	panic("implement me")
}

func (f *FilmServer) ListFilm(message *filmpb.ListFilmMessage, server filmpb.FilmService_ListFilmServer) error {
	panic("implement me")
}
