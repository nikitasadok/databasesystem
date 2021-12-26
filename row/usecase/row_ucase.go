package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/nikitasadok/database-system/domain"
	"sort"
	"strings"
	"time"
)

type rowUsecase struct {
	rowRepo        domain.RowRepository
	colsRepo       domain.ColumnDataRepository
	contextTimeout time.Duration
}

func NewRowUsecase(rowRepo domain.RowRepository, colsRepo domain.ColumnDataRepository, timeout time.Duration) domain.RowUsecase {
	return &rowUsecase{
		rowRepo:        rowRepo,
		colsRepo:       colsRepo,
		contextTimeout: timeout,
	}
}

func (ru *rowUsecase) Create(c context.Context, r *domain.Row) error {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	cols, err := ru.colsRepo.GetAllColumnsForTable(ctx, r.TableIDString)
	if err != nil {
		return err
	}

	for i := range r.Cells {
		colData, ok := cols[r.Cells[i].GetColumnName()]
		if !ok {
			return fmt.Errorf("unknown column name %s\n", r.Cells[i].GetColumnName())
		}
		r.Cells[i].SetColumnData(colData)
	}

	return ru.rowRepo.Create(ctx, r)
}

func (ru *rowUsecase) Update(c context.Context, r *domain.Row, id string) error {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	rowBefore, err := ru.rowRepo.GetByID(ctx, id, r.TableIDString)
	if err != nil {
		return err
	}

	if !isValidUpdate(rowBefore, *r) {
		return errors.New("Invalid update")
	}

	return ru.rowRepo.Update(ctx, r, id)
}

func (ru *rowUsecase) Delete(c context.Context, id string, tableID string) error {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	return ru.rowRepo.Delete(ctx, id, tableID)
}

func (ru *rowUsecase) GetByID(ctx context.Context, id string, tableID string) (domain.Row, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	return ru.rowRepo.GetByID(ctx, id, tableID)
}

func (ru *rowUsecase) GetTableRows(ctx context.Context, tableID string) ([]domain.Row, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	return ru.rowRepo.GetTableRows(ctx, tableID)
}

func isValidUpdate(rowBefore, rowAfter domain.Row) bool {
	if len(rowBefore.CellsMongo) != len(rowAfter.Cells) {
		return false
	}

	sort.SliceStable(rowBefore.Cells, func(i, j int) bool {
		res := strings.Compare(rowBefore.CellsMongo[i].ColData.Name, rowBefore.CellsMongo[j].ColData.Name)
		if res == -1 {
			return true
		}
		return false
	})

	sort.SliceStable(rowAfter.Cells, func(i, j int) bool {
		res := strings.Compare(rowAfter.Cells[i].GetColumnName(), rowAfter.Cells[j].GetColumnName())
		if res == -1 {
			return true
		}
		return false
	})

	for i := range rowBefore.Cells {
		if rowBefore.CellsMongo[i].ColData.Name != rowAfter.Cells[i].GetColumnName() ||
			rowBefore.CellsMongo[i].ColData.Datatype != rowAfter.Cells[i].GetDatatype() {
			fmt.Println(i)
			return false
		}
	}

	return true
}
