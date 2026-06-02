package mysql

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// IsDuplicateKey reports MySQL duplicate-key errors (SQLSTATE 23000, errno 1062).
func IsDuplicateKey(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}

	return false
}
