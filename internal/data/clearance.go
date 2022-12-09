package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

type Clearance struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

func ValidateDescription(v *validator.Validator, description string) {
	v.Check(description != "", "description", "must be provided")
	v.Check(len(description) <= 256, "title", "must not be more than 256 bytes")
}

type ClearanceModel struct {
	DB *sql.DB
}

func (m *ClearanceModel) Insert(c *Clearance) error {
	qry := `
		INSERT INTO clearances (description)
		VALUES ($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, qry, c.Description).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ClearanceModel) Get(id int64) (*Clearance, error) {
	qry := `
		SELECT description
		FROM clearances
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, qry, id)

	c := Clearance{ID: id}

	if err := row.Scan(&c.Description); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *ClearanceModel) Update(c Clearance) error {
	qry := `
		UPDATE clearances
		SET description = $1
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.QueryContext(ctx, qry, c.Description, c.ID)
	return err
}

func (m *ClearanceModel) Delete(id int64) error {
	qry := `
		DELETE FROM clearances
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (m *ClearanceModel) GetAll() ([]*Clearance, error) {
	qry := `
		SELECT id, description
		FROM clearances
		ORDER BY description`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clearances []*Clearance

	for rows.Next() {
		var c Clearance
		err := rows.Scan(
			&c.ID,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}
		clearances = append(clearances, &c)
	}

	return clearances, nil
}
