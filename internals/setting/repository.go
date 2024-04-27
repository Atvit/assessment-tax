package setting

import (
	"database/sql"
	"github.com/Atvit/assessment-tax/internals/models"
	"time"
)

type Repository interface {
	Get() (*models.DeductionConfig, error)
	Update(id uint, entity *models.DeductionConfigEntity) (*models.DeductionConfig, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) Get() (*models.DeductionConfig, error) {
	rows, err := r.db.Query("SELECT * FROM tax_deduction_configs ORDER BY updated_at DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var config models.DeductionConfig
	for rows.Next() {
		err := rows.Scan(&config.ID, &config.Personal, &config.KReceipt, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func (r repository) Update(id uint, entity *models.DeductionConfigEntity) (*models.DeductionConfig, error) {
	query := `UPDATE tax_deduction_configs SET personal = $1, kreceipt = $2, updated_at = $3 WHERE id = $4 RETURNING *`
	row := r.db.QueryRow(query, entity.Personal, entity.KReceipt, time.Now(), id)

	var result models.DeductionConfig
	err := row.Scan(&result.ID, &result.Personal, &result.KReceipt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
