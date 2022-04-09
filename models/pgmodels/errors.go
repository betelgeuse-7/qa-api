package pgmodels

import (
	"errors"
	"fmt"
)

func err_QUERY_BUILDING_FAIL() error {
	return errors.New("query building failed")
}

func err_DB_EXEC_FAIL(execErr error) error {
	return fmt.Errorf("*sqlx.DB.Exec() fail: %s", execErr.Error())
}

func err_DB_NO_ROWS() error {
	return fmt.Errorf("no record found")
}
