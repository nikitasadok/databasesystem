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

type mongoTableRepository struct {
	collectionTables    *mongo.Collection
	collectionDatabases *mongo.Collection
}

func NewMongoDatabaseRepository(collectionTables, collectionDatabases *mongo.Collection) domain.TableRepository {
	return &mongoTableRepository{
		collectionTables:    collectionTables,
		collectionDatabases: collectionDatabases,
	}
}

func (m *mongoTableRepository) Create(ctx context.Context, t *domain.Table) error {
	t.ID = primitive.NewObjectID()
	t.DatabaseID, _ = primitive.ObjectIDFromHex(t.DatabaseIDString)
	if _, err := m.collectionTables.InsertOne(ctx, bson.M{"_id": t.ID, "name": t.Name, "createdAt": t.CreatedAt, "updatedAt": t.UpdatedAt, "database_id": t.DatabaseID}); err != nil {
		return err
	}

	_, err := m.collectionDatabases.UpdateOne(ctx, bson.M{"_id": t.DatabaseID}, bson.M{"$inc": bson.M{"tableCount": 1}})
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoTableRepository) Update(ctx context.Context, t *domain.Table) error {
	update := bson.D{{"$set", bson.D{
		{"updatedAt", time.Now()},
		{"name", t.Name},
	}}}

	_, err := m.collectionTables.UpdateOne(ctx, bson.M{"_id": t.ID, "database_id": t.DatabaseID}, update)
	if err != nil {
		return err
	}

	return nil
}

// todo delete via mongo query
func (m *mongoTableRepository) Delete(ctx context.Context, id string, dbID string) error {
	idConverted, _ := primitive.ObjectIDFromHex(id)
	dbIDConverted, _ := primitive.ObjectIDFromHex(dbID)
	_, err := m.collectionTables.DeleteOne(ctx, bson.M{"_id": idConverted, "database_id": dbIDConverted})
	if err != nil {
		return err
	}

	_, err = m.collectionDatabases.UpdateOne(ctx, bson.M{"_id": dbIDConverted}, bson.M{"$inc": bson.M{"tableCount": -1}})
	return err
}

func (m *mongoTableRepository) GetByID(ctx context.Context, id string, dbID string) (domain.Table, error) {
	dbIDConv, _ := primitive.ObjectIDFromHex(dbID)
	idConv, _ := primitive.ObjectIDFromHex(id)
	var t domain.Table

	err := m.collectionTables.FindOne(ctx, bson.M{"_id": idConv, "database_id": dbIDConv}).Decode(&t)
	if err != nil {
		return domain.Table{}, err
	}

	fmt.Println("ROW CNT: ", t.RowsCount)

	return t, nil
}

func (m *mongoTableRepository) GetByName(c context.Context, name, dbID string) (domain.Table, error) {
	var res domain.Table

	if err := m.collectionTables.FindOne(c, bson.M{"name": name, "database_id": dbID}).Decode(&res); err != nil {
		return domain.Table{}, err
	}

	return res, nil
}

/*
func (m *mongoTableRepository) GetAllColumns(ctx context.Context, id string, dbID string) ([]domain.ColumnData, error) {
	dbIDConv, _ := primitive.ObjectIDFromHex(dbID)
	idConv, _ := primitive.ObjectIDFromHex(id)
	var cols []domain.ColumnData
}
*/
// func (m *mongoDatabaseRepository) GetByName(ctx context.Context, name string) (domain.Database, error) {
// 	var res domain.Database

// 	if err := m.collectionTables.FindOne(ctx, bson.M{"name": name}).Decode(&res); err != nil {
// 		return domain.NewDatabase(), err
// 	}

// 	return res, nil
// }
