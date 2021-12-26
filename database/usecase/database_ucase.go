package usecase

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/nikitasadok/database-system/domain"
)

type databaseUsecase struct {
	databaseRepo   domain.DatabaseRepository
	contextTimeout time.Duration
}

func NewDatabaseUsecase(dr domain.DatabaseRepository, timeout time.Duration) domain.DatabaseUsecase {
	return &databaseUsecase{
		databaseRepo:   dr,
		contextTimeout: timeout,
	}
}

func (d *databaseUsecase) Create(c context.Context, db *domain.Database) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	if !d.isValid(db) {
		return errors.New("the database with the specified name already exists")
	}

	db.CreatedAt = time.Now()
	db.UpdatedAt = time.Now()
	return d.databaseRepo.Create(ctx, db)
}

func (d *databaseUsecase) Get(c context.Context, id string) (domain.Database, error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	res, err := d.databaseRepo.Get(ctx, id)
	if err != nil {
		return domain.Database{}, err
	}

	return res, nil

}
func (d *databaseUsecase) Update(c context.Context, db *domain.Database) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	db.UpdatedAt = time.Now()
	return d.databaseRepo.Update(ctx, db)
}
func (d *databaseUsecase) Drop(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	existingDb, err := d.Get(ctx, id)
	if err != nil {
		return err
	}

	if reflect.ValueOf(existingDb).IsZero() {
		return err
	}

	return d.databaseRepo.Drop(ctx, id)
}

func (d *databaseUsecase) GetByName(c context.Context, name string) (domain.Database, error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	res, err := d.databaseRepo.GetByName(ctx, name)
	if err != nil {
		return domain.Database{}, err
	}

	return res, nil
}

func (d *databaseUsecase) isValid(db *domain.Database) bool {
	existingDb, err := d.GetByName(context.Background(), db.Name)
	if err != nil {
		return true
	}

	if existingDb.Name == db.Name {
		return false
	}

	return true
}
