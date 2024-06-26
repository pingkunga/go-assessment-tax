package deductions

import (
	"errors"
	"fmt"

	"github.com/dustin/go-humanize"
	repo "github.com/pingkunga/assessment-tax/postgres"
)

type DeductionService struct {
	repo DeductionRepository
}

type DeductionRepository interface {
	GetDeductionConfig(pKey string) (float64, error)
	SetDeductionConfig(pKey string, newValue float64) error
	DeductionConfigs() ([]repo.DeductionConfig, error)
	DeductionConfigByType(deductionType string) (repo.DeductionConfig, error)
}

func NewService(repo DeductionRepository) *DeductionService {
	return &DeductionService{repo: repo}
}

const personalDeduction = "personal"
const kReceipt = "k-receipt"

func (h *DeductionService) SetPersonalDeduction(request DebuctionRequest) (PersonalDeductionResponse, error) {
	DeductionConfig, err := h.repo.DeductionConfigByType(personalDeduction)
	if err != nil {
		return PersonalDeductionResponse{}, errors.New("personal deduction not found")
	}

	if request.Amount < DeductionConfig.DeductionMin {
		return PersonalDeductionResponse{}, fmt.Errorf("personal deduction must be greater than %s", humanize.FormatFloat("#,###.##", DeductionConfig.DeductionMin))
	}

	if request.Amount > DeductionConfig.DeductionMax {
		return PersonalDeductionResponse{}, fmt.Errorf("personal deduction must be less than %s", humanize.FormatFloat("#,###.##", DeductionConfig.DeductionMax))
	}

	errSet := h.repo.SetDeductionConfig(personalDeduction, request.Amount)
	if errSet != nil {
		return PersonalDeductionResponse{}, errSet
	}

	return PersonalDeductionResponse{PersonalDeduction: request.Amount}, nil

}

func (h *DeductionService) SetKReceiptDeduction(request DebuctionRequest) (KReceiptResponse, error) {
	DeductionConfig, err := h.repo.DeductionConfigByType(kReceipt)
	if err != nil {
		return KReceiptResponse{}, errors.New("k-receipt deduction not found")
	}

	if request.Amount < DeductionConfig.DeductionMin {
		return KReceiptResponse{}, fmt.Errorf("k-receipt deduction must be greater than %s", humanize.FormatFloat("#,###.##", DeductionConfig.DeductionMin))
	}

	if request.Amount > DeductionConfig.DeductionMax {
		return KReceiptResponse{}, fmt.Errorf("k-receipt deduction must be less than %s", humanize.FormatFloat("#,###.##", DeductionConfig.DeductionMax))
	}

	errSet := h.repo.SetDeductionConfig(kReceipt, request.Amount)
	if errSet != nil {
		return KReceiptResponse{}, errSet
	}

	return KReceiptResponse{KReceipt: request.Amount}, nil

}

func (h *DeductionService) GetPersonalDeduction() (PersonalDeductionResponse, error) {
	DeductionConfig, err := h.repo.DeductionConfigByType(personalDeduction)
	if err != nil {
		return PersonalDeductionResponse{}, errors.New("personal deduction not found")
	}

	return PersonalDeductionResponse{PersonalDeduction: DeductionConfig.DeductionAmount}, nil
}

func (h *DeductionService) DeductionConfigs() ([]repo.DeductionConfig, error) {
	deductions, err := h.repo.DeductionConfigs()
	if err != nil {
		return nil, err
	}

	return deductions, nil
}
