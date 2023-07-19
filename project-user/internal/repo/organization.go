/**
 * @Author: lenovo
 * @Description:
 * @File:  organization
 * @Version: 1.0.0
 * @Date: 2023/07/19 20:39
 */

package repo

import (
	"context"
	"projectManager/project-user/internal/data/organization"
	"projectManager/project-user/internal/database"
)

type Organization interface {
	SaveOrganization(conn database.DbConn, ctx context.Context, org *organization.Organization) error
}
