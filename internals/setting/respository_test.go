package setting

import (
	"database/sql"
	"errors"
	"github.com/Atvit/assessment-tax/internals/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var mockDBErr = errors.New("could not open database connection")

func TestRepository_Get(t *testing.T) {
	mockRow := models.DeductionConfig{
		ID:        1,
		Personal:  60000.00,
		KReceipt:  70000.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	r := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "personal", "kreceipt", "created_at", "updated_at"}).
			AddRow(mockRow.ID, mockRow.Personal, mockRow.KReceipt, mockRow.CreatedAt, mockRow.UpdatedAt)
		mock.ExpectQuery(getStmt).WillReturnRows(rows)

		result, err := r.Get()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, 60000.00, result.Personal)
		assert.Equal(t, 70000.00, result.KReceipt)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(getStmt).WillReturnError(sql.ErrNoRows)

		result, err := r.Get()

		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("error scan rows", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "personal", "kreceipt", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil)
		mock.ExpectQuery(getStmt).WillReturnRows(rows)

		mock.ExpectQuery(getStmt).WillReturnRows(rows)

		result, err := r.Get()

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}

func TestRepository_UpdatePersonalDeduction(t *testing.T) {
	mockRow := models.DeductionConfig{
		ID:        1,
		Personal:  60000.00,
		KReceipt:  70000.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	r := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(updatePersonalStmt).
			WithArgs(60000.00, sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "personal", "kreceipt", "created_at", "updated_at"}).
				AddRow(mockRow.ID, mockRow.Personal, mockRow.KReceipt, mockRow.CreatedAt, mockRow.UpdatedAt))

		result, err := r.UpdatePersonalDeduction(1, 60000.00)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 60000.00, result.Personal)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(updatePersonalStmt).WillReturnError(mockDBErr)

		result, err := r.UpdatePersonalDeduction(1, 60000.00)

		assert.Nil(t, result)
		assert.Error(t, mockDBErr, err)
	})
}

func TestRepository_UpdateKReceiptDeduction(t *testing.T) {
	mockRow := models.DeductionConfig{
		ID:        1,
		Personal:  60000.00,
		KReceipt:  70000.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	r := NewRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(updateKReceiptStmt).
			WithArgs(70000.00, sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "personal", "kreceipt", "created_at", "updated_at"}).
				AddRow(1, mockRow.Personal, mockRow.KReceipt, mockRow.CreatedAt, mockRow.UpdatedAt))

		result, err := r.UpdateKReceiptDeduction(1, 70000.00)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 70000.00, result.KReceipt)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(updateKReceiptStmt).WillReturnError(mockDBErr)

		result, err := r.UpdateKReceiptDeduction(1, 70000.00)

		assert.Nil(t, result)
		assert.Equal(t, mockDBErr, err)
	})
}
