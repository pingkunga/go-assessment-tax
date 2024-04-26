package tax

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

func CalculateTax(tax TaxRequest) TaxResponse {

	netIncome := calculateNetIncome(tax)

	var taxAmount float64
	var taxLevels []TaxLevel
	for _, step := range TaxStepList() {
		stepTaxAmount := 0.0
		if netIncome <= 0 {
			stepTaxAmount = 0
		} else if netIncome > step.MaxIncome {
			stepTaxAmount = (step.MaxIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		} else {
			stepTaxAmount = (netIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		}

		taxLevels = append(taxLevels, TaxLevel{Level: step.FriendlyMessage, Tax: stepTaxAmount})
		taxAmount = taxAmount + stepTaxAmount
	}

	taxAmount = taxAmount - tax.WHT

	if taxAmount < 0 {
		return TaxResponse{Tax: 0, TaxRefund: taxAmount * -1, TaxLevels: taxLevels}
	}
	return TaxResponse{Tax: taxAmount, TaxRefund: 0, TaxLevels: taxLevels}
}

func calculateNetIncome(tax TaxRequest) float64 {
	netIncome := tax.TotalIncome - PersonalDeduction()

	for _, allowance := range tax.Allowances {
		netIncome = netIncome - allowanceAmountGuard(allowance)
	}

	return netIncome
}

func allowanceAmountGuard(allowance Allowance) float64 {

	if allowance.AllowanceType == DEDUCTION_DONATION && allowance.Amount > DEDUCTION_DONATION_MAX {
		return DEDUCTION_DONATION_MAX
	}
	return allowance.Amount
}

const DEDUCTION_DONATION = "donation"
const DEDUCTION_DONATION_MAX = 100000

const DEDUCTION_K_RECEIPT = "k-receipt"

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

func PersonalDeduction() float64 {
	return 60000
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
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s ขึ้นไป", humanize.FormatFloat("#,###.##", step.MinIncome))
		} else {
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s-%s", humanize.FormatFloat("#,###.##", step.MinIncome), humanize.FormatFloat("#,###.", step.MaxIncome))
		}
	}

	return taxStep
}

func ImportTaxCSV(pathToFile string) (taxls []TaxRequest, err error) {

	//read file from path e.g. Uploads/taxes.csv
	csvFile, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Error Open CSV file at: " + pathToFile)
	}
	defer csvFile.Close()

	//parse csv file
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Error read CSV file at: " + pathToFile)
	}
	//loop through each row
	var taxRequests []TaxRequest
	for i, row := range csvData {
		if i == 0 {
			if err := ValidateHeader(row); err != nil {
				log.Fatal(err)
				return nil, errors.New(fmt.Sprintf("Error validate header at row %d", i))
			} else {
				continue
			}
		}
		totalIncome, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(fmt.Sprintf("Error parse TotalIncome at row %d", i))
		}

		wht, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(fmt.Sprintf("Error parse WHT at row %d", i))
		}
		allowances, err := createAllowances(row, i)
		if err != nil {
			log.Fatal(err)
			return nil, errors.New(fmt.Sprintf("Error create Allowances at row %d", i))
		}

		TaxRequest := TaxRequest{
			TotalIncome: totalIncome,
			WHT:         wht,
			Allowances:  allowances,
		}
		taxRequests = append(taxRequests, TaxRequest)
	}

	//validate each row
	return taxRequests, nil
}

func createAllowances(row []string, rowId int) (allowances []Allowance, err error) {

	//donation
	allowanceDonation, err := createAllowance(DEDUCTION_DONATION, row[2], rowId)
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

func createAllowance(pType string, pValue string, rowId int) (allowance Allowance, err error) {
	allowanceType := pType
	allowanceAmount, err := ParseFloatForImport(pValue, "Allowance Amount", rowId)
	allowance = Allowance{AllowanceType: allowanceType, Amount: allowanceAmount}
	return allowance, err
}

func ValidateHeader(header []string) (err error) {
	//totalIncome,wht,donation,k-receipt
	if len(header) != 4 {
		errors.Join(err, errors.New("Invalid header"))
	}
	if header[0] != "totalIncome" {
		errors.Join(err, errors.New("Invalid header totalIncome"))
	}
	if header[1] != "wht" {
		errors.Join(err, errors.New("Invalid header wht"))
	}
	if header[2] != "donation" {
		errors.Join(err, errors.New("Invalid header donation"))
	}
	// if header[3] != "k-receipt" {
	// 	errors.Join(err, errors.New("Invalid header k-receipt"))
	// }

	return err
}

func ParseFloatForImport(pValue string, pFieldName string, rowId int) (floatValue float64, err error) {
	floatResult, err := strconv.ParseFloat(pValue, 64)
	if err != nil {
		return 0.0, errors.Join(err, errors.New(fmt.Sprintf("Error parse %s at row %d", pFieldName, rowId)))
	}
	floatValue = floatResult
	return floatValue, nil
}
