package usecase

import (
	"context"
	"github.com/nikitasadok/database-system/domain"
	"time"
)

type columnDataUsecase struct {
	columnDataRepo domain.ColumnDataRepository
	contextTimeout time.Duration
}

func (c *columnDataUsecase) Create(cntx context.Context, data *domain.ColumnData) error {
	ctx, cancel := context.WithTimeout(cntx, c.contextTimeout)
	defer cancel()

	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return c.columnDataRepo.Create(ctx, data)
}

func (c *columnDataUsecase) GetByID(cntx context.Context, id string, tableID string) (domain.ColumnData, error) {
	ctx, cancel := context.WithTimeout(cntx, c.contextTimeout)
	defer cancel()

	res, err := c.columnDataRepo.GetByID(ctx, id, tableID)
	return res, err
}

func (c *columnDataUsecase) DeleteAllColumnsForTable(cntx context.Context, tableID string) error {
	ctx, cancel := context.WithTimeout(cntx, c.contextTimeout)
	defer cancel()

	return c.columnDataRepo.DeleteAllColumnsForTable(ctx, tableID)
}
