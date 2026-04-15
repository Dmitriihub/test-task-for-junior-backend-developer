package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	recJSON, err := marshalRecurrence(task.Recurrence)
	if err != nil {
		return nil, fmt.Errorf("marshal recurrence: %w", err)
	}

	const query = `
		INSERT INTO tasks (title, description, status, recurrence, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, status, recurrence, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Status, recJSON, task.CreatedAt, task.UpdatedAt,
	)
	return scanTask(row)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, recurrence, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	found, err := scanTask(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	recJSON, err := marshalRecurrence(task.Recurrence)
	if err != nil {
		return nil, fmt.Errorf("marshal recurrence: %w", err)
	}

	const query = `
		UPDATE tasks
		SET title       = $1,
		    description = $2,
		    status      = $3,
		    recurrence  = $4,
		    updated_at  = $5
		WHERE id = $6
		RETURNING id, title, description, status, recurrence, created_at, updated_at
	`

	updated, err := scanTask(r.pool.QueryRow(ctx, query,
		task.Title, task.Description, task.Status, recJSON, task.UpdatedAt, task.ID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, recurrence, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}

	return tasks, rows.Err()
}

// ---- helpers ----------------------------------------------------------------

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
		recJSON []byte
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&recJSON,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	if len(recJSON) > 0 && string(recJSON) != "null" {
		var rec taskdomain.Recurrence
		if err := json.Unmarshal(recJSON, &rec); err != nil {
			return nil, fmt.Errorf("unmarshal recurrence: %w", err)
		}
		task.Recurrence = &rec
	}

	return &task, nil
}

func marshalRecurrence(r *taskdomain.Recurrence) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}
