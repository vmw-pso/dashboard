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
	UpdatedAt     time.Time `json:"updatedAt"`
	Version       int64     `json:"version"`
	Closed        bool      `json:"closed"`
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
	qry := `
		INSERT INTO resource_requests(customer, start_date, end_date, hours_per_week, skills, opportunity_id, engagement_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, version, closed`

	args := []interface{}{
		rr.Customer,
		rr.StartDate,
		rr.EndDate,
		rr.HoursPerWeek,
		pq.Array(rr.Skills),
		rr.OpportunityID,
		rr.EngagementID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, qry, args...).Scan(&rr.ID, &rr.CreatedAt, &rr.UpdatedAt, &rr.Version, &rr.Closed)
}

func (m *ResourceRequestModel) Get(id int64) (*ResourceRequest, error) {
	if id < 1 {
		return nil, ErrNotFound
	}

	qry := `
		SELECT id, customer, start_date, end_data, hours_per_week, skills, opportunity_id, engagement_id, created_at, updated_at, version, closed
		FROM resource_requests
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rr ResourceRequest

	err := m.DB.QueryRowContext(ctx, qry, id).Scan(
		&rr.ID,
		&rr.Customer,
		&rr.StartDate,
		&rr.EndDate,
		&rr.HoursPerWeek,
		pq.Array(&rr.Skills),
		&rr.OpportunityID,
		&rr.EngagementID,
		&rr.CreatedAt,
		&rr.UpdatedAt,
		&rr.Version,
		&rr.Closed,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &rr, nil
}

func (m *ResourceRequestModel) Update(rr *ResourceRequest) error {
	qry := `
		UPDATE resource_requests
		SET customer=$1, start_date=$2, end_date=$3, hours_per_week=$4, skills=$5, updated_at=$6, version=version+1, closed=$7
		WHERE id=$8 AND updated_at=$9
		RETURNING updated_at`

	args := []interface{}{
		rr.Customer,
		rr.StartDate,
		rr.EndDate,
		rr.HoursPerWeek,
		pq.Array(rr.Skills),
		time.Now(),
		rr.Closed,
		rr.ID,
		rr.UpdatedAt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, qry, args...).Scan(&rr.UpdatedAt)
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

func (m *ResourceRequestModel) Delete(id int64) error {
	if id < 1 {
		return ErrNotFound
	}

	qry := `
		DELETE FROM resource_requests
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

func (m *ResourceRequestModel) GetAll(customer string, skills []string, closed bool, filters Filters) ([]*ResourceRequest, Metadata, error) {
	qry := fmt.Sprintf(`
		SELECT count(*) OVER(), id, customer, start_date, end_date, hours_per_week, skills, opportunity_id, engagement_id, created_at, updated_at, version
		FROM resource_requests
		WHERE (to_tsvector('simple', customer) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (closed = $2 OR $2 = false)
		AND (skills @> $3 OR $3 = '{}')
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{customer, closed, pq.Array(skills), filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, qry, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	resourceRequests := []*ResourceRequest{}

	for rows.Next() {
		var rr ResourceRequest
		err := rows.Scan(
			&totalRecords,
			&rr.ID,
			&rr.Customer,
			&rr.StartDate,
			&rr.EndDate,
			&rr.HoursPerWeek,
			&rr.Skills,
			&rr.OpportunityID,
			&rr.EngagementID,
			&rr.CreatedAt,
			&rr.UpdatedAt,
			&rr.Version,
			&rr.Closed,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		resourceRequests = append(resourceRequests, &rr)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return resourceRequests, metadata, nil
}
