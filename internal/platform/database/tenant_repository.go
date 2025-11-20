package database

import (
	"context"
	"database/sql"
	"keep-your-house-clean/internal/domain"
	"time"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) domain.TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (name, domain, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		tenant.Name,
		tenant.Domain,
		tenant.Status,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	).Scan(&tenant.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	query := `
		SELECT id, name, domain, status, created_at, updated_at, deleted_at
		FROM tenants
		WHERE id = $1 AND deleted_at IS NULL
	`

	var tenant domain.Tenant
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Domain,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
		&tenant.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *TenantRepository) GetByDomain(ctx context.Context, domainParam string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, domain, status, created_at, updated_at, deleted_at
		FROM tenants
		WHERE domain = $1 AND deleted_at IS NULL
	`

	var tenant domain.Tenant
	err := r.db.QueryRowContext(ctx, query, domainParam).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Domain,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
		&tenant.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *TenantRepository) FetchAll(ctx context.Context) ([]domain.Tenant, error) {
	query := `
		SELECT id, name, domain, status, created_at, updated_at, deleted_at
		FROM tenants
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		var tenant domain.Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Domain,
			&tenant.Status,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
			&tenant.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants SET
			name = $1,
			domain = $2,
			status = $3,
			updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		tenant.Name,
		tenant.Domain,
		tenant.Status,
		tenant.UpdatedAt,
		tenant.ID,
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

func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	now := time.Now()
	query := `UPDATE tenants SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, now, id)
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
