package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTax(t *testing.T) {
	var tests = []struct {
		name            string
		requestBody     []byte
		expectedStatus  int
		expectedTax     float64
		mockCalculateFn func(t *Tax) (float64, float64, error)
	}{
		{
			name:            "invalid request",
			requestBody:     []byte(`[]`),
			expectedStatus:  http.StatusBadRequest,
			mockCalculateFn: nil,
		},
		{
			name:            "valid request",
			requestBody:     []byte(`{"totalIncome": 50000, "wht": 5000, "allowances": [{"allowanceType": "donation", "amount": 1000}]}`),
			expectedStatus:  http.StatusOK,
			expectedTax:     1000,
			mockCalculateFn: func(t *Tax) (float64, float64, error) { return 1000, 0, nil },
		},
		{
			name:            "invalid totalIncome",
			requestBody:     []byte(`{"totalIncome": -100, "wht": 0}`),
			expectedStatus:  http.StatusBadRequest,
			mockCalculateFn: nil,
		},
		{
			name:            "invalid WHT greater than totalIncome",
			requestBody:     []byte(`{"totalIncome": 5000, "wht": 6000}`),
			expectedStatus:  http.StatusBadRequest,
			mockCalculateFn: nil,
		},
		{
			name:            "tax calculation error",
			requestBody:     []byte(`{"totalIncome": 50000}`),
			expectedStatus:  http.StatusBadRequest,
			mockCalculateFn: func(t *Tax) (float64, float64, error) { return 0, 0, errors.New("calculation error") },
		},
	}

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tax/calculations", bytes.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zap.NewNop()
			validate := validator.New()
			h := &handler{
				logger:   logger,
				validate: validate,
			}

			originalCalculate := Calculate
			if tt.mockCalculateFn != nil {
				Calculate = tt.mockCalculateFn
			}
			defer func() { Calculate = originalCalculate }()

			if assert.NoError(t, h.CalculateTax(c)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if rec.Code == http.StatusOK {
					var resp Response
					if err := json.Unmarshal(rec.Body.Bytes(), &resp); err == nil {
						assert.Equal(t, tt.expectedTax, resp.Tax)
					}
				}
			}
		})
	}
}
