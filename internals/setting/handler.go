package setting

import (
	"github.com/Atvit/assessment-tax/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type PersonalDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required,lte=100000,gte=10000"`
}

type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KReceiptDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required,lte=100000,gt=0"`
}

type KReceiptDeductionResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

type Handler interface {
	UpdatePersonalDeduction(c echo.Context) error
	UpdateKReceiptDeduction(c echo.Context) error
}

type handler struct {
	logger     *zap.Logger
	validate   *validator.Validate
	repository Repository
}

func NewHandler(logger *zap.Logger, validate *validator.Validate, repository Repository) Handler {
	return &handler{
		logger:     logger,
		validate:   validate,
		repository: repository,
	}
}

func (h handler) UpdatePersonalDeduction(c echo.Context) error {
	var req PersonalDeductionRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Error("binding request failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("validate request body failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, utils.ErrResponse{
			Error: utils.GetValidateErrMsg(err),
		})
	}

	result, err := h.repository.UpdatePersonalDeduction(1, req.Amount)
	if err != nil {
		h.logger.Error("update personal allowance failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, PersonalDeductionResponse{
		PersonalDeduction: result.Personal,
	})
}

func (h handler) UpdateKReceiptDeduction(c echo.Context) error {
	var req KReceiptDeductionRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Error("binding request failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("validate request body failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, utils.ErrResponse{
			Error: utils.GetValidateErrMsg(err),
		})
	}

	result, err := h.repository.UpdateKReceiptDeduction(1, req.Amount)
	if err != nil {
		h.logger.Error("update k-receipt allowance failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, KReceiptDeductionResponse{
		KReceipt: result.KReceipt,
	})
}
