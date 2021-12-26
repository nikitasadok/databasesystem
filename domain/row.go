package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Row struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Cells      []Cell             `json:"cells" bson:"-"`
	CellsMongo []struct {
		Value   interface{} `bson:"value"`
		ColData ColumnData  `bson:"column"`
	} `json:"-" bson:"cells"`
	TableID       primitive.ObjectID `json:"-" bson:"table_id"`
	TableIDString string             `json:"-" bson:"-"`
}

type RowRepository interface {
	Create(ctx context.Context, r *Row) error
	Update(ctx context.Context, r *Row, id string) error
	Delete(ctx context.Context, id string, tableID string) error
	GetByID(ctx context.Context, id string, tableID string) (Row, error)
	GetTableRows(ctx context.Context, tableID string) ([]Row, error)
}

type RowUsecase interface {
	Create(ctx context.Context, r *Row) error
	Update(ctx context.Context, r *Row, id string) error
	Delete(ctx context.Context, id string, tableID string) error
	GetByID(ctx context.Context, id string, tableID string) (Row, error)
	GetTableRows(ctx context.Context, tableID string) ([]Row, error)
}

//type TableRepository interface {
//	Create(ctx context.Context, t *Table) error
//	Update(ctx context.Context, t *Table) error
//	Delete(ctx context.Context, id string, dbID string) error
//	GetByID(ctx context.Context, id string, dbID string) (Table, error)
//}
//
//func NewRow(cells []Cell) *Row {
//	return &Row{
//		Cells: cells,
//	}
//}

/*func (r *Row) appendNullDataOfTpe(tpe int) {
	r.Cells = append(r.Cells, GetConcreteCellForType(tpe))
}*/
