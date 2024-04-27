package tax

import (
	"github.com/Atvit/assessment-tax/internals/setting"
	"github.com/Atvit/assessment-tax/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType" validate:"omitempty,oneof=donation k-receipt"`
	Amount        float64 `json:"amount" validate:"omitempty,gte=0"`
}

type Request struct {
	TotalIncome float64     `json:"totalIncome" validate:"required,gte=0"`
	Wht         float64     `json:"wht" validate:"omitempty,gte=0,ltefield=TotalIncome"`
	Allowances  []Allowance `json:"allowances" validate:"dive"`
}

type Response struct {
	Tax       float64    `json:"tax"`
	TaxLevel  []TaxLevel `json:"taxLevel,omitempty"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
}

type Handler interface {
	CalculateTax(c echo.Context) error
}

type handler struct {
	logger      *zap.Logger
	validate    *validator.Validate
	settingRepo setting.Repository
}

func NewHandler(
	logger *zap.Logger,
	validate *validator.Validate,
	settingRepo setting.Repository,
) Handler {
	return handler{
		logger:      logger,
		validate:    validate,
		settingRepo: settingRepo,
	}
}

func (h handler) CalculateTax(c echo.Context) error {
	var req Request

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

	allowanceSetting, err := h.settingRepo.Get()
	if err != nil {
		h.logger.Error("get allowance setting failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	var taxAllowances []TaxAllowance
	for _, allowances := range req.Allowances {
		taxAllowances = append(taxAllowances, TaxAllowance{
			AllowanceType: allowances.AllowanceType,
			Amount:        allowances.Amount,
		})
	}
	taxAmount, refundAmount, taxLevels, err := Calculate(&Tax{
		Income:     req.TotalIncome,
		Wht:        req.Wht,
		Allowances: taxAllowances,
		AllowanceSetting: TaxAllowanceSetting{
			Personal: allowanceSetting.Personal,
			KReceipt: allowanceSetting.KReceipt,
		},
	})
	if err != nil {
		h.logger.Error("tax calculation failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, utils.ErrResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Tax:       taxAmount,
		TaxLevel:  taxLevels,
		TaxRefund: refundAmount,
	})
}
