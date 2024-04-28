package tax

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/Atvit/assessment-tax/internals/models"
	mockSetting "github.com/Atvit/assessment-tax/mocks/setting"
	"github.com/Atvit/assessment-tax/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CalculateTax(t *testing.T) {
	type testcase struct {
		requestBody     []byte
		expectedStatus  int
		expectedBody    string
		mockCalculateFn func(t *Tax) (float64, float64, []TaxLevel, error)
	}

	e := echo.New()
	logger := zap.NewNop()
	validate := validator.New()

	t.Run("invalid request", func(t *testing.T) {
		tc := testcase{
			requestBody:     []byte(`[]`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":"code=400, message=Unmarshal type error: expected=tax.Request, got=array, field=, offset=1, internal=json: cannot unmarshal array into Go value of type tax.Request"}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("valid request", func(t *testing.T) {
		tc := testcase{
			requestBody:     []byte(`{"totalIncome": 500000.0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) { return 29000, 0, getMockTaxLevels(), nil },
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"tax":29000,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":0},{"level":"500,001-1,000,000","tax":0},{"level":"1,000,001-2,000,000","tax":0},{"level":"2,000,001 ขึ้นไป","tax":0}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("k-receipt allowance", func(t *testing.T) {
		tc := testcase{
			requestBody: []byte(`{"totalIncome": 500000.0,"wht": 0.0,"allowances": [{"allowanceType": "k-receipt","amount": 200000.0},{"allowanceType": "donation","amount": 100000.0}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) {
				return 14000, 0, getMockTaxLevels(TaxLevel{level2, utils.ToPointer(14000.0)}), nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"tax":14000,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":14000},{"level":"500,001-1,000,000","tax":0},{"level":"1,000,001-2,000,000","tax":0},{"level":"2,000,001 ขึ้นไป","tax":0}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("invalid totalIncome", func(t *testing.T) {
		tc := testcase{
			requestBody:     []byte(`{"totalIncome": -100, "wht": 0}`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":[{"field":"TotalIncome","message":"the value of TotalIncome must be greater than or equal 0"}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("invalid WHT greater than totalIncome", func(t *testing.T) {
		tc := testcase{
			requestBody:     []byte(`{"totalIncome": 5000, "wht": 6000}`),
			mockCalculateFn: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":[{"field":"Wht","message":"the value of Wht value must be lower than or equal value of field TotalIncome"}]}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("tax calculation error", func(t *testing.T) {
		tc := testcase{
			requestBody:     []byte(`{"totalIncome": 50000}`),
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) { return 0, 0, nil, errors.New("calculation error") },
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    `{"error":"calculation error"}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("return tax refund field", func(t *testing.T) {
		tc := testcase{
			requestBody: []byte(`{"totalIncome": 150000.0, "wht": 10000.0, "allowances": [{"allowanceType": "donation", "amount": 200000.0}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) {
				return 0, 10000, getMockTaxLevels(
					TaxLevel{level2, utils.ToPointer(35000.0)},
					TaxLevel{level3, utils.ToPointer(75000.0)},
					TaxLevel{level4, utils.ToPointer(68000.0)},
				), nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"tax":0,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":35000},{"level":"500,001-1,000,000","tax":75000},{"level":"1,000,001-2,000,000","tax":68000},{"level":"2,000,001 ขึ้นไป","tax":0}],"taxRefund":10000}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})

	t.Run("get tax setting failed", func(t *testing.T) {
		tc := testcase{
			requestBody: []byte(`{"totalIncome": 150000.0, "wht": 10000.0, "allowances": [{"allowanceType": "donation", "amount": 200000.0}]}`),
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) {
				return 0, 10000, getMockTaxLevels(), nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"sql: no rows in result set"}`,
		}

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewReader(tc.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		errNoRows := sql.ErrNoRows
		settingRepo.On("Get").Return(nil, errNoRows).Once()

		if assert.NoError(t, h.CalculateTax(c)) {
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, rec.Body.String(), tc.expectedBody)
		}
	})
}

func TestHandler_UploadCSV(t *testing.T) {
	type testcase struct {
		fileContent     string
		mockReadError   error
		expectedStatus  int
		expectedBody    string
		mockCalculateFn func(t *Tax) (float64, float64, []TaxLevel, error)
	}

	e := echo.New()
	logger := zap.NewNop()
	validate := validator.New()

	t.Run("valid CSV file", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000",
			mockReadError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"taxes":[{"totalIncome":500000,"tax":29000},{"totalIncome":600000,"tax":0,"taxRefund":2000},{"totalIncome":750000,"tax":11250}]}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("empty CSV file", func(t *testing.T) {
		tc := testcase{
			fileContent:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"empty csv file given"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("CSV file with no records", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n",
			mockReadError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"empty csv file given"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("upload file error", func(t *testing.T) {
		tc := testcase{
			fileContent:    "",
			mockReadError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"multipart: NextPart: EOF"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("get tax setting failed", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"sql: no rows in result set"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		errNoRows := sql.ErrNoRows
		settingRepo.On("Get").Return(nil, errNoRows).Once()

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("negative donation", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n500000,0,0\n600000,40000,-20000\n750000,50000,15000",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":[{"field":"Donation","message":"the value of Donation must be greater than or equal 0"}]}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}

		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("negative income", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n-500000,0,0\n600000,40000,20000\n750000,50000,15000",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":[{"field":"TotalIncome", "message":"the value of TotalIncome must be greater than or equal 0"}]}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}

		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("wht higher than income", func(t *testing.T) {
		tc := testcase{
			fileContent:    "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,950000,15000",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":[{"field":"Wht", "message":"the value of Wht value must be lower than or equal value of field TotalIncome"}]}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}

		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("calculation error", func(t *testing.T) {
		tc := testcase{
			fileContent:     "totalIncome,wht,donation\n500000,0,0\n600000,40000,-20000\n750000,50000,15000",
			expectedStatus:  http.StatusBadRequest,
			mockCalculateFn: func(t *Tax) (float64, float64, []TaxLevel, error) { return 0, 0, nil, errors.New("calculation error") },
			expectedBody:    `{"error":"calculation error"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		settingRepo := new(mockSetting.Repository)

		originalCalculate := Calculate
		Calculate = tc.mockCalculateFn
		defer func() { Calculate = originalCalculate }()

		settingRepo.On("Get").Return(&models.DeductionConfig{ID: 1, Personal: 60000, KReceipt: 50000}, nil).Once()

		h := &handler{
			logger:      logger,
			validate:    validate,
			settingRepo: settingRepo,
		}

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}

		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

	t.Run("incorrect filename", func(t *testing.T) {
		tc := testcase{
			fileContent:    "",
			mockReadError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"http: no such file"}`,
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "taxes.csv")
		part.Write([]byte(tc.fileContent))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &handler{
			logger:   logger,
			validate: validate,
		}

		if err := h.UploadCSV(c); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectedStatus, rec.Code)
		assert.JSONEq(t, tc.expectedBody, rec.Body.String())
	})

}
