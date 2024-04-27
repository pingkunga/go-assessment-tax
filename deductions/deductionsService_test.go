package deductions

import (
	"testing"

	repo "github.com/pingkunga/assessment-tax/postgres"
	"github.com/stretchr/testify/assert"
)

func TestDeductionsService_SetPersonalDeduction(t *testing.T) {

	t.Run("given peronal dection 1000 but min 10,000 / max 100,000, should return error", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 1000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 1000,
		}

		// Call the SetPersonalDeduction method
		response, err := service.SetPersonalDeduction(request)

		assert.Error(t, err)
		assert.EqualValues(t, "personal deduction must be greater than 10,000.00", err.Error())
		assert.EqualValues(t, 0, response.PersonalDeduction)
	})

	t.Run("given peronal dection 100,001 but min 10,000 / max 100,000, should return error", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100001,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 100001,
		}

		// Call the SetPersonalDeduction method
		response, err := service.SetPersonalDeduction(request)

		assert.Error(t, err)
		assert.EqualValues(t, "personal deduction must be less than 100,000.00", err.Error())
		assert.EqualValues(t, 0, response.PersonalDeduction)
	})

	t.Run("given peronal dection 10,000 but min 10,000 / max 100,000, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 10000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 10000,
		}

		// Call the SetPersonalDeduction method
		response, err := service.SetPersonalDeduction(request)

		assert.NoError(t, err)
		assert.EqualValues(t, 10000, response.PersonalDeduction)
	})

	t.Run("given peronal dection 100,000 but min 10,000 / max 100,000, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 100000,
		}

		// Call the SetPersonalDeduction method
		response, err := service.SetPersonalDeduction(request)

		assert.NoError(t, err)
		assert.EqualValues(t, 100000, response.PersonalDeduction)
	})
}

func TestDeductionsService_SetKReceiptDeduction(t *testing.T) {

	t.Run("given k-receipt dection 70000 but min 0 / max 100,000, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 70000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 70000,
		}

		// Call the SetKReceiptDeduction method
		response, err := service.SetKReceiptDeduction(request)

		assert.NoError(t, err)
		assert.EqualValues(t, 70000, response.KReceipt)

	})

	t.Run("given k-receipt dection 0 but min 0 / max 100,000, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 0,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 0,
		}

		// Call the SetKReceiptDeduction method
		response, err := service.SetKReceiptDeduction(request)

		assert.NoError(t, err)
		assert.EqualValues(t, 0, response.KReceipt)
	})

	t.Run("given k-receipt dection 100,000 but min 0 / max 100,000, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 100000,
		}

		// Call the SetKReceiptDeduction method
		response, err := service.SetKReceiptDeduction(request)

		assert.NoError(t, err)
		assert.EqualValues(t, 100000, response.KReceipt)
	})

	t.Run("given k-receipt dection 100,001 but min 0 / max 100,000, should return error", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100001,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 100001,
		}

		// Call the SetKReceiptDeduction method
		response, err := service.SetKReceiptDeduction(request)

		assert.Error(t, err)
		assert.EqualValues(t, "k-receipt deduction must be less than 100,000.00", err.Error())
		assert.EqualValues(t, 0, response.KReceipt)
	})

	t.Run("given k-receipt dection 100,001 but min 0 / max 100,000, should return error", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100001,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Create a new DeductionRequest
		request := DebuctionRequest{
			Amount: 100001,
		}

		// Call the SetKReceiptDeduction method
		response, err := service.SetKReceiptDeduction(request)

		assert.Error(t, err)
		assert.EqualValues(t, "k-receipt deduction must be less than 100,000.00", err.Error())
		assert.EqualValues(t, 0, response.KReceipt)
	})
}

func TestDeductionsService_GetDeductionConfigs(t *testing.T) {
	t.Run("given no deduction config, should return error", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			DeductionConfigsStub: []repo.DeductionConfig{},
			errStub:              assert.AnError,
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Call the DeductionConfigs method
		deductions, err := service.DeductionConfigs()

		assert.Error(t, err)
		assert.Nil(t, deductions)
	})

	t.Run("given deduction config, should return success", func(t *testing.T) {
		// Create a new StubRepository
		repo := &StubRepository{
			DeductionConfigsStub: []repo.DeductionConfig{
				{
					DeductionType: "personal",
					DeductionMin:  10000,
					DeductionMax:  100000,
				},
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)

		// Call the DeductionConfigs method
		deductions, err := service.DeductionConfigs()

		assert.NoError(t, err)
		assert.NotNil(t, deductions)
	})
}

// Mock Repository
type StubRepository struct {
	GetDeductionStub     float64
	SetDeductionStub     float64
	DeductionConfigsStub []repo.DeductionConfig
	DeductionConfigStub  repo.DeductionConfig
	errStub              error
}

func (r *StubRepository) GetDeductionConfig(pKey string) (float64, error) {
	return r.GetDeductionStub, nil
}

func (r *StubRepository) SetDeductionConfig(pKey string, newValue float64) error {
	r.SetDeductionStub = newValue
	return nil
}

func (r *StubRepository) DeductionConfigs() ([]repo.DeductionConfig, error) {
	return r.DeductionConfigsStub, r.errStub
}

func (r *StubRepository) DeductionConfigByType(deductionType string) (repo.DeductionConfig, error) {
	return r.DeductionConfigStub, nil
}
