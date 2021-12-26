package docs

import (
	"github.com/nikitasadok/database-system/domain"
)

// swagger:route POST /database/{database_id}/table/{table_id}/row row row-create
// Create row record for specified database and table.
// responses:
//   200: rowResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route PUT /database/{database_id}/table/{table_id}/row/{id} row row-edit
// Edit row record for specified database and table.
// responses:
//   200: rowResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route DELETE /database/{database_id}/table/{table_id}/row/{id} row row-delete
// Delete row record for specified database and table.
// responses:
//   200: databaseResponseDelete
//   400: databaseResponseError
//   422: databaseResponseError

// Valid table record
// swagger:response rowResponse
type rowResponseWrapper struct {
	// in:body
	Body domain.Row
}

// swagger:parameters row-create
type rowCreateRequestWrapper struct {
	// in:path
	// swagger:name database_id
	DatabaseID string `json:"database_id"`
	// in:path
	TableID string `json:"table_id"`
	// in:body
	Body struct {
		Cells []struct {
			Value      interface{}       `json:"value"`
			ColumnData domain.ColumnData `json:"column"`
		} `json:"cells"`
	}
}

// swagger:parameters row-edit row-get row-delete
type rowRenameRequestWrapper struct {
	// in: path
	// swagger:name id
	ID string `json:"id"`
	// in:path
	// swagger:name database_id
	DatabaseID string `json:"database_id"`
	// in:path
	TableID string `json:"table_id"`
	// in:body
	Body []struct {
		Value      interface{}       `json:"value"`
		ColumnData domain.ColumnData `json:"column"`
	}
}
