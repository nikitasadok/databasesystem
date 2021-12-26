package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/tools/godoc/util"
)

type (
	IntegerCell struct {
		Value  int         `json:"value" bson:"value"`
		Column *ColumnData `json:"column" bson:"column"`
	}

	CharCell struct {
		Value  rune        `json:"value" bson:"value"`
		Column *ColumnData `json:"column" bson:"column"`
	}

	StringCell struct {
		Value  string      `json:"value" bson:"value"`
		Column *ColumnData `json:"column" bson:"column"`
	}

	RealCell struct {
		Value  float64     `json:"value" bson:"value"`
		Column *ColumnData `json:"column" bson:"column"`
	}

	TextFileCell struct {
		Value    []byte      `json:"value" bson:"value"`
		Column   *ColumnData `json:"column" bson:"column"`
		Filename string      `json:"filename" bson:"filename"`
	}

	IntegerIntervalCell struct {
		Value  IntegerIntervalCellValue
		Column *ColumnData `json:"column" bson:"column"`
	}

	IntegerIntervalCellValue struct {
		Start int `json:"start" bson:"start"`
		End   int `json:"end" bson:"end"`
	}
)

func (i *IntegerIntervalCell) GetValue() interface{} {
	return i.Value
}

func (i *IntegerIntervalCell) SetValue(val interface{}) error {
	conv, ok := val.(map[string]interface{})
	if !ok {
		start, end, err := i.parseDBValueFromDB(val)
		if err != nil {
			return err
		}

		i.Value = IntegerIntervalCellValue{
			Start: start,
			End:   end,
		}

		return nil
	}
	start, ok := conv["start"]
	if !ok {
		fmt.Println("here1")
		return fmt.Errorf("cannot convert Value %v to type IntegerIntervalCellValue", val)
	}
	end, ok := conv["end"]
	if !ok {
		fmt.Println("here2")
		return fmt.Errorf("cannot convert Value %v to type IntegerIntervalCellValue", val)
	}

	startConv, _ := start.(float64)
	endConv, _ := end.(float64)

	startConvInt := int(startConv)
	endConvInt := int(endConv)

	if startConvInt > endConvInt {
		return fmt.Errorf("start of the inverval (%d) is greater than the end (%d)", startConvInt, endConvInt)
	}

	i.Value = IntegerIntervalCellValue{
		Start: startConvInt,
		End:   endConvInt,
	}
	return nil
}

func (i *IntegerIntervalCell) parseDBValueFromDB(val interface{}) (start, end int, err error) {
	conv, ok := val.(primitive.D)
	if !ok {
		return 0, 0, errors.New("cannot cast")
	}

	var startTemp, endTemp int32
	for i := range conv {
		if conv[i].Key == "start" {
			startTemp = conv[i].Value.(int32)
		} else {
			endTemp = conv[i].Value.(int32)
		}
	}

	start = int(startTemp)
	end = int(endTemp)

	return start, end, nil
}

func (i *IntegerIntervalCell) GetColumnName() string {
	return i.Column.Name
}

func (i *IntegerIntervalCell) SetColumnData(col ColumnData) {
	i.Column = &col
}

func (i *IntegerIntervalCell) GetDatatype() string {
	return i.Column.Datatype
}

func (t *TextFileCell) GetValue() interface{} {
	return t.Value
}

func (t *TextFileCell) SetValue(val interface{}) error {
	conv, ok := val.(string)
	if !ok {
		return t.parseFromDB(val)
	}

	decodedString, err := base64.StdEncoding.DecodeString(conv)
	if err != nil {
		return fmt.Errorf("invalid base64 supplied")
	}

	if !util.IsText(decodedString) {
		return errors.New("provided file is not valid UTF-8 file")
	}

	t.Value = []byte(conv)
	return nil
}

func (t *TextFileCell) parseFromDB(val interface{}) error {
	conv, ok := val.(primitive.Binary)
	if !ok {
		return fmt.Errorf("cannot convert Value %v into primitive.Binary", val)
	}

	t.Value = conv.Data
	return nil
}

func (t *TextFileCell) GetColumnName() string {
	return t.Column.Name
}

func (t *TextFileCell) SetColumnData(col ColumnData) {
	t.Column = &col
}

func (t *TextFileCell) GetDatatype() string {
	return t.Column.Datatype
}

func (c *RealCell) GetColumnName() string {
	return c.Column.Name
}

func (c *RealCell) SetColumnData(col ColumnData) {
	c.Column = &col
}

func (t *RealCell) GetDatatype() string {
	return t.Column.Datatype
}

func (c *StringCell) GetColumnName() string {
	return c.Column.Name
}

func (c *StringCell) SetColumnData(col ColumnData) {
	c.Column = &col
}

func (c *StringCell) GetDatatype() string {
	return c.Column.Datatype
}

func (c *CharCell) GetColumnName() string {
	return c.Column.Name
}

func (c *CharCell) SetColumnData(col ColumnData) {
	c.Column = &col
}

func (t *CharCell) GetDatatype() string {
	return t.Column.Datatype
}

func NewIntegerCell(val int, colData *ColumnData) IntegerCell {
	return IntegerCell{Value: val, Column: colData}
}

func NewCharCell(val rune, colData *ColumnData) CharCell {
	return CharCell{Value: val, Column: colData}
}

func NewStringCell(val string, colData *ColumnData) StringCell {
	return StringCell{Value: val, Column: colData}
}

func NewRealCell(val float64, colData *ColumnData) RealCell {
	return RealCell{Value: val, Column: colData}
}

func (c *IntegerCell) GetValue() interface{} {
	return c.Value
}

func (c *IntegerCell) SetValue(val interface{}) error {
	conv, ok := val.(float64)
	if !ok {
		convToInt, ok := val.(int32)
		if !ok {
			return fmt.Errorf("cannot convert Value %v to type float type %T", val, val)
		}
		c.Value = int(convToInt)
		return nil
	}

	convInt := int(conv)

	c.Value = convInt
	return nil
}

func (c *IntegerCell) GetColumnName() string {
	return c.Column.Name
}

func (c *IntegerCell) SetColumnData(col ColumnData) {
	c.Column = &col
}

func (t *IntegerCell) GetDatatype() string {
	return t.Column.Datatype
}

func (c *CharCell) GetValue() interface{} {
	return c.Value
}

func (c *CharCell) SetValue(val interface{}) error {
	conv, ok := val.(string)
	if !ok {
		return c.handleDBUpload(val)
		// return fmt.Errorf("cannot convert Value %v to type byte", val)
	}

	c.Value = ([]rune(conv))[0]
	return nil
}

func (c *CharCell) handleDBUpload(val interface{}) error {
	conv, ok := val.(rune)
	if !ok {
		return fmt.Errorf("cannot convert Value %v to type byte when loading from DB", val)
	}
	
	c.Value = conv
	return nil
}

func (c *StringCell) GetValue() interface{} {
	return c.Value
}

func (c *StringCell) SetValue(val interface{}) error {
	conv, ok := val.(string)
	if !ok {
		return fmt.Errorf("cannot convert Value %v to type string", val)
	}

	c.Value = conv
	return nil
}

func (c *RealCell) GetValue() interface{} {
	return c.Value
}

func (c *RealCell) SetValue(val interface{}) error {
	conv, ok := val.(float64)
	if !ok {
		return fmt.Errorf("cannot convert Value %v to type float64", val)
	}

	c.Value = conv
	return nil
}
