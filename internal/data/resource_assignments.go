package data

import "time"

type ResourceAssignment struct {
	ResourceRequestID int64     `json:"resourceRequestId"`
	ResourceID        int64     `json:"resourceId"`
	HoursPerWeek      int64     `json:"hoursPerWeek"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Version           int64     `json:"version"`
	Completed         bool      `json:"completed"`
}
