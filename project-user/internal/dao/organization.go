/**
 * @Author: lenovo
 * @Description:
 * @File:  organization
 * @Version: 1.0.0
 * @Date: 2023/07/19 20:44
 */

package dao

import (
	"context"
	"projectManager/project-user/internal/data/organization"
	"projectManager/project-user/internal/database"
	"projectManager/project-user/internal/database/gorms"
)

type OrganizationConn struct {
	conn *gorms.GormConn
}

func (o *OrganizationConn) FindOrganizationsByMemID(ctx context.Context, id int64) ([]*organization.Organization, error) {
	var orgs []*organization.Organization
	err := o.conn.Session(ctx).Where("member_id = ?", id).Find(&orgs).Error
	return orgs, err
}

func NewOrganizationDao() *OrganizationConn {
	return &OrganizationConn{
		conn: gorms.New(),
	}
}

func (o *OrganizationConn) SaveOrganization(conn database.DbConn, ctx context.Context, org *organization.Organization) error {
	o.conn = conn.(*gorms.GormConn)
	return o.conn.Tx(ctx).Create(org).Error
}
