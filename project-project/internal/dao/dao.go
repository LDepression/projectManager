/**
 * @Author: lenovo
 * @Description:
 * @File:  dao
 * @Version: 1.0.0
 * @Date: 2023/07/19 21:42
 */

package dao

import (
	"projectManager/project-project/internal/database"
	"projectManager/project-project/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DbConn
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTran(),
	}
}

func (t *TransactionImpl) Action(f func(tx database.DbConn) error) error {
	t.conn.Begin()
	if err := f(t.conn); err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}
