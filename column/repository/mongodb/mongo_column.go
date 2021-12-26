package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/nikitasadok/database-system/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoColumnRepository struct {
	collectionColumns *mongo.Collection
}

func NewMongoColumnRepository(collection *mongo.Collection) domain.ColumnDataRepository {
	return &mongoColumnRepository{collectionColumns: collection}
}

func (m *mongoColumnRepository) Create(ctx context.Context, col *domain.ColumnData) error {
	col.ID = primitive.NewObjectID()
	// col.TableID, _ = primitive.ObjectIDFromHex(col.TableIDString)
	col.CreatedAt = time.Now()
	col.UpdatedAt = time.Now()
	_, err := m.collectionColumns.InsertOne(ctx, col)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoColumnRepository) GetByID(ctx context.Context, id string, tableID string) (domain.ColumnData, error) {
	tableIDConv, _ := primitive.ObjectIDFromHex(tableID)
	idConv, _ := primitive.ObjectIDFromHex(id)
	var t domain.ColumnData

	err := m.collectionColumns.FindOne(ctx, bson.M{"_id": idConv, "table_id": tableIDConv}).Decode(&t)
	if err != nil {
		return domain.ColumnData{}, err
	}

	return t, nil
}

func (m *mongoColumnRepository) GetAllColumnsForTable(ctx context.Context, tableID string) (map[string]domain.ColumnData, error) {
	tableIDConv, _ := primitive.ObjectIDFromHex(tableID)
	var cols []domain.ColumnData

	results, err := m.collectionColumns.Find(ctx, bson.M{"table_id": tableIDConv})
	if err != nil {
		return nil, err
	}

	err = results.All(ctx, &cols)
	if err != nil {
		return nil, err
	}

	fmt.Printf("COLS: %+v\n", cols)

	res := make(map[string]domain.ColumnData)
	for i := range cols {
		res[cols[i].Name] = cols[i]
	}

	return res, nil
}

func (m *mongoColumnRepository) DeleteAllColumnsForTable(ctx context.Context, tableID string) error {
	tableIDConverted, _ := primitive.ObjectIDFromHex(tableID)

	_, err := m.collectionColumns.DeleteMany(ctx, bson.M{"table_id": tableIDConverted})
	if err != nil {
		return err
	}

	return nil
}
