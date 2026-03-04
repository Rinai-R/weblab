package mysql

import (
	"errors"
	"weblab/internal/dao"

	drivermysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func normalizeErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dao.ErrNotFound
	}

	var mysqlErr *drivermysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return dao.ErrDuplicated
		case 1452:
			return dao.ErrNotFound
		}
	}

	return err
}
