package deductions

import (
	repo "github.com/pingkunga/assessment-tax/postgres"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{repo: repo}
}

func (h *Service) SetPersonalDeduction(request DebuctionRequest) (PersonalDeductionResponse, error) {
	err := h.repo.SetDeductionConfig("personal", request.Amount)
	if err != nil {
		return PersonalDeductionResponse{}, err
	}

	return PersonalDeductionResponse{PersonalDeduction: request.Amount}, nil

}

func (h *Service) SetKPlustDeduction(request DebuctionRequest) (KReceiptResponse, error) {
	err := h.repo.SetDeductionConfig("k-receipt", request.Amount)
	if err != nil {
		return KReceiptResponse{}, err
	}

	return KReceiptResponse{KReceipt: request.Amount}, nil

}

func (h *Service) AllowanceConfigs() ([]repo.DeductionConfig, error) {
	allowances, err := h.repo.AllowanceConfigs()
	if err != nil {
		return nil, err
	}

	return allowances, nil
}
