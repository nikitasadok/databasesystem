package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	DatabaseID       primitive.ObjectID `json:"database_id" bson:"database_id"`
	DatabaseIDString string             `json:"-" bson:"-"`
	Name             string             `json:"name" bson:"name"`
	Columns          []*ColumnData      `json:"columns,omitempty" bson:"columns"`
	Rows      		 []*Row             `json:"rows,omitempty" bson:"rows"`
	// example: 10
	RowsCount 		 uint64             `json:"rowsCount" bson:"rowsCount"`
	CreatedAt 		 time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type TablesProduct struct {
	Rows []Row `json:"rows"`
}

type TableRepository interface {
	Create(ctx context.Context, t *Table) error
	Update(ctx context.Context, t *Table) error
	Delete(ctx context.Context, id string, dbID string) error
	GetByID(ctx context.Context, id string, dbID string) (Table, error)
	GetByName(ctx context.Context, name,dbID string) (Table, error)
}

type TableUsecase interface {
	Create(ctx context.Context, t *Table) error
	Update(ctx context.Context, t *Table) error
	Delete(ctx context.Context, id string, dbID string) error
	GetByID(ctx context.Context, id string, dbID string) (Table, error)
	GetTableProduct(ctx context.Context, t1ID, t2ID, dbID string) (TablesProduct, error)
	GetByName(ctx context.Context, name,dbID string) (Table, error)
}

func NewTable(name string) *Table {
	return &Table{
		Name:    name,
		Columns: []*ColumnData{},
		Rows:    []*Row{},
	}
}

// func (t *Table) AddColumn(col *ColumnData) {
// 	t.columns = append(t.columns, col)
// 	for i := 0; i < len(t.rows); i++ {
// 		t.rows[i].appendNullDataOfTpe(col.Datatype)
// 	}
// }

// func (t *Table) AddRow(row *Row) {
// 	t.rows = append(t.rows, row)
// }
