package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/vmw-pso/delivery-dashboard/back-end/internal/validator"
)

type Position struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func ValidateTitle(v *validator.Validator, title string) {
	v.Check(title != "", "title", "must be provided")
	v.Check(len(title) <= 256, "title", "must not be more than 256 bytes")
}

type PositionModel struct {
	DB *sql.DB
}

func (m *PositionModel) Insert(p *Position) error {
	qry := `
		INSERT INTO positions (title)
		VALUES ($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, qry, p.Title).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *PositionModel) Get(id int64) (*Position, error) {
	qry := `
		SELECT title
		FROM positions
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, qry, id)

	p := Position{ID: id}

	if err := row.Scan(&p.Title); err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PositionModel) Update(p Position) error {
	qry := `
		UPDATE positions
		SET title = $1
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.QueryContext(ctx, qry, p.Title, p.ID)
	return err
}

func (m *PositionModel) Delete(id int64) error {
	qry := `
		DELETE FROM positions
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

func (m *PositionModel) GetAll() ([]*Position, error) {
	qry := `
		SELECT id, title
		FROM positions
		ORDER BY title`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*Position

	for rows.Next() {
		var p Position
		err := rows.Scan(
			&p.ID,
			&p.Title,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, &p)
	}

	return positions, nil
}
