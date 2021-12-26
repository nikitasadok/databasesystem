package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikitasadok/database-system/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

type DatabaseHandler struct {
	DUsecase domain.DatabaseUsecase
}

func NewDatabaseHandler(r *gin.RouterGroup, us domain.DatabaseUsecase) {
	handler := &DatabaseHandler{
		DUsecase: us,
	}

	r.GET("/database/:id", handler.GetByID)
	r.PUT("/database/:id", handler.Update)
	r.DELETE("/database/:id", handler.Delete)
	r.POST("/database", handler.Create)
}

func (h *DatabaseHandler) Create(c *gin.Context) {
	var db domain.Database
	err := c.Bind(&db)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	err = h.DUsecase.Create(ctx, &db)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, db)
}

func (h *DatabaseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	db, err := h.DUsecase.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, db)
}

func (h *DatabaseHandler) Update(c *gin.Context) {
	// id := c.Param("id")

	var db domain.Database
	if err := c.Bind(&db); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()

	if err := h.DUsecase.Update(ctx, &db); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *DatabaseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	if err := h.DUsecase.Drop(ctx, id); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
