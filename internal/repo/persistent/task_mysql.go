package persistent

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/mysql"
)

// TaskRepo -.
type TaskRepo struct {
	*mysql.MySQL
}

// NewTaskRepo -.
func NewTaskRepo(db *mysql.MySQL) *TaskRepo {
	return &TaskRepo{db}
}

// Store -.
func (r *TaskRepo) Store(ctx context.Context, task *entity.Task) error {
	query, args, err := r.Builder.
		Insert("tasks").
		Columns("id, user_id, title, description, status, created_at, updated_at").
		Values(task.ID, task.UserID, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *TaskRepo) GetByID(ctx context.Context, userID, taskID string) (entity.Task, error) {
	query, args, err := r.Builder.
		Select("id, user_id, title, description, status, created_at, updated_at").
		From("tasks").
		Where(sq.Eq{"id": taskID}).
		ToSql()
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskRepo - GetByID - r.Builder: %w", err)
	}

	var task entity.Task

	err = r.Pool.QueryRow(ctx, query, args...).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if errors.Is(err, dbsql.ErrNoRows) {
			return entity.Task{}, entity.ErrTaskNotFound
		}

		return entity.Task{}, fmt.Errorf("TaskRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	if task.UserID != userID {
		return entity.Task{}, entity.ErrTaskForbidden
	}

	return task, nil
}

// List -.
func (r *TaskRepo) List(ctx context.Context, userID string, filter repo.TaskFilter) ([]entity.Task, int, error) {
	countBuilder := r.Builder.
		Select("COUNT(*)").
		From("tasks").
		Where(sq.Eq{"user_id": userID})

	if filter.Status != nil {
		countBuilder = countBuilder.Where(sq.Eq{"status": *filter.Status})
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - countBuilder: %w", err)
	}

	var total int

	err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - count query: %w", err)
	}

	dataBuilder := r.Builder.
		Select("id, user_id, title, description, status, created_at, updated_at").
		From("tasks").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset)

	if filter.Status != nil {
		dataBuilder = dataBuilder.Where(sq.Eq{"status": *filter.Status})
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - dataBuilder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("TaskRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0, filter.Limit)

	for rows.Next() {
		var t entity.Task

		err = rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("TaskRepo - List - rows.Scan: %w", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, total, nil
}

// Update -.
func (r *TaskRepo) Update(ctx context.Context, task *entity.Task) error {
	query, args, err := r.Builder.
		Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("status", task.Status).
		Set("updated_at", task.UpdatedAt).
		Where(sq.Eq{"id": task.ID, "user_id": task.UserID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - r.Pool.Exec: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("TaskRepo - Update - result.RowsAffected: %w", err)
	}

	if affected == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}

// Delete -.
func (r *TaskRepo) Delete(ctx context.Context, userID, taskID string) error {
	query, args, err := r.Builder.
		Delete("tasks").
		Where(sq.Eq{"id": taskID, "user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - r.Pool.Exec: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("TaskRepo - Delete - result.RowsAffected: %w", err)
	}

	if affected == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}
