package docs

import (
"github.com/nikitasadok/database-system/domain"
)

// swagger:route POST /database/{database_id}/table table table-create
// Create table record for specified database.
// responses:
//   200: tableResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route GET /database/{database_id}/table/{id} table table-get
// Get table record by ID for specified database.
// responses:
//   200: tableResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route GET /database/{database_id}/tablesProduct table table-product
// Get table record by ID for specified database.
// responses:
//   200: tableProductResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route PUT /database/{database_id}/table/{id} table table-rename
// Rename table record by ID for specified database.
// responses:
//   200: tableResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route DELETE /database/{database_id}/table/{id} table table-delete
// Delete table record by ID for specified database.
// responses:
//   200: databaseResponseDelete
//   400: databaseResponseError
//   422: databaseResponseError

// Valid table record
// swagger:response tableResponse
type tableResponseWrapper struct {
	// in:body
	Body domain.Table
}

// swagger:parameters table-create
type tableCreateRequestWrapper struct {
	// in:path
	// swagger:name database_id
	DatabaseID string `json:"database_id"`
	// in:body
	Body struct {
		Name string `json:"name"`
		Columns []struct{
			Datatype string `json:"datatype"`
			Name string `json:"name"`
		} `json:"columns"`
	}
}

// swagger:parameters table-rename table-get table-delete
type tableRenameRequestWrapper struct {
	// in: path
	// swagger:name id
	ID string `json:"id"`
	// in:path
	// swagger:name database_id
	DatabaseID string `json:"database_id"`
	// in:body
	Body struct {
		Name string `json:"name"`
	}
}

// swagger:parameters table-get table-delete
type tableGetRequestWrapper struct {
	// in:path
	// min: 12 max: 12
	ID string `json:"id"`
	// in:path
	DatabaseID string `json:"database_id"`
}

// swagger:response tableProductResponse
type tableProductGetResponseWrapper struct {
	//in: body
	Body domain.TablesProduct
}

// swagger:parameters table-product
type tableProductGetRequestWrapper struct {
	//in:path
	DatabaseID string `json:"database_id"`
	//in:query
	T1 string `json:"t1_id"`
	T2 string `json:"t2_id"`
}
