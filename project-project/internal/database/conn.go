/**
 * @Author: lenovo
 * @Description:
 * @File:  conn
 * @Version: 1.0.0
 * @Date: 2023/07/19 21:14
 */

package database

type DbConn interface {
	Begin()
	Rollback()
	Commit()
}
