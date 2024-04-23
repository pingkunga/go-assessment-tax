package postgres

import (
	"log"
)

type DeductionConfig struct {
	Id              int     `json:"id"`
	DeductionType   string  `json:"deductionType"`
	DeductionAmount float64 `json:"deductionAmount"`
	DeductionMin    float64 `json:"deductionMin"`
	DeductionMax    float64 `json:"adeductionMax"`
}

func (p *Repository) GetDeductionConfig(pKey string) (float64, error) {

	stmt, err := p.Db.Prepare("SELECT deduction_amount FROM TAX_DEDUCTIONCONFIG WHERE deduction_type = $1")
	if err != nil {
		log.Fatal("can't Prepare SQL for "+pKey, err)
	}
	rowId := pKey
	row := stmt.QueryRow(rowId)
	var personalDeduction float64

	err = row.Scan(&personalDeduction)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return personalDeduction, nil
}

func (p *Repository) SetDeductionConfig(pKey string, newValue float64) error {
	stmt, err := p.Db.Prepare("UPDATE TAX_DEDUCTIONCONFIG SET deduction_amount = $1 WHERE deduction_type = $2")
	if err != nil {
		log.Fatal("can't Prepare SQL for Update PersonalDeduction", err)
	}

	rowId := pKey
	_, err = stmt.Exec(newValue, rowId)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (p *Repository) AllowanceConfigs() ([]DeductionConfig, error) {
	rows, err := p.Db.Query("SELECT * FROM TAX_DEDUCTIONCONFIG")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	allowanceConfigs := []DeductionConfig{}
	for rows.Next() {
		allowanceConfig := DeductionConfig{}
		err := rows.Scan(&allowanceConfig.Id, &allowanceConfig.DeductionType, &allowanceConfig.DeductionAmount, &allowanceConfig.DeductionMin, &allowanceConfig.DeductionMax)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		allowanceConfigs = append(allowanceConfigs, allowanceConfig)
	}
	return allowanceConfigs, nil
}
