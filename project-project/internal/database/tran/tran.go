/**
 * @Author: lenovo
 * @Description:
 * @File:  tran
 * @Version: 1.0.0
 * @Date: 2023/07/19 21:07
 */

package tran

import "projectManager/project-project/internal/database"

// Transaction 事务的操作一定和数据库有关 注入数据库的连接 gorm.db
type Transaction interface {
	Action(func(conn database.DbConn) error) error
}
