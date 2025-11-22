package database

import (
	"context"
	"database/sql"
	"keep-your-house-clean/internal/domain"
	"time"

	"github.com/lib/pq"
)

type ComplimentRepository struct {
	db *sql.DB
}

func NewComplimentRepository(db *sql.DB) domain.ComplimentRepository {
	return &ComplimentRepository{db: db}
}

func (r *ComplimentRepository) Create(ctx context.Context, compliment *domain.Compliment) error {
	query := `
		INSERT INTO compliments (
			title, description, points, from_user_id, to_user_id,
			tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at, viewed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		compliment.Title,
		compliment.Description,
		compliment.Points,
		compliment.FromUserID,
		compliment.ToUserID,
		compliment.TenantID,
		compliment.CreatedAt,
		compliment.CreatedById,
		compliment.UpdatedAt,
		compliment.UpdatedById,
		compliment.DeletedAt,
		compliment.ViewedAt,
	).Scan(&compliment.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *ComplimentRepository) GetByID(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error) {
	query := `
		SELECT id, title, description, points, from_user_id, to_user_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at, viewed_at
		FROM compliments
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	var compliment domain.Compliment
	err := r.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&compliment.ID,
		&compliment.Title,
		&compliment.Description,
		&compliment.Points,
		&compliment.FromUserID,
		&compliment.ToUserID,
		&compliment.TenantID,
		&compliment.CreatedAt,
		&compliment.CreatedById,
		&compliment.UpdatedAt,
		&compliment.UpdatedById,
		&compliment.DeletedAt,
		&compliment.ViewedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &compliment, nil
}

func (r *ComplimentRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.Compliment, error) {
	query := `
		SELECT id, title, description, points, from_user_id, to_user_id,
		       tenant_id, created_at, created_by_id, updated_at, updated_by_id, deleted_at, viewed_at
		FROM compliments
		WHERE deleted_at IS NULL AND tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compliments []domain.Compliment
	for rows.Next() {
		var compliment domain.Compliment
		err := rows.Scan(
			&compliment.ID,
			&compliment.Title,
			&compliment.Description,
			&compliment.Points,
			&compliment.FromUserID,
			&compliment.ToUserID,
			&compliment.TenantID,
			&compliment.CreatedAt,
			&compliment.CreatedById,
			&compliment.UpdatedAt,
			&compliment.UpdatedById,
			&compliment.DeletedAt,
			&compliment.ViewedAt,
		)
		if err != nil {
			return nil, err
		}
		compliments = append(compliments, compliment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return compliments, nil
}

func (r *ComplimentRepository) GetLastReceivedByUser(ctx context.Context, userID int64, tenantID int64) (*domain.ComplimentWithUser, error) {
	query := `
		SELECT c.id, c.title, c.description, c.points, c.from_user_id, c.to_user_id,
		       c.tenant_id, c.created_at, c.created_by_id, c.updated_at, c.updated_by_id, c.deleted_at, c.viewed_at,
		       u.name as from_user_name
		FROM compliments c
		LEFT JOIN users u ON c.from_user_id = u.id AND u.deleted_at IS NULL
		WHERE c.deleted_at IS NULL AND c.tenant_id = $1 AND c.to_user_id = $2
		ORDER BY c.created_at DESC
		LIMIT 1
	`

	var compliment domain.ComplimentWithUser
	err := r.db.QueryRowContext(ctx, query, tenantID, userID).Scan(
		&compliment.ID,
		&compliment.Title,
		&compliment.Description,
		&compliment.Points,
		&compliment.FromUserID,
		&compliment.ToUserID,
		&compliment.TenantID,
		&compliment.CreatedAt,
		&compliment.CreatedById,
		&compliment.UpdatedAt,
		&compliment.UpdatedById,
		&compliment.DeletedAt,
		&compliment.ViewedAt,
		&compliment.FromUserName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &compliment, nil
}

func (r *ComplimentRepository) GetUserComplimentsHistory(ctx context.Context, userID int64, tenantID int64) ([]domain.ComplimentWithUser, error) {
	query := `
		SELECT c.id, c.title, c.description, c.points, c.from_user_id, c.to_user_id,
		       c.tenant_id, c.created_at, c.created_by_id, c.updated_at, c.updated_by_id, c.deleted_at, c.viewed_at,
		       CASE 
		         WHEN c.from_user_id = $2 THEN to_user.name
		         WHEN c.to_user_id = $2 THEN from_user.name
		       END as from_user_name
		FROM compliments c
		LEFT JOIN users from_user ON c.from_user_id = from_user.id AND from_user.deleted_at IS NULL
		LEFT JOIN users to_user ON c.to_user_id = to_user.id AND to_user.deleted_at IS NULL
		WHERE c.deleted_at IS NULL 
			AND c.tenant_id = $1 
			AND (c.from_user_id = $2 OR c.to_user_id = $2)
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compliments []domain.ComplimentWithUser
	for rows.Next() {
		var compliment domain.ComplimentWithUser
		err := rows.Scan(
			&compliment.ID,
			&compliment.Title,
			&compliment.Description,
			&compliment.Points,
			&compliment.FromUserID,
			&compliment.ToUserID,
			&compliment.TenantID,
			&compliment.CreatedAt,
			&compliment.CreatedById,
			&compliment.UpdatedAt,
			&compliment.UpdatedById,
			&compliment.DeletedAt,
			&compliment.ViewedAt,
			&compliment.FromUserName,
		)
		if err != nil {
			return nil, err
		}
		compliments = append(compliments, compliment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return compliments, nil
}

func (r *ComplimentRepository) GetUnviewedReceivedCompliments(ctx context.Context, userID int64, tenantID int64) ([]domain.ComplimentWithUser, error) {
	query := `
		SELECT c.id, c.title, c.description, c.points, c.from_user_id, c.to_user_id,
		       c.tenant_id, c.created_at, c.created_by_id, c.updated_at, c.updated_by_id, c.deleted_at, c.viewed_at,
		       u.name as from_user_name
		FROM compliments c
		LEFT JOIN users u ON c.from_user_id = u.id AND u.deleted_at IS NULL
		WHERE c.deleted_at IS NULL 
			AND c.tenant_id = $1 
			AND c.to_user_id = $2
			AND c.viewed_at IS NULL
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compliments []domain.ComplimentWithUser
	for rows.Next() {
		var compliment domain.ComplimentWithUser
		err := rows.Scan(
			&compliment.ID,
			&compliment.Title,
			&compliment.Description,
			&compliment.Points,
			&compliment.FromUserID,
			&compliment.ToUserID,
			&compliment.TenantID,
			&compliment.CreatedAt,
			&compliment.CreatedById,
			&compliment.UpdatedAt,
			&compliment.UpdatedById,
			&compliment.DeletedAt,
			&compliment.ViewedAt,
			&compliment.FromUserName,
		)
		if err != nil {
			return nil, err
		}
		compliments = append(compliments, compliment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return compliments, nil
}

func (r *ComplimentRepository) MarkAsViewed(ctx context.Context, ids []int64, userID int64, tenantID int64) error {
	if len(ids) == 0 {
		return nil
	}

	now := time.Now()
	query := `
		UPDATE compliments 
		SET viewed_at = $1 
		WHERE id = ANY($2::bigint[]) 
			AND tenant_id = $3 
			AND to_user_id = $4 
			AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, now, pq.Array(ids), tenantID, userID)
	return err
}

func (r *ComplimentRepository) Delete(ctx context.Context, id int64, tenantID int64) error {
	now := time.Now()
	query := `UPDATE compliments SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3 AND deleted_at IS NULL`

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

