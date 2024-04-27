package setting

import (
	"bytes"
	"database/sql"
	"github.com/Atvit/assessment-tax/internals/models"
	mockSetting "github.com/Atvit/assessment-tax/mocks/setting"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdatePersonalDeduction(t *testing.T) {
	type testcase struct {
		requestBody    []byte
		expectedStatus int
		expectedBody   string
	}

	e := echo.New()
	logger := zap.NewNop()
	validate := validator.New()
	repo := new(mockSetting.Repository)

	h := &handler{
		logger:     logger,
		validate:   validate,
		repository: repo,
	}

	t.Run("valid request", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 70000.0}`),
			expectedStatus: 200,
			expectedBody:   `{"personalDeduction":70000}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo.On("Update", mock.AnythingOfType("uint"), mock.AnythingOfType("*models.DeductionConfigEntity")).Return(&models.DeductionConfig{ID: 1, Personal: 70000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`[]`),
			expectedStatus: 400,
			expectedBody:   `{"error":"code=400, message=Unmarshal type error: expected=setting.PersonalDeductionRequest, got=array, field=, offset=1, internal=json: cannot unmarshal array into Go value of type setting.PersonalDeductionRequest"}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("amount greter than 100000", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 1000000.0}`),
			expectedStatus: 400,
			expectedBody:   `{"error":[{"field":"Amount","message":"the value of Amount must be less than or equal 100000"}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("amount less than 10000", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 100.0}`),
			expectedStatus: 400,
			expectedBody:   `{"error":[{"field":"Amount","message":"the value of Amount must be greater than or equal 10000"}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("amount equal 10000", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 10000.0}`),
			expectedStatus: 200,
			expectedBody:   `{"personalDeduction":10000}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo.On("Update", mock.AnythingOfType("uint"), mock.AnythingOfType("*models.DeductionConfigEntity")).Return(&models.DeductionConfig{ID: 1, Personal: 10000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("amount equal 100000", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 100000.0}`),
			expectedStatus: 200,
			expectedBody:   `{"personalDeduction":100000}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repo.On("Update", mock.AnythingOfType("uint"), mock.AnythingOfType("*models.DeductionConfigEntity")).Return(&models.DeductionConfig{ID: 1, Personal: 100000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("update db error", func(t *testing.T) {
		tc := testcase{
			requestBody:    []byte(`{"amount": 100000.0}`),
			expectedStatus: 500,
			expectedBody:   `{"error":"sql: no rows in result set"}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		errNoRows := sql.ErrNoRows
		repo.On("Update", mock.AnythingOfType("uint"), mock.AnythingOfType("*models.DeductionConfigEntity")).Return(nil, errNoRows).Once()

		if assert.NoError(t, h.UpdatePersonalDeduction(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		}
	})
}
