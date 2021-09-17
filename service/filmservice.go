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

func (f *FilmServer) CreateFilm(ctx context.Context, message *filmpb.CreateFilmMessage) (*filmpb.CreateFilmResponse, error) {

	film := message.GetFilm()

	if film.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Name cannot be empty!"))
	}

	if film.GetId() != "" || film.GetScore() != 0 || film.GetDownvotes() != 0 || film.GetUpvotes() != 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Film must be created only with name e votes equals to zero!"))
	}

	data := filmpb.Film{
		Name:      film.GetName(),
		Upvotes:   film.GetUpvotes(),
		Downvotes: film.GetDownvotes(),
		Score:     film.GetScore(),
	}

	exists := f.Collection.FindOne(ctx, bson.M{"name": data.Name})

	if exists.Err() == nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("Already exists a film with specified name"))
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

func (f *FilmServer) DeleteFilm(ctx context.Context, message *filmpb.DeleteFilmMessage) (*filmpb.DeleteFilmResponse, error) {
	oid, err := primitive.ObjectIDFromHex(message.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error to convert ObjectId: %v", err))
	}

	_, err = f.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find film with ObjectId: %v", err))
	}

	return &filmpb.DeleteFilmResponse{Success: true}, nil
}

func (f *FilmServer) ListFilm(message *filmpb.ListFilmMessage, stream filmpb.FilmService_ListFilmServer) error {
	data := &models.Film{}

	films, err := f.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	defer films.Close(context.Background())

	for films.Next(context.Background()) {
		err := films.Decode(data)

		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Error to decode data: %v", err))
		}

		stream.Send(&filmpb.ListFilmResponse{Film: &filmpb.Film{
			Id:        data.ID.Hex(),
			Name:      data.Name,
			Upvotes:   data.Upvotes,
			Downvotes: data.Downvotes,
			Score:     data.Score,
		}})
	}

	if err := films.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknow internal error: %v", err))
	}

	return nil
}
