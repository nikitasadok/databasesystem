package mongodb

import (
	"context"
	"github.com/nikitasadok/database-system/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDatabaseRepository struct {
	collection       *mongo.Collection
	collectionTables *mongo.Collection
}

func NewMongoDatabaseRepository(collection, collectionTables *mongo.Collection) domain.DatabaseRepository {
	return &mongoDatabaseRepository{
		collection:       collection,
		collectionTables: collectionTables,
	}
}

func (m *mongoDatabaseRepository) Create(ctx context.Context, db *domain.Database) error {
	db.ID = primitive.NewObjectID()
	_, err := m.collection.InsertOne(ctx, db)

	return err
}

func (m *mongoDatabaseRepository) Get(ctx context.Context, id string) (domain.Database, error) {
	var res domain.Database

	convertedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Database{}, err
	}

	if err := m.collection.FindOne(ctx, bson.M{"_id": convertedID}).Decode(&res); err != nil {
		return domain.NewDatabase(), err
	}

	results, err := m.collectionTables.Find(ctx, bson.M{"database_id": convertedID})
	if err != nil {
		return domain.NewDatabase(), err
	}

	if err := results.All(ctx, &res.Tables); err != nil {
		return domain.NewDatabase(), err
	}

	return res, nil
}

func (m *mongoDatabaseRepository) Update(ctx context.Context, db *domain.Database) error {
	update := bson.D{{"$set", bson.D{
		{"name", db.Name},
		{"updatedAt", db.UpdatedAt},
	}}}
	_, err := m.collection.UpdateOne(ctx, bson.M{"_id": db.ID}, update)
	return err
}

func (m *mongoDatabaseRepository) Drop(ctx context.Context, id string) error {
	idConv, _ := primitive.ObjectIDFromHex(id)
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": idConv})
	if err != nil {
		return err
	}

	_, err = m.collectionTables.DeleteMany(ctx, bson.M{"database_id": idConv})
	return err
}

func (m *mongoDatabaseRepository) GetByName(ctx context.Context, name string) (domain.Database, error) {
	var res domain.Database

	if err := m.collection.FindOne(ctx, bson.M{"name": name}).Decode(&res); err != nil {
		return domain.NewDatabase(), err
	}

	return res, nil
}
