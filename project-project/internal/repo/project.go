/**
 * @Author: lenovo
 * @Description:
 * @File:  project
 * @Version: 1.0.0
 * @Date: 2023/07/21 20:47
 */

package repo

import (
	"context"
	"projectManager/project-project/internal/data/pro"
	"projectManager/project-project/internal/database"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*pro.ProjectMemberUnion, int64, error)
	FindCollectProjectByMemId(ctx context.Context, id int64, page int64, size int64) ([]*pro.ProjectMemberUnion, int64, error)
	SaveProject(conn database.DbConn, ctx context.Context, pr *pro.Project) error
	SaveProjectMember(conn database.DbConn, ctx context.Context, pm *pro.ProjectMember) error
	FindProjectByPIdAndMemId(ctx context.Context, projectCode int64, memberId int64) (*pro.ProjectMemberUnion, error)
	FindCollectByPidAndMemId(ctx context.Context, projectCode int64, memberId int64) (bool, error)
	UpdateDeletedProject(ctx context.Context, projectCode int64, deleted bool) error
	SaveProjectCollect(ctx context.Context, pc *pro.ProjectCollection) error
	DeleteProjectCollect(ctx context.Context, id int64, code int64) error
	UpdateProject(ctx context.Context, proj *pro.Project) error
}

type ProjectTemplateRepo interface {
	FindProjectTemplateSystem(ctx context.Context, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
	FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
	FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
}
