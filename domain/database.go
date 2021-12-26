package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Database struct {

	// example: "0000000000aa"
	ID         primitive.ObjectID `json:"id" bson:"_id"`

	// required: true
	// example: production
	Name       string             `json:"name" bson:"name"`
	Tables     []*Table           `json:"tables,omitempty" bson:"tables"`
	TableCount int                `json:"tableCount" bson:"tableCount"`
	CreatedAt  time.Time          `json:"-" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"-" bson:"updatedAt"`
}

type DatabaseRepository interface {
	Create(ctx context.Context, db *Database) error
	Get(ctx context.Context, id string) (Database, error)
	Update(ctx context.Context, db *Database) error
	Drop(ctx context.Context, id string) error
	GetByName(ctx context.Context, name string) (Database, error)
}

type DatabaseUsecase interface {
	Create(ctx context.Context, db *Database) error
	Get(ctx context.Context, id string) (Database, error)
	Update(ctx context.Context, db *Database) error
	Drop(ctx context.Context, id string) error
	GetByName(ctx context.Context, name string) (Database, error)
}

func NewDatabase() Database {
	return Database{}
}
