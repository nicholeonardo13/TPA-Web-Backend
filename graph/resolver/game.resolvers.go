package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nicholeonardo13/gqlgen-todos/database"
	"github.com/nicholeonardo13/gqlgen-todos/graph/middleware"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
)

func (r *queryResolver) Games(ctx context.Context) ([]*model.Game, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	user := middleware.ForContext(ctx)
	if user == nil {
		return []*model.Game{}, fmt.Errorf("access denied")
	}

	var games []*model.Game
	db.Preload("Genres").Preload("Developers").Preload("Publisher").Find(&games)

	return games, nil
}

func (r *queryResolver) GamesByGenre(ctx context.Context, genre *string) ([]*model.Game, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var games []*model.Game
	var genres []*model.Genre
	db.Where("genre_name=?", genre).Find(&genres)
	db.Model(&genres).Association("Games").Find(&games)

	return games, nil
}

func (r *queryResolver) GameByTitle(ctx context.Context, title *string) (*model.Game, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var game model.Game
	db.First(&game, "title=?", title)

	return &game, nil
}
