package docs

import (
	"github.com/nikitasadok/database-system/domain"
	"github.com/nikitasadok/database-system/row/delivery/http"
)

// swagger:route POST /database database database-create
// Create database record.
// responses:
//   200: databaseResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route GET /database/{id} database database-get
// Get database record by ID.
// responses:
//   200: databaseResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route PUT /database/{id} database database-rename
// Rename database record by ID.
// responses:
//   200: databaseResponse
//   400: databaseResponseError
//   422: databaseResponseError

// swagger:route DELETE /database/{id} database database-delete
// Delete database record by ID.
// responses:
//   200: databaseResponseDelete
//   400: databaseResponseError
//   422: databaseResponseError

// Valid database record
// swagger:response databaseResponse
type databaseResponseWrapper struct {
	// in:body
	Body domain.Database
}

// swagger:parameters database-create
type databaseCreateRequestWrapper struct {
	// in:body
	Body struct {
		Name string `json:"name"`
	}
}

// swagger:parameters database-rename
type databaseRenameRequestWrapper struct {
	// in: path
	ID string `json:"id"`
	// in:body
	Body struct {
		Name string `json:"name"`
	}
}

// swagger:parameters database-get database-delete
type databaseGetRequestWrapper struct {
	// in:path
	// min: 12 max: 12
	ID string `json:"id"`
}

// swagger:response databaseResponseError
type reponseErrorWrapper struct {
	// in:body
	Body http.ResponseError
}

// swagger:response databaseResponseDelete
type responseDeleteWrapper struct {
	// in:body
	Body struct {
		Status string `json:"status"`
	}
}