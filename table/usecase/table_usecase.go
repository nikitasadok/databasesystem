package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"

	"github.com/nikitasadok/database-system/domain"
)

type tableUsecase struct {
	tableRepo      domain.TableRepository
	colRepo        domain.ColumnDataRepository
	rowRepo        domain.RowRepository
	contextTimeout time.Duration
}

func NewTableUsecase(tr domain.TableRepository, colRepo domain.ColumnDataRepository, rowRepo domain.RowRepository, timeout time.Duration) domain.TableUsecase {
	return &tableUsecase{
		tableRepo:      tr,
		colRepo:        colRepo,
		rowRepo:        rowRepo,
		contextTimeout: timeout,
	}
}

func (d *tableUsecase) Create(c context.Context, t *domain.Table) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	if !d.isValid(t) || !d.areColsValid(t.Columns) {
		return errors.New("either the table with this name already exists in db or the list of columns is empty")
	}

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	err := d.tableRepo.Create(ctx, t)
	if err != nil {
		return err
	}

	for i := range t.Columns {
		t.Columns[i].TableID = t.ID
		err := d.colRepo.Create(ctx, t.Columns[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *tableUsecase) GetByID(c context.Context, id string, dbID string) (domain.Table, error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	res, err := d.tableRepo.GetByID(ctx, id, dbID)
	if err != nil {
		return domain.Table{}, err
	}

	cols, err := d.colRepo.GetAllColumnsForTable(ctx, id)
	for k := range cols {
		val := cols[k]
		res.Columns = append(res.Columns, &val)
	}

	rows, err := d.rowRepo.GetTableRows(ctx, id)
	fmt.Printf("ROWSS: %+v\n", rows)
	for k := range rows {
		val := &rows[k]
		res.Rows = append(res.Rows, val)
	}
	return res, nil

}
func (d *tableUsecase) Update(c context.Context, t *domain.Table) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	if !d.isValid(t) {
		return errors.New("the table with this name already exists")
	}
	return d.tableRepo.Update(ctx, t)
}

func (d *tableUsecase) Delete(c context.Context, id string, dbID string) error {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	existingTable, err := d.GetByID(ctx, id, dbID)
	if err != nil {
		return err
	}

	if reflect.ValueOf(existingTable).IsZero() {
		return err
	}

	if err := d.tableRepo.Delete(ctx, id, dbID); err != nil {
		return err
	}

	if err := d.colRepo.DeleteAllColumnsForTable(ctx, id); err != nil {
		return err
	}

	return nil
}
func (d *tableUsecase) GetTableProduct(ctx context.Context, t1ID, t2ID, dbID string) (domain.TablesProduct, error) {
	var res domain.TablesProduct
	t1Rows, err := d.rowRepo.GetTableRows(ctx, t1ID)
	if err != nil {
		return domain.TablesProduct{}, err
	}

	t2Rows, err := d.rowRepo.GetTableRows(ctx, t2ID)
	if err != nil {
		return domain.TablesProduct{}, err
	}

	for i := range t1Rows {
		for j := range t2Rows {
			tmp := append(t1Rows[i].Cells, t2Rows[j].Cells...)
			row := domain.Row{
				ID:      primitive.NewObjectID(),
				Cells:   tmp,
				TableID: primitive.NewObjectID(),
			}
			res.Rows = append(res.Rows, row)
		}
	}

	return res, nil

}

func (d *tableUsecase) isValid(t *domain.Table) bool {
	existingTable, err := d.GetByName(context.Background(), t.Name,t.DatabaseIDString)
	if err != nil {
		return true
	}

	if existingTable.Name == t.Name {
		return false
	}

	if len(t.Columns) == 0 {
		return false
	}

	return true
}

func (d *tableUsecase) areColsValid(cols []*domain.ColumnData) bool {
	seen := make(map[string]struct{})

	for i := range cols {
		if _, ok := seen[cols[i].Name]; ok {
			return false
		}
		if !domain.IsValidDataType(cols[i].Datatype) {
			return false
		}
		seen[cols[i].Name] = struct{}{}
	}

	return true
}

func (d *tableUsecase) GetByName(c context.Context, name,dbID string) (domain.Table, error) {
	ctx, cancel := context.WithTimeout(c, d.contextTimeout)
	defer cancel()

	res, err := d.tableRepo.GetByName(ctx, name,dbID)
	if err != nil {
		return domain.Table{}, err
	}
	return res, nil
}
