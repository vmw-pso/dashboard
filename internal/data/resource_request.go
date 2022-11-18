package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

type ResourceRequest struct {
	ID            int64     `json:"id"`
	Customer      string    `json:"customer"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	HoursPerWeek  int64     `json:"hoursPerWeek"`
	Skills        []string  `json:"skills"`
	OpportunityID string    `json:"projectID,omitempty"`
	EngagementID  string    `json:"engagementID,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatesAt     time.Time `json:"updatedAt"`
	Version       int64     `json:"version"`
}

func ValidateCustomer(v *validator.Validator, customer string) {
	v.Check(customer != "", "customer", "must be provided")
}

func ValidateSkills(v *validator.Validator, skills []string) {
	v.Check(len(skills) > 0, "skills", "at least one must be provided")
}

func ValidateResourceRequest(v *validator.Validator, rr ResourceRequest) {
	ValidateCustomer(v, rr.Customer)
	ValidateSkills(v, rr.Skills)
	v.Check(validator.Unique(rr.Skills), "skills", "must not contain duplicate values")
}

type ResourceRequestModel struct {
	DB *sql.DB
}

func (m *ResourceRequestModel) Insert(rr *ResourceRequest) error {
	return errors.New("not implemented")
}

func (m *ResourceRequestModel) Get(id int64) (*ResourceRequest, error) {
	return nil, errors.New("not implemented")
}

func (m *ResourceRequestModel) GetAll(customer string, skills []string, filters Filters) ([]*ResourceRequest, Metadata, error) {
	return nil, Metadata{}, errors.New("not implemented")
}

func (m *ResourceRequestModel) Update(rr *ResourceRequest) error {
	return errors.New("not implemented")
}

func (m *ResourceRequestModel) Delete(id int64) error {
	return errors.New("not implemented")
}
