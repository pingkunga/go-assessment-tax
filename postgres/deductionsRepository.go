package postgres

import "log"

func (p *Repository) GetAllowanceConfig(pKey string) (float64, error) {

	stmt, err := p.Db.Prepare("SELECT allowance_amount FROM TAX_ALLOWANCECONFIG WHERE allowance_type = $1")
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

func (p *Repository) SetAllowanceConfig(pKey string, newValue float64) error {
	stmt, err := p.Db.Prepare("UPDATE TAX_ALLOWANCECONFIG SET allowance_amount = $1 WHERE allowance_type = $2")
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
