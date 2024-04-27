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
		expectedBody    string
		mockCalculateFn func(t *Tax) (float64, float64, error)
	}{
		{
			name:            "invalid request",
			requestBody:     []byte(`[]`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":"code=400, message=Unmarshal type error: expected=tax.Request, got=array, field=, offset=1, internal=json: cannot unmarshal array into Go value of type tax.Request"}`,
		},
		{
			name:            "valid request",
			requestBody:     []byte(`{"totalIncome": 50000, "wht": 5000, "allowances": [{"allowanceType": "donation", "amount": 1000}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, error) { return 1000, 0, nil },
			expectedStatus:  http.StatusOK,
			expectedTax:     1000,
			expectedBody:    `{"tax":1000}`,
		},
		{
			name:            "invalid totalIncome",
			requestBody:     []byte(`{"totalIncome": -100, "wht": 0}`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":[{"field":"TotalIncome","message":"the value of TotalIncome must be greater than or equal 0"}]}`,
		},
		{
			name:            "invalid WHT greater than totalIncome",
			requestBody:     []byte(`{"totalIncome": 5000, "wht": 6000}`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":[{"field":"Wht","message":"the value of Wht value must be lower than or equal value of field TotalIncome"}]}`,
		},
		{
			name:            "tax calculation error",
			requestBody:     []byte(`{"totalIncome": 50000}`),
			mockCalculateFn: func(t *Tax) (float64, float64, error) { return 0, 0, errors.New("calculation error") },
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":"calculation error"}`,
		},
		{
			name:            "return tax refund field",
			requestBody:     []byte(`{"totalIncome": 150000.0, "wht": 10000, "allowances": [{"allowanceType": "donation", "amount": 200000.0}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, error) { return 0, 10000, nil },
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"tax":0,"taxRefund":10000}`,
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
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
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
