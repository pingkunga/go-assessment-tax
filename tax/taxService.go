package tax

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	deductions "github.com/pingkunga/assessment-tax/deductions"
	repo "github.com/pingkunga/assessment-tax/postgres"
)

type TaxService struct {
	deductsvc deductions.IDeductionService
}

func NewService(debuctionSVC deductions.IDeductionService) *TaxService {
	return &TaxService{deductsvc: debuctionSVC}
}

func (s *TaxService) CalculateTax(tax TaxRequest) (TaxResponse, error) {

	netIncome, err := s.calculateNetIncome(tax)
	if err != nil {
		return TaxResponse{}, err
	}

	var taxAmount float64
	var taxLevels []TaxLevel
	for _, step := range TaxStepList() {
		stepTaxAmount := 0.0
		if netIncome <= 0 {
			stepTaxAmount = 0
		} else if netIncome > step.MaxIncome {
			stepTaxAmount = round2digits(step.MaxIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		} else {
			stepTaxAmount = round2digits(netIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		}

		taxLevels = append(taxLevels, TaxLevel{Level: step.FriendlyMessage, Tax: stepTaxAmount})
		taxAmount = taxAmount + stepTaxAmount
	}

	taxAmount = taxAmount - tax.WHT

	if taxAmount < 0 {
		return TaxResponse{Tax: 0, TaxRefund: taxAmount * -1, TaxLevels: taxLevels}, nil
	}
	return TaxResponse{Tax: taxAmount, TaxRefund: 0, TaxLevels: taxLevels}, nil
}

// Note
// Go Round https://yourbasic.org/golang/round-float-2-decimal-places/
func round2digits(value float64) float64 {
	return float64(int(value*100)) / 100
}

func (s *TaxService) calculateNetIncome(tax TaxRequest) (float64, error) {
	netIncome := tax.TotalIncome - s.PersonalDeduction()

	deductionls, err := s.deductsvc.DeductionConfigs()
	if err != nil {
		return 0, errors.New("error get deduction config")
	}

	for _, allowance := range tax.Allowances {
		netIncome = netIncome - allowanceAmountGuard(allowance, deductionls)
	}

	return netIncome, nil
}

func allowanceAmountGuard(allowance Allowance, deductionls []repo.DeductionConfig) float64 {

	for _, deduction := range deductionls {
		if allowance.AllowanceType == deduction.DeductionType && allowance.Amount > deduction.DeductionMax {
			return deduction.DeductionMax
		}
	}

	return allowance.Amount
}

const DEDUCTION_DONATION = "donation"
const DEDUCTION_DONATION_MAX = 100000

const DEDUCTION_K_RECEIPT = "k-receipt"
const DEDUCTION_K_RECEIPT_MAX = 100000

var ErrTotalIncomeLessThanZero = errors.New("TotalIncome is less than 0")
var ErrWHTLessThanZero = errors.New("WHT is less than 0")
var ErrWHTMoreThanTotalIncome = errors.New("WHT is more than TotalIncome")
var ErrAllowanceAmountLessThanZero = errors.New("Allowance amount is less than 0")
var ErrNotSupportAllowanceType = errors.New("Allowance type support: donation, k-receipt")

func ValidateTaxRequest(tax TaxRequest) (err error) {
	if tax.TotalIncome < 0 {
		err = errors.Join(err, ErrTotalIncomeLessThanZero)
	}

	if tax.WHT < 0 {
		err = errors.Join(err, ErrWHTLessThanZero)
	}

	if tax.TotalIncome < tax.WHT {
		err = errors.Join(err, ErrWHTMoreThanTotalIncome)
	}

	for _, allowance := range tax.Allowances {

		if allowance.AllowanceType != DEDUCTION_DONATION && allowance.AllowanceType != DEDUCTION_K_RECEIPT {
			err = errors.Join(err, ErrNotSupportAllowanceType)
		}

		if allowance.Amount < 0 {
			err = errors.Join(err, ErrAllowanceAmountLessThanZero)
		}
	}

	return
}

func (s *TaxService) PersonalDeduction() float64 {
	//return 60000

	//Get Personal Deduction from Deduction Service
	personalDeduction, err := s.deductsvc.GetPersonalDeduction()
	if err != nil {
		return 60000
	}
	return personalDeduction.PersonalDeduction
}

/*
การคำนวนภาษีตามขั้นบันใด
รายได้ 0 - 150,000 ได้รับการยกเว้น
150,001 - 500,000 อัตราภาษี 10%
500,001 - 1,000,000 อัตราภาษี 15%
1,000,001 - 2,000,000 อัตราภาษี 20%
มากกว่า 2,000,000 อัตราภาษี 35%
*/
func TaxStepList() []TaxStep {
	taxStep := []TaxStep{
		{seq: 1, MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
		{seq: 2, MinIncome: 150001, MaxIncome: 500000, TaxRate: 0.1},
		{seq: 3, MinIncome: 500001, MaxIncome: 1000000, TaxRate: 0.15},
		{seq: 4, MinIncome: 1000001, MaxIncome: 2000000, TaxRate: 0.2},
		{seq: 5, MinIncome: 2000001, MaxIncome: 999999999, TaxRate: 0.35},
	}

	//Update FriendlyMessage
	for i, step := range taxStep {
		if taxStep[i].MaxIncome == 999999999 {
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s ขึ้นไป", humanize.FormatFloat("#,###.", step.MinIncome))
		} else {
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s-%s", humanize.FormatFloat("#,###.", step.MinIncome), humanize.FormatFloat("#,###.", step.MaxIncome))
		}
	}

	return taxStep
}

func (s *TaxService) CalculateTaxBatch(taxs []TaxRequest) (TaxBatchsResponse, error) {
	var taxBatchs []TaxBatchResponse
	for _, tax := range taxs {
		taxResponse, err := s.CalculateTax(tax)
		if err != nil {
			return TaxBatchsResponse{}, err
		}
		taxBatchs = append(taxBatchs, TaxBatchResponse{TotalIncome: tax.TotalIncome, Tax: taxResponse.Tax, TaxRefund: taxResponse.TaxRefund})
	}
	return TaxBatchsResponse{Taxes: taxBatchs}, nil
}

func ImportTaxCSV(pathToFile string) (taxls []TaxRequest, errTax error) {

	if _, err := os.Stat(pathToFile); err != nil {
		return nil, errors.Join(errTax, fmt.Errorf("error Open CSV file at: %s", pathToFile))
	}

	csvFile, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	//parse csv file
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return nil, errors.Join(errTax, fmt.Errorf("error read CSV file at: %s", pathToFile))
	}

	//loop through each row
	var taxRequests []TaxRequest
	var header []string
	for i, row := range csvData {
		if i == 0 {
			header = row
			errHead := ValidateHeader(header)
			errTax = concatenateError(errTax, errHead)
			continue
		}

		totalIncome, err := ParseFloatForImport(header[0], row[0], i)
		errTax = concatenateError(errTax, err)

		wht, err := ParseFloatForImport(header[1], row[1], i)
		errTax = concatenateError(errTax, err)

		allowances, err := createAllowances(header, row, i)
		errTax = concatenateError(errTax, err)

		//Add Result
		if errTax == nil {
			TaxRequest := TaxRequest{
				TotalIncome: totalIncome,
				WHT:         wht,
				Allowances:  allowances,
			}
			taxRequests = append(taxRequests, TaxRequest)
		}
	}

	if errTax != nil {
		return nil, errTax
	}

	//validate each row
	return taxRequests, nil
}

func concatenateError(err error, err2 error) error {
	if err != nil {
		return errors.Join(err, err2)
	}
	return err2
}

func createAllowances(heading []string, row []string, rowId int) (allowances []Allowance, err error) {

	//donation
	allowanceDonation, err := createAllowance(DEDUCTION_DONATION, heading[2], row[2], rowId)
	if err != nil {
		return nil, err
	}
	allowances = append(allowances, allowanceDonation)
	//k-receipt
	// allowanceKReceipt, err := createAllowance(DEDUCTION_K_RECEIPT, row[3], rowId)
	// if err != nil {
	// 	return nil, err
	// }
	// allowances = append(allowances, allowanceKReceipt)

	return allowances, nil
}

func createAllowance(pHead string, pType string, pValue string, rowId int) (allowance Allowance, err error) {
	allowanceType := pType
	allowanceAmount, err := ParseFloatForImport(pHead, strings.TrimSpace(pValue), rowId)
	allowance = Allowance{AllowanceType: allowanceType, Amount: allowanceAmount}
	return allowance, err
}

func ValidateHeader(header []string) (err error) {
	//totalIncome,wht,donation,k-receipt
	if len(header) != 3 {
		err = errors.Join(err, errors.New("invalid header count"))
	}
	if strings.TrimSpace(header[0]) != "totalIncome" {
		err = errors.Join(err, errors.New("invalid header totalIncome"))
	}
	if strings.TrimSpace(header[1]) != "wht" {
		err = errors.Join(err, errors.New("invalid header wht"))
	}
	if strings.TrimSpace(header[2]) != "donation" {
		err = errors.Join(err, errors.New("invalid header donation"))
	}
	// if strings.TrimSpace(header[3]) != "k-receipt" {
	// 	err = errors.Join(err, errors.New("Invalid header k-receipt"))
	// }

	return
}

func ParseFloatForImport(pFieldName string, pValue string, rowId int) (floatValue float64, err error) {
	floatResult, errCast := strconv.ParseFloat(strings.TrimSpace(pValue), 64)
	if errCast != nil {
		return 0.0, errors.Join(err, fmt.Errorf("error parse %s at row %d", pFieldName, rowId))
	}
	floatValue = floatResult
	return floatValue, nil
}
