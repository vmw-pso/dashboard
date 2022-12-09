package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Positions        PositionModel
	Clearances       ClearanceModel
	Resources        ResourceModel
	ResourceRequests ResourceRequestModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Positions:        PositionModel{DB: db},
		Clearances:       ClearanceModel{DB: db},
		Resources:        ResourceModel{DB: db},
		ResourceRequests: ResourceRequestModel{DB: db},
	}
}
