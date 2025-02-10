package graph

import (
	"ozon-GraphQL/graph/model"
	"ozon-GraphQL/internal/database"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo             database.Repository
	CommentObservers map[string]chan *model.Comment
	mu               sync.Mutex
}

func NewResolver(Repo database.Repository) *Resolver {
	return &Resolver{
		Repo:             Repo,
		CommentObservers: make(map[string]chan *model.Comment),
	}
}
