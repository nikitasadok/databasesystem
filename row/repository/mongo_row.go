package repository

import (
	"context"
	"fmt"
	"github.com/nikitasadok/database-system/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mongoRowRepository struct {
	collectionTables *mongo.Collection
	collectionRows   *mongo.Collection
}

func NewMongoRowRepository(collectionTables, collectionRows *mongo.Collection) domain.RowRepository {
	return &mongoRowRepository{
		collectionTables: collectionTables,
		collectionRows:   collectionRows,
	}
}

func (m mongoRowRepository) Create(ctx context.Context, r *domain.Row) error {
	r.ID = primitive.NewObjectID()
	r.TableID, _ = primitive.ObjectIDFromHex(r.TableIDString)
	r.CellsMongo = make([]struct {
		Value   interface{}       `bson:"value"`
		ColData domain.ColumnData `bson:"column"`
	}, len(r.Cells))
	for i := range r.Cells {
		r.CellsMongo[i].Value = r.Cells[i].GetValue()
		r.CellsMongo[i].ColData.Name = r.Cells[i].GetColumnName()
		r.CellsMongo[i].ColData.Datatype = r.Cells[i].GetDatatype()
	}
	if _, err := m.collectionRows.InsertOne(ctx, bson.M{"_id": r.ID, "table_id": r.TableID, "cells": r.CellsMongo}); err != nil {
		return err
	}

	_, err := m.collectionTables.UpdateOne(ctx, bson.M{"_id": r.TableID}, bson.M{"$inc": bson.M{"rowsCount": 1}})
	if err != nil {
		return err
	}

	return nil
}

func (m mongoRowRepository) Update(ctx context.Context, r *domain.Row, id string) error {
	update := bson.D{{"$set", bson.D{
		{"updatedAt", time.Now()},
		{"cells", r.Cells},
	}}}

	tableID, _ := primitive.ObjectIDFromHex(r.TableIDString)
	ID, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("r.ID: ", r.ID)
	fmt.Println("tableID: ", tableID)
	o, err := m.collectionRows.UpdateOne(ctx, bson.M{"_id": ID, "table_id": tableID}, update)
	fmt.Printf("OOOO: %+v\n", o)
	if err != nil {
		fmt.Println("ERRROR: ", err)
		return err
	}

	return nil
}

func (m mongoRowRepository) Delete(ctx context.Context, id string, tableID string) error {
	idConverted, _ := primitive.ObjectIDFromHex(id)
	tableIDConverted, _ := primitive.ObjectIDFromHex(tableID)
	_, err := m.collectionRows.DeleteOne(ctx, bson.M{"_id": idConverted, "table_id": tableIDConverted})
	if err != nil {
		return err
	}

	_, err = m.collectionTables.UpdateOne(ctx, bson.M{"_id": tableIDConverted}, bson.M{"$inc": bson.M{"rowsCount": -1}})
	return err
}

func (m mongoRowRepository) GetByID(ctx context.Context, id string, tableID string) (domain.Row, error) {
	tableIDConv, _ := primitive.ObjectIDFromHex(tableID)
	idConv, _ := primitive.ObjectIDFromHex(id)
	var r domain.Row

	err := m.collectionRows.FindOne(ctx, bson.M{"_id": idConv, "table_id": tableIDConv}).Decode(&r)
	if err != nil {
		return domain.Row{}, err
	}

	return r, nil
}

func (m mongoRowRepository) GetTableRows(ctx context.Context, tableID string) ([]domain.Row, error) {
	tableIDConv, _ := primitive.ObjectIDFromHex(tableID)

	var rows []domain.Row
	c, err := m.collectionRows.Find(ctx, bson.M{"table_id": tableIDConv})
	if err != nil {
		fmt.Println("Err: ", err)
		return nil, err
	}

	err = c.All(ctx, &rows)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}

	for i := range rows {
		for j := range rows[i].CellsMongo {
			cell := doCast(rows[i].CellsMongo[j], rows[i].CellsMongo[j].ColData)
			rows[i].Cells = append(rows[i].Cells, cell)
		}
	}

	return rows, nil
}

func doCast(cell struct {
	Value   interface{}       `bson:"value"`
	ColData domain.ColumnData `bson:"column"`
}, colData domain.ColumnData) domain.Cell {
	var cellUpd domain.Cell

	switch colData.Datatype {
	case domain.INTEGER:
		cellUpd = &domain.IntegerCell{}
	case domain.CHAR:
		cellUpd = &domain.CharCell{}
	case domain.REAL:
		cellUpd = &domain.RealCell{}
	case domain.STRING:
		cellUpd = &domain.StringCell{}
	case domain.INTEGERINTERVAL:
		cellUpd = &domain.IntegerIntervalCell{}
	case domain.TEXTFILE:
		cellUpd = &domain.TextFileCell{}
	default:
		return nil
	}

	err := cellUpd.SetValue(cell.Value)
	if err != nil {
		fmt.Println("err: ", err)
		return nil
	}
	cellUpd.SetColumnData(colData)

	return cellUpd
}
