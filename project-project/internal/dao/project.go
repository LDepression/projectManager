/**
 * @Author: lenovo
 * @Description:
 * @File:  project
 * @Version: 1.0.0
 * @Date: 2023/07/21 20:55
 */

package dao

import (
	"context"
	"fmt"
	"projectManager/project-project/internal/data/pro"
	"projectManager/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p *ProjectDao) FindCollectProjectByMemId(ctx context.Context, memId int64, page int64, size int64) ([]*pro.ProjectMemberUnion, int64, error) {
	var pms []*pro.ProjectMemberUnion
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project where id in(select project_code from ms_project_collection where member_code =?) order by sort limit ?,? ")
	raw := session.Raw(sql, memId, index, size)
	raw.Scan(&pms)

	var total int64
	query := fmt.Sprintf("member_code=?")
	err := session.Model(&pro.ProjectCollection{}).Where(query, memId).Count(&total).Error
	return pms, total, err
}

func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*pro.ProjectMemberUnion, int64, error) {
	var pms []*pro.ProjectMemberUnion
	session := p.conn.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where  a.id = b.project_code and b.member_code =? %s order by sort limit ?,? ", condition)

	raw := session.Raw(sql, memId, index, size)
	raw.Scan(&pms)

	var total int64
	query := fmt.Sprintf("select count(*) from ms_project a, ms_project_member b where  a.id = b.project_code and b.member_code =? %s ", condition)
	tx := session.Raw(query, memId)
	err := tx.Scan(&total).Error
	return pms, total, err

}

func NewProjectRepo() *ProjectDao {
	return &ProjectDao{conn: gorms.New()}
}
