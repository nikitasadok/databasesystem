package domain

type Cell interface {
	// swagger:name value
	GetValue() interface{}
	SetValue(val interface{}) error
	// swagger:ignore
	GetColumnName() string
	SetColumnData(col ColumnData)
	// swagger:name datatype
	GetDatatype() string
}

func GetConcreteCellForType(datatype string) Cell {
	switch datatype {
	case INTEGER:
		return &IntegerCell{}
	case CHAR:
		return &CharCell{}
	case STRING:
		return &StringCell{}
	case REAL:
		return &RealCell{}
	}

	return nil
}
