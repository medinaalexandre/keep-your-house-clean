package database

import (
	"context"
	"database/sql"
	"keep-your-house-clean/internal/domain"
	"time"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (
			title, description, points, status, scheduled_to, scheduled_by_id,
			frequency_value, frequency_unit, completed, completed_by_id,
			tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Points,
		task.Status,
		task.ScheduledTo,
		task.ScheduledById,
		task.FrequencyValue,
		task.FrequencyUnit,
		task.Completed,
		task.CompletedById,
		task.TenantID,
		task.CreatedAt,
		task.CreatedById,
		task.UpdatedAt,
		task.UpdatedById,
		task.DeletedAt,
	).Scan(&task.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.Task, error) {
	query := `
		SELECT id, title, description, points, status, scheduled_to, scheduled_by_id,
		       frequency_value, frequency_unit, completed, completed_by_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Points,
			&task.Status,
			&task.ScheduledTo,
			&task.ScheduledById,
			&task.FrequencyValue,
			&task.FrequencyUnit,
			&task.Completed,
			&task.CompletedById,
			&task.TenantID,
			&task.CreatedAt,
			&task.CreatedById,
			&task.UpdatedAt,
			&task.UpdatedById,
			&task.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
	query := `
		SELECT id, title, description, points, status, scheduled_to, scheduled_by_id,
		       frequency_value, frequency_unit, completed, completed_by_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at
		FROM tasks
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	var task domain.Task
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Points,
		&task.Status,
		&task.ScheduledTo,
		&task.ScheduledById,
		&task.FrequencyValue,
		&task.FrequencyUnit,
		&task.Completed,
		&task.CompletedById,
		&task.TenantID,
		&task.CreatedAt,
		&task.CreatedById,
		&task.UpdatedAt,
		&task.UpdatedById,
		&task.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks SET
			title = $1,
			description = $2,
			points = $3,
			status = $4,
			scheduled_to = $5,
			scheduled_by_id = $6,
			frequency_value = $7,
			frequency_unit = $8,
			completed = $9,
			completed_by_id = $10,
			updated_at = $11,
			updated_by_id = $12
		WHERE id = $13 AND tenant_id = $14 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Points,
		task.Status,
		task.ScheduledTo,
		task.ScheduledById,
		task.FrequencyValue,
		task.FrequencyUnit,
		task.Completed,
		task.CompletedById,
		task.UpdatedAt,
		task.UpdatedById,
		task.ID,
		task.TenantID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64, tenantID int64) error {
	now := time.Now()
	query := `UPDATE tasks SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, now, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TaskRepository) GetUpcomingTasks(ctx context.Context, tenantID int64, limit int, offset int) ([]domain.Task, error) {
	query := `
		SELECT id, title, description, points, status, scheduled_to, scheduled_by_id,
		       frequency_value, frequency_unit, completed, completed_by_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL AND tenant_id = $1 AND completed = false
		ORDER BY scheduled_to ASC NULLS FIRST
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Points,
			&task.Status,
			&task.ScheduledTo,
			&task.ScheduledById,
			&task.FrequencyValue,
			&task.FrequencyUnit,
			&task.Completed,
			&task.CompletedById,
			&task.TenantID,
			&task.CreatedAt,
			&task.CreatedById,
			&task.UpdatedAt,
			&task.UpdatedById,
			&task.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetCompletedTasksHistory(ctx context.Context, tenantID int64, limit int) ([]domain.TaskWithUser, error) {
	query := `
		SELECT t.id, t.title, t.description, t.points, t.status, t.scheduled_to, t.scheduled_by_id,
		       t.frequency_value, t.frequency_unit, t.completed, t.completed_by_id,
		       t.tenant_id, t.created_at, t.created_by_id, t.updated_at, t.updated_by_id, t.deleted_at,
		       u.name as completed_by_name
		FROM tasks t
		LEFT JOIN users u ON t.completed_by_id = u.id AND u.deleted_at IS NULL
		WHERE t.deleted_at IS NULL AND t.tenant_id = $1 AND t.completed = true
		ORDER BY t.updated_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.TaskWithUser
	for rows.Next() {
		var task domain.TaskWithUser
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Points,
			&task.Status,
			&task.ScheduledTo,
			&task.ScheduledById,
			&task.FrequencyValue,
			&task.FrequencyUnit,
			&task.Completed,
			&task.CompletedById,
			&task.TenantID,
			&task.CreatedAt,
			&task.CreatedById,
			&task.UpdatedAt,
			&task.UpdatedById,
			&task.DeletedAt,
			&task.CompletedByName,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) FindTaskCreatedAfterCompletion(ctx context.Context, originalTask *domain.Task, completionTime time.Time) (*domain.Task, error) {
	query := `
		SELECT id, title, description, points, status, scheduled_to, scheduled_by_id,
		       frequency_value, frequency_unit, completed, completed_by_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL 
			AND tenant_id = $1 
			AND title = $2 
			AND description = $3
			AND completed = false
			AND created_at >= $4
			AND created_at <= $5
		ORDER BY created_at ASC
		LIMIT 1
	`

	timeWindowStart := completionTime.Add(-time.Minute)
	timeWindowEnd := completionTime.Add(time.Minute)

	var task domain.Task
	err := r.db.QueryRowContext(
		ctx,
		query,
		originalTask.TenantID,
		originalTask.Title,
		originalTask.Description,
		timeWindowStart,
		timeWindowEnd,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Points,
		&task.Status,
		&task.ScheduledTo,
		&task.ScheduledById,
		&task.FrequencyValue,
		&task.FrequencyUnit,
		&task.Completed,
		&task.CompletedById,
		&task.TenantID,
		&task.CreatedAt,
		&task.CreatedById,
		&task.UpdatedAt,
		&task.UpdatedById,
		&task.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}
