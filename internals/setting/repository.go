package setting

import (
	"database/sql"
	"github.com/Atvit/assessment-tax/internals/models"
	"time"
)

type Repository interface {
	Get() (*models.DeductionConfig, error)
	UpdatePersonalDeduction(id uint, value float64) (*models.DeductionConfig, error)
	UpdateKReceiptDeduction(id uint, value float64) (*models.DeductionConfig, error)
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

func (r repository) UpdatePersonalDeduction(id uint, value float64) (*models.DeductionConfig, error) {
	query := `UPDATE tax_deduction_configs SET personal = $1, updated_at = $2 WHERE id = $3 RETURNING *`
	row := r.db.QueryRow(query, value, time.Now(), id)

	var result models.DeductionConfig
	err := row.Scan(&result.ID, &result.Personal, &result.KReceipt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r repository) UpdateKReceiptDeduction(id uint, value float64) (*models.DeductionConfig, error) {
	query := `UPDATE tax_deduction_configs SET kreceipt = $1, updated_at = $2 WHERE id = $3 RETURNING *`
	row := r.db.QueryRow(query, value, time.Now(), id)

	var result models.DeductionConfig
	err := row.Scan(&result.ID, &result.Personal, &result.KReceipt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
