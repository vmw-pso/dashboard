package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

type Sex int64

const (
	Unknown Sex = iota
	Male
	Female
	Withheld
)

func (s Sex) String() string {
	switch s {
	case Male:
		return "Male"
	case Female:
		return "Female"
	case Withheld:
		return "Not Specified"
	default:
		return "Unknown"
	}
}

type Resource struct {
	ID             int64    `json:"resourceId"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Position       string   `json:"position"`
	Clearance      string   `json:"clearance"`
	Specialties    []string `json:"specialties,omitempty"`
	Certifications []string `json:"certifications,omitempty"`
	Active         bool     `json:"active"`
	Sex            string   `json:"sex"`
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
	v.Check(validator.PermittedValue(clearance, clearances...), "clearance", "must be one of ('None', 'Baseline', 'NV1', 'NV2')")
}

func ValidateSex(v *validator.Validator, sex string) {
	sexes := []string{"Unknown", "Male", "Female", "Not Specified"}
	v.Check(validator.PermittedValue(sex, sexes...), "sex", "must be one of ('Unknown', 'Male', 'Female', 'Not Specified')")
}

func ValidateResource(v *validator.Validator, r Resource) {
	ValidateID(v, int(r.ID))
	ValidateFirstName(v, r.FirstName)
	ValidateLastName(v, r.LastName)
	ValidatePosition(v, r.Position)
	ValidateClearance(v, r.Clearance)
	ValidateSex(v, r.Sex)
	v.Check(validator.Unique(r.Specialties), "specialties", "must not contain duplicate values")
	v.Check(validator.Unique(r.Certifications), "certification", "must not contain duplicate values")
}

type ResourceModel struct {
	DB *sql.DB
}

func (m *ResourceModel) Insert(r *Resource) error {
	qry := `
		INSERT INTO resources
		(id, first_name, last_name, position_id, clearance_id, specialties, certifications, active, sex)
		VALUES ($1, $2, $3, (SELECT id FROM positions WHERE title = $4), (SELECT id FROM clearances WHERE description = $5), $6, $7, $8, $9)
		RETURNING id`

	args := []interface{}{r.ID, r.FirstName, r.LastName, r.Position, r.Clearance, pq.Array(r.Specialties), pq.Array(r.Certifications), r.Active, r.Sex}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, qry, args...).Scan(&r.ID)
}

func (m *ResourceModel) Get(id int64) (*Resource, error) {
	if id < 1 {
		return nil, ErrNotFound
	}

	qry := `
		SELECT resources.id, resources.first_name, resources.last_name, positions.title, clearances.description, resources.specialties, resources.certifications, resources.active, resources.sex
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
		&r.Sex,
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
		SET first_name = $1, last_name = $2, position_id = (SELECT id FROM positions WHERE title = $3), clearance_id = (SELECT id FROM clearances WHERE description = $4), specialties = $5, certifications = $6, active = $7, sex = $8
		WHERE id = $9`

	args := []interface{}{
		r.FirstName,
		r.LastName,
		r.Position,
		r.Clearance,
		pq.Array(r.Specialties),
		pq.Array(r.Certifications),
		r.Active,
		r.Sex,
		r.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, qry, args...).Scan()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
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

func (m *ResourceModel) GetAll(specialties []string, certifications []string, active bool, filters Filters) ([]*Resource, Metadata, error) {
	qry := fmt.Sprintf(`
		SELECT count(*) OVER(), resources.id, resources.first_name, resources.last_name, positions.title, clearances.description, resources.specialties, resources.certifications, resources.active, resources.sex
		FROM ((resources
			INNER JOIN positions ON positions.id = resources.position_id)
			INNER JOIN clearances ON clearances.id = resources.clearance_id)
		WHERE (specialties @> $1 OR $1 = '{}')
		AND (certifications @> $2 OR $2 = '{}')
		AND (active = $3 OR $3 = true)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{pq.Array(specialties), pq.Array(certifications), active, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, qry, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	resources := []*Resource{}

	for rows.Next() {
		var resource Resource
		err := rows.Scan(
			&totalRecords,
			&resource.ID,
			&resource.FirstName,
			&resource.LastName,
			&resource.Position,
			&resource.Clearance,
			pq.Array(&resource.Specialties),
			pq.Array(&resource.Certifications),
			&resource.Active,
			&resource.Sex,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return resources, metadata, nil
}
