/**
 * @Author: lenovo
 * @Description:
 * @File:  login_service
 * @Version: 1.0.0
 * @Date: 2023/07/16 20:04
 */

package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"projectManager/project-common/encrypts"
	"projectManager/project-common/errs"
	"projectManager/project-common/tms"
	"projectManager/project-grpc/project"
	"projectManager/project-grpc/user/login"
	"projectManager/project-project/internal/dao"
	"projectManager/project-project/internal/data/menu"
	"projectManager/project-project/internal/data/pro"
	"projectManager/project-project/internal/data/task"
	"projectManager/project-project/internal/database"
	"projectManager/project-project/internal/database/tran"
	"projectManager/project-project/internal/repo"
	"projectManager/project-project/internal/rpc"
	"projectManager/project-project/pkg/model"
	"strconv"
	"time"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer
	cache                  repo.Cache
	tran                   tran.Transaction
	menuRepo               repo.MenuRepo
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
}

func New() *ProjectService {
	return &ProjectService{
		cache:                  dao.Rc,
		tran:                   dao.NewTransaction(),
		menuRepo:               dao.NewMenuRepo(),
		projectRepo:            dao.NewProjectRepo(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
	}
}
func (p *ProjectService) Index(ctx context.Context, message *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("index db FindMenus error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	childs := menu.CovertChild(pms)
	var mms []*project.MenuMessage
	copier.Copy(&mms, childs)
	return &project.IndexResponse{Menus: mms}, nil
}
func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.MyProjectResponse, error) {
	memberId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	var pms []*pro.ProjectMemberUnion
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, page, pageSize, "")
	}
	if msg.SelectBy == "archive" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, page, pageSize, "and archive =1")
	}
	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, page, pageSize, "and deleted =1")
	}
	if msg.SelectBy == "collect" {
		pms, total, err = p.projectRepo.FindCollectProjectByMemId(ctx, memberId, page, pageSize)
		for _, v := range pms {
			v.Collected = model.Collected
		}
	} else {
		collectPms, _, err := p.projectRepo.FindCollectProjectByMemId(ctx, memberId, page, pageSize)
		if err != nil {
			zap.L().Error("project FindProjectByMemId error ", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		var cMap = make(map[int64]*pro.ProjectMemberUnion)
		for _, v := range collectPms {
			cMap[v.Id] = v
		}
		for _, v := range pms {
			if cMap[v.ProjectCode] != nil {
				v.Collected = model.Collected
			}
		}
	}

	if err != nil {
		zap.L().Error("project FindProjectByMemId error ", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if pms == nil {
		return &project.MyProjectResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}
	var pmm []*project.ProjectMessage
	copier.Copy(&pmm, pms)
	for _, v := range pmm {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKEY)
		pam := pro.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKEY)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &project.MyProjectResponse{Pm: pmm, Total: total}, nil
}

func (ps *ProjectService) FindProjectTemplate(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectTemplateResponse, error) {
	//1.根据viewType去查询项目模板表 得到list
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKEY)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	page := msg.Page
	pageSize := msg.PageSize
	var pts []pro.ProjectTemplate
	var total int64
	var err error
	if msg.ViewType == -1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateAll(ctx, organizationCode, page, pageSize)
	}
	if msg.ViewType == 0 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateCustom(ctx, msg.MemberId, organizationCode, page, pageSize)
	}
	if msg.ViewType == 1 {
		pts, total, err = ps.projectTemplateRepo.FindProjectTemplateSystem(ctx, page, pageSize)
	}
	if err != nil {
		zap.L().Error("project FindProjectTemplate FindProjectTemplateSystem error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//2.模型转换，拿到模板id列表 去 任务步骤模板表 去进行查询
	tsts, err := ps.taskStagesTemplateRepo.FindInProTemIds(ctx, pro.ToProjectTemplateIds(pts))
	if err != nil {
		zap.L().Error("project FindProjectTemplate FindInProTemIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var ptas []*pro.ProjectTemplateAll
	for _, v := range pts {
		//写代码 该谁做的事情一定要交出去
		ptas = append(ptas, v.Convert(task.CovertProjectMap(tsts)[v.Id]))
	}
	//3.组装数据
	var pmMsgs []*project.ProjectTemplateMessage
	copier.Copy(&pmMsgs, ptas)
	return &project.ProjectTemplateResponse{Ptm: pmMsgs, Total: total}, nil
}

func (ps *ProjectService) SaveProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.SaveProjectMessage, error) {
	organizationCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, model.AESKEY)
	organizationCode, _ := strconv.ParseInt(organizationCodeStr, 10, 64)
	templateCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, model.AESKEY)
	templateCode, _ := strconv.ParseInt(templateCodeStr, 10, 64)
	//1. 保存项目表
	pr := &pro.Project{
		Name:              msg.Name,
		Description:       msg.Description,
		TemplateCode:      int(templateCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://img2.baidu.com/it/u=792555388,2449797505&fm=253&fmt=auto&app=138&f=JPEG?w=667&h=500",
		Deleted:           model.NoDeleted,
		Archive:           model.NoArchive,
		OrganizationCode:  organizationCode,
		AccessControlType: model.Open,
		TaskBoardTheme:    model.Simple,
	}
	err := ps.tran.Action(func(conn database.DbConn) error {
		err := ps.projectRepo.SaveProject(conn, ctx, pr)
		if err != nil {
			zap.L().Error("project SaveProject SaveProject error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		pm := &pro.ProjectMember{
			ProjectCode: pr.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   "",
		}
		//2. 保存项目和成员的关联表
		err = ps.projectRepo.SaveProjectMember(conn, ctx, pm)
		if err != nil {
			zap.L().Error("project SaveProject SaveProjectMember error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	code, _ := encrypts.EncryptInt64(pr.Id, model.AESKEY)
	rsp := &project.SaveProjectMessage{
		Id:               pr.Id,
		Code:             code,
		OrganizationCode: organizationCodeStr,
		Name:             pr.Name,
		Cover:            pr.Cover,
		CreateTime:       tms.FormatByMill(pr.CreateTime),
		TaskBoardTheme:   pr.TaskBoardTheme,
	}
	return rsp, nil
}

// FindProjectDetail 1. 查项目表
// 2. 项目和成员的关联表 查到项目的拥有者 去member表中查名字
// 3 查收藏项目 判断收藏状态
func (ps *ProjectService) FindProjectDetail(ctx context.Context, msg *project.ProjectRpcMessage) (*project.ProjectDetailMessage, error) {
	memberId := msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectCode, _ := strconv.Atoi(msg.ProjectCode)
	projectAndMember, err := ps.projectRepo.FindProjectByPIdAndMemId(c, int64(projectCode), memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindProjectByPIdAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	ownerId := projectAndMember.IsOwner
	member, err := rpc.LoginServiceClient.FindMemInfoById(c, &login.UserMessage{MemId: ownerId})
	if err != nil {
		zap.L().Error("project rpc FindProjectDetail FindMemInfoById error", zap.Error(err))
		return nil, err
	}
	//去user模块去找了
	//TODO 优化 收藏的时候 可以放入redis
	isCollect, err := ps.projectRepo.FindCollectByPidAndMemId(c, int64(projectCode), memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindCollectByPidAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if isCollect {
		projectAndMember.Collected = model.Collected
	}
	var detailMsg = &project.ProjectDetailMessage{}
	copier.Copy(detailMsg, projectAndMember)
	detailMsg.OwnerAvatar = member.Avatar
	detailMsg.OwnerName = member.Name
	detailMsg.Code, _ = encrypts.EncryptInt64(projectAndMember.Id, model.AESKEY)
	detailMsg.AccessControlType = projectAndMember.GetAccessControlType()
	detailMsg.OrganizationCode, _ = encrypts.EncryptInt64(projectAndMember.OrganizationCode, model.AESKEY)
	detailMsg.Order = int32(projectAndMember.Sort)
	detailMsg.CreateTime = tms.FormatByMill(projectAndMember.CreateTime)
	return detailMsg, nil
}
func (ps *ProjectService) UpdateDeletedProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.DeletedProjectResponse, error) {
	projectCode, _ := strconv.ParseInt(msg.ProjectCode, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := ps.projectRepo.UpdateDeletedProject(c, projectCode, msg.Deleted)
	if err != nil {
		zap.L().Error("project RecycleProject DeleteProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.DeletedProjectResponse{}, nil
}

func (ps *ProjectService) UpdateProject(ctx context.Context, msg *project.UpdateProjectMessage) (*project.UpdateProjectResponse, error) {
	projectCode, _ := strconv.ParseInt(msg.ProjectCode, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	proj := &pro.Project{
		Id:                 projectCode,
		Name:               msg.Name,
		Description:        msg.Description,
		Cover:              msg.Cover,
		TaskBoardTheme:     msg.TaskBoardTheme,
		Prefix:             msg.Prefix,
		Private:            int(msg.Private),
		OpenPrefix:         int(msg.OpenPrefix),
		OpenBeginTime:      int(msg.OpenBeginTime),
		OpenTaskPrivate:    int(msg.OpenTaskPrivate),
		Schedule:           msg.Schedule,
		AutoUpdateSchedule: int(msg.AutoUpdateSchedule),
	}
	err := ps.projectRepo.UpdateProject(c, proj)
	if err != nil {
		zap.L().Error("project UpdateProject::UpdateProject error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &project.UpdateProjectResponse{}, nil
}
