/**
 * @Author: lenovo
 * @Description:
 * @File:  menu
 * @Version: 1.0.0
 * @Date: 2023/07/21 16:19
 */

package dao

import (
	"context"
	"projectManager/project-project/internal/data/menu"
	"projectManager/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func NewMenuRepo() *MenuDao {
	return &MenuDao{conn: gorms.New()}
}
func (m MenuDao) FindMenus(ctx context.Context) (pms []*menu.ProjectMenu, err error) {
	session := m.conn.Session(ctx)
	err = session.Order("pid,sort asc, id asc").Find(&pms).Error
	return
}
