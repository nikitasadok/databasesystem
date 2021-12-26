package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColumnData struct {
	ID            primitive.ObjectID `json:"-" bson:"_id"`
	TableID       primitive.ObjectID `json:"-" bson:"table_id"`
	TableIDString string             `json:"-" bson:"-"`
	Name          string             `json:"name" bson:"name"`
	Datatype      string             `json:"datatype" bson:"datatype"`
	CreatedAt     time.Time          `json:"-" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"-" bson:"updatedAt"`
}

type ColumnDataRepository interface {
	Create(ctx context.Context, data *ColumnData) error
	GetByID(ctx context.Context, id string, tableID string) (ColumnData, error)
	GetAllColumnsForTable(ctx context.Context, tableID string) (map[string]ColumnData, error)
	DeleteAllColumnsForTable(ctx context.Context, tableID string) error
}

type ColumnDataUsecase interface {
	Create(ctx context.Context, data *ColumnData) error
	GetByID(ctx context.Context, id string, tableID string) (ColumnData, error)
	DeleteAllColumnsForTable(ctx context.Context, tableID string) error
}

func NewColumnData(name string, datatype string) *ColumnData {
	return &ColumnData{
		Name:     name,
		Datatype: datatype,
	}
}
