package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nicholeonardo13/gqlgen-todos/database"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
)

func (r *queryResolver) Countries(ctx context.Context) ([]*model.Country, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	var countries []*model.Country
	db.Find(&countries)

	return countries, nil
}
