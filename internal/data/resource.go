package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
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
	Active         bool     `json:"active"`
}

func ValidateID(v *validator.Validator, id int) {
	v.Check(id != 0, "id", "must be provided")
	v.Check(id > 0, "id", "cannot be a negative number")
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
	positions := []string{"Associate Consultant I", "Associate Consultant II", "Consultant", "Senior Consultant", "Staff Consultant", "Consulting Architect", "Staff Consulting Architect"}
	v.Check(validator.PermittedValue(position, positions...), "position", "does not exist")
}
func ValidateClearance(v *validator.Validator, clearance string) {
	clearances := []string{"None", "Baseline", "NV1", "NV2", "TSPV"}
	v.Check(validator.PermittedValue(clearance, clearances...), "clearance", "does not exist")
}

func ValidateResource(v *validator.Validator, r *Resource) {
	ValidateID(v, int(r.ID))
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
		(id, first_name, last_name, position_id, clearance_id, specialties, certifications, active)
		VALUES ($1, $2, $3, (SELECT id FROM positions WHERE title = $4), (SELECT id FROM clearances WHERE description = $5), $6, $7, $8)
		RETURNING id`

	args := []interface{}{r.ID, r.FirstName, r.LastName, r.Position, r.Clearance, pq.Array(r.Specialties), pq.Array(r.Certifications), r.Active}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, qry, args...).Scan(&r.ID)
}

func (m *ResourceModel) Get(id int64) (*Resource, error) {
	if id < 1 {
		return nil, ErrNotFound
	}

	qry := `
		SELECT resources.id, resources.first_name, resources.last_name, positions.title, clearances.description, resources.specialties, resources.certifications, resources.active
		FROM ((resources
		INNER JOIN positions ON positions.id = resources.position_id)
		INNER JOIN clearances ON clearances.id = resources.clearance_id)
		WHERE resources.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var r Resource

	err := m.DB.QueryRowContext(ctx, qry, id).Scan(
		&r.ID,
		&r.FirstName,
		&r.LastName,
		&r.Position,
		&r.Clearance,
		pq.Array(&r.Specialties),
		pq.Array(&r.Certifications),
		&r.Active,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &r, nil
}

func (m *ResourceModel) Update(r *Resource) error {
	qry := `
		UPDATE resources
		SET first_name = $1, last_name = $2, position_id = (SELECT id FROM positions WHERE title = $3), clearance_id = (SELECT id FROM clearances WHERE description = $4), specialties = $5, certificates = $6, active = $7
		WHERE id = $8`

	args := []interface{}{
		r.FirstName,
		r.LastName,
		r.Position,
		r.Clearance,
		pq.Array(r.Specialties),
		pq.Array(r.Certifications),
		r.Active,
		r.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, qry, args...).Scan()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditCOnflict
		default:
			return err
		}
	}
	return nil
}

func (m *ResourceModel) Delete(id int64) error {
	if id < 1 {
		return ErrNotFound
	}

	qry := `
		DELETE FROM resources
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// TODO: Complete this with filters passed as query parameters
func (m *ResourceModel) GetAll(specialty int, clearance int, filters Filters) ([]*Resource, error) {
	return nil, errors.New("not implemented")
}
