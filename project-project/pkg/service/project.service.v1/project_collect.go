/**
 * @Author: lenovo
 * @Description:
 * @File:  project_collect
 * @Version: 1.0.0
 * @Date: 2023/08/01 15:27
 */

package project_service_v1

import (
	"context"
	"go.uber.org/zap"
	"projectManager/project-common/errs"
	"projectManager/project-grpc/project"
	"projectManager/project-project/internal/data/pro"
	"projectManager/project-project/pkg/model"
	"strconv"
	"time"
)

func (ps *ProjectService) UpdateCollectProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.CollectProjectResponse, error) {
	projectCode, _ := strconv.ParseInt(msg.ProjectCode, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var err error
	if "collect" == msg.CollectType {
		pc := &pro.ProjectCollection{
			ProjectCode: projectCode,
			MemberCode:  msg.MemberId,
			CreateTime:  time.Now().UnixMilli(),
		}
		err = ps.projectRepo.SaveProjectCollect(c, pc)
	}
	if "cancel" == msg.CollectType {
		err = ps.projectRepo.DeleteProjectCollect(c, msg.MemberId, projectCode)
	}
	if err != nil {
		zap.L().Error("project UpdateCollectProject SaveProjectCollect error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.CollectProjectResponse{}, nil
}
