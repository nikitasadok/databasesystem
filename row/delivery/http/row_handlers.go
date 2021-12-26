package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nikitasadok/database-system/domain"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type RowHandler struct {
	TUsecase domain.TableUsecase
	RUsecase domain.RowUsecase
}

func NewRowHandler(r *gin.RouterGroup, tu domain.TableUsecase, ru domain.RowUsecase) {
	handler := &RowHandler{
		TUsecase: tu,
		RUsecase: ru,
	}

	r.POST("/database/:dbID/table/:tableID/row", handler.Create)
	r.PUT("/database/:id/table/:tableID/row/:rowID", handler.Update)
	r.DELETE("/database/:id/table/:tableID/row/:rowID", handler.Delete)

}

func (rh *RowHandler) Create(c *gin.Context) {
	// "/database/:id/table/:tableID/row"
	var toRec map[string]*json.RawMessage
	var r domain.Row
	if err := c.Bind(&toRec); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
		return
	}

	gg := []struct {
		Value      interface{}       `json:"value"`
		ColumnData domain.ColumnData `json:"column"`
	}{}

	if err := json.Unmarshal(*toRec["cells"], &gg); err != nil {
		c.JSON(http.StatusBadRequest, "Error unmarshaling")
		return
	}

	for i := range gg {
		convertedValue, err := doCast(gg[i], gg[i].ColumnData)
		if err != nil {
			c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
			return
		}
		r.Cells = append(r.Cells, convertedValue)
	}

	fmt.Printf("%+v\n", r.Cells)

	tableID := c.Param("tableID")
	r.TableIDString = tableID

	err := rh.RUsecase.Create(c.Request.Context(), &r)
	if err != nil {
		fmt.Println("error: ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println("here")
	c.JSON(http.StatusOK, "ok")
}

func (rh *RowHandler) Update(c *gin.Context) {
	rowID := c.Param("rowID")
	tableID := c.Param("tableID")
	var toRec map[string]*json.RawMessage
	var r domain.Row
	if err := c.Bind(&toRec); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
		return
	}

	gg := []struct {
		Value      interface{}       `json:"value"`
		ColumnData domain.ColumnData `json:"column"`
	}{}

	if err := json.Unmarshal(*toRec["cells"], &gg); err != nil {
		c.JSON(http.StatusBadRequest, "cannot unmarshal cells")
		return
	}

	for i := range gg {
		convertedValue, err := doCast(gg[i], gg[i].ColumnData)
		if err != nil  {
			c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
			return
		}
		r.Cells = append(r.Cells, convertedValue)
	}

	r.TableIDString = tableID

	fmt.Printf("%+v\n", r.Cells)

	err := rh.RUsecase.Update(c.Request.Context(), &r, rowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "uar")
}

func (rh *RowHandler) Delete(c *gin.Context) {
	rowID := c.Param("rowID")
	tableID := c.Param("tableID")

	if err := rh.RUsecase.Delete(c.Request.Context(), rowID, tableID); err != nil {
		c.JSON(http.StatusBadRequest, "bad")
		return
	}

	c.JSON(http.StatusOK, "ok")
}

func doCast(cell struct {
	Value      interface{}       `json:"value"`
	ColumnData domain.ColumnData `json:"column"`
}, colData domain.ColumnData) (domain.Cell, error) {
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
		return nil, fmt.Errorf("unknown datatype %s", colData.Datatype)
	}

	err := cellUpd.SetValue(cell.Value)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert value %v to datatype %s", cell.Value, colData.Datatype)
	}
	cellUpd.SetColumnData(colData)
	fmt.Printf("cellUpd: %+v\n", cellUpd)

	return cellUpd, nil
}