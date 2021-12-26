package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikitasadok/database-system/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

type TableHandler struct {
	TUsecase domain.TableUsecase
	RUsecase domain.RowUsecase
}

func NewTableHandler(r *gin.RouterGroup, us domain.TableUsecase, usRow domain.RowUsecase) {
	handler := &TableHandler{
		TUsecase: us,
		RUsecase: usRow,
	}

	r.POST("/database/:dbID/table", handler.Create)
	//s	r.POST("/database/:id/table/:tableID/addRow", handler.AddRow)
	r.GET("/database/:id/table/:tableID", handler.GetByID)
	r.PUT("/database/:id/table", handler.Update)
	r.DELETE("/database/:id/table/:tableID", handler.Delete)
	r.GET("/database/:id/tablesProduct", handler.TablesProduct)

	//	r.GET("/database/:id", handler.GetByID
	//	r.PUT("/database/:id", handler.Update)
	//	r.POST("/database", handler.Create)
}

func (h *TableHandler) Create(c *gin.Context) {
	var t domain.Table
	err := c.Bind(&t)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	t.DatabaseIDString = c.Param("dbID")

	ctx := c.Request.Context()
	err = h.TUsecase.Create(ctx, &t)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *TableHandler) GetByID(c *gin.Context) {
	tableID := c.Param("tableID")
	id := c.Param("id")

	ctx := c.Request.Context()
	db, err := h.TUsecase.GetByID(ctx, tableID, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, db)
}

func (h *TableHandler) Update(c *gin.Context) {
	var t domain.Table
	if err := c.Bind(&t); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()

	if err := h.TUsecase.Update(ctx, &t); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *TableHandler) Delete(c *gin.Context) {
	dbID := c.Param("id")
	tableID := c.Param("tableID")

	ctx := c.Request.Context()
	if err := h.TUsecase.Delete(ctx, tableID, dbID); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *TableHandler) TablesProduct(c *gin.Context) {
	dbID := c.Param("id")
	t1ID := c.Query("t1_id")
	t2ID := c.Query("t2_id")

	ctx := c.Request.Context()
	tableProduct, err := h.TUsecase.GetTableProduct(ctx, t1ID, t2ID, dbID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	fmt.Printf("%+v", tableProduct.Rows[0].Cells[0].GetDatatype())
	c.JSON(http.StatusOK, tableProduct)

}
