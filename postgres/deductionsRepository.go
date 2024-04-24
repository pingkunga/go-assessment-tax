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

func (p *Repository) DeductionConfigs() ([]DeductionConfig, error) {
	rows, err := p.Db.Query("SELECT * FROM TAX_DEDUCTIONCONFIG")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	deductionConfigs := []DeductionConfig{}
	for rows.Next() {
		deductionConfig := DeductionConfig{}
		err := rows.Scan(&deductionConfig.Id, &deductionConfig.DeductionType, &deductionConfig.DeductionAmount, &deductionConfig.DeductionMin, &deductionConfig.DeductionMax)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		deductionConfigs = append(deductionConfigs, deductionConfig)
	}
	return deductionConfigs, nil
}

func (p *Repository) DeductionConfigByType(pKey string) (DeductionConfig, error) {
	stmt, err := p.Db.Prepare("SELECT * FROM TAX_DEDUCTIONCONFIG WHERE deduction_type = $1")
	if err != nil {
		log.Fatal("can't Prepare SQL for "+pKey, err)
	}
	rowId := pKey
	row := stmt.QueryRow(rowId)
	deductionConfig := DeductionConfig{}

	err = row.Scan(&deductionConfig.Id, &deductionConfig.DeductionType, &deductionConfig.DeductionAmount, &deductionConfig.DeductionMin, &deductionConfig.DeductionMax)
	if err != nil {
		log.Fatal(err)
		return DeductionConfig{}, err
	}

	return deductionConfig, nil
}
