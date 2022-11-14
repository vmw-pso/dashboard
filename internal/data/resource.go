package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

type Resource struct {
	ID             int64    `json:"resourceId"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Position       string   `json:"position"`
	Clearance      string   `json:"clearance"`
	Specialties    []string `json:"specialities,omitempty"`
	Certifications []string `json:"certifications,omitempty"`
}

func ValidateFirstName(v *validator.Validator, firstName string) {
	v.Check(firstName != "", "firstName", "must be provided")
	v.Check(len(firstName) < 256, "firstName", "must not be more than 256 bytes")
}

func ValidateLastName(v *validator.Validator, lastName string) {
	v.Check(lastName != "", "lastName", "must be provided")
	v.Check(len(lastName) < 256, "lastName", "must not be more than 256 bytes")
}

func ValidatePosition(v *validator.Validator, position string) {
	positions := []string{"AC1, AC2, C, SC, CA, SCA"}
	v.Check(validator.PermittedValue(position, positions...), "position", "does not exist")
}
func ValidateClearance(v *validator.Validator, clearance string) {
	clearances := []string{"None", "Baseline", "NV1", "NV2", "TSPV"}
	v.Check(validator.PermittedValue(clearance, clearances...), "clearance", "does not exist")
}

func ValidateResource(v *validator.Validator, r *Resource) {
	ValidateFirstName(v, r.FirstName)
	ValidateLastName(v, r.LastName)
	ValidatePosition(v, r.Position)
	ValidateClearance(v, r.Clearance)
	v.Check(validator.Unique(r.Specialties), "specialties", "must not contain duplicate values")
	v.Check(validator.Unique(r.Certifications), "certification", "must not contain duplicate values")
}

type ResourceModel struct {
	DB *sql.DB
}

func (m *ResourceModel) Insert(r *Resource) error {
	qry := `
		INSERT INTO resources
		(first_name, last_name, position, clearance, specialites, certification)
		VALUES ($1, $2, (SELECT id FROM positions WHERE title = $3), (SELECT id FROM clearances WHERE description = $4), $5, $6)
		RETURNING id`

	args := []interface{}{r.FirstName, r.LastName, r.Position, r.Clearance, r.Specialties, r.Certifications}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, qry, args...).Scan(&r.ID)
}
