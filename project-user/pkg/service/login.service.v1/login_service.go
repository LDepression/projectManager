/**
 * @Author: lenovo
 * @Description:
 * @File:  login_service
 * @Version: 1.0.0
 * @Date: 2023/07/16 20:04
 */

package login_service_v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"log"
	common "projectManager/project-common"
	"projectManager/project-common/encrypts"
	"projectManager/project-common/errs"
	"projectManager/project-common/logs"
	"projectManager/project-grpc/user/login"
	"projectManager/project-user/internal/dao"
	"projectManager/project-user/internal/data/member"
	"projectManager/project-user/internal/data/organization"
	"projectManager/project-user/internal/database"
	"projectManager/project-user/internal/database/tran"
	"projectManager/project-user/internal/repo"
	"projectManager/project-user/pkg/model"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.Organization

	tran tran.Transaction
}

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		tran:             dao.NewTransaction(),
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, req *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	//1.获取参数
	mobile := req.Mobile
	//2.校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.生成验证码(随机4位或者6位
	code := "123456"
	//4.调用短信平台,放入短信平台(第三方, 放入Go程中执行,接口可以快速响应)
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功,发送短信 INFO")
		logs.LG.Debug("短信平台调用成功,发送短信 debug")
		zap.L().Error("短信平台调用成功,发送短信 error")

		//redis 假设可能缓存在mysql,mongo当中,如果我们直接写的话,就没有低耦合了
		/*
				我们应该抽象一层出去,只需要实现其接口就好了.
			这样我们也不至于换成其他中间件的时候改代码
		*/
		//5.储存验证码 redis当中 过期15分钟
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := ls.cache.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Second); err != nil {
			log.Printf("验证码存入redis出错,cause by: %v\n", err)
		}
		log.Printf("将手机号和验证码存入redis成功: REGISTER_%s : %s", mobile, code)
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, req *login.RegisterMessage) (*login.RegisterResponse, error) {
	c := context.Background()
	//1.可以去校验一下参数
	//2.验证码是正确的
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+req.Mobile)
	if errors.Is(err, redis.Nil) {
		return nil, errs.GrpcError(model.ErrCaptcha)
	}
	if err != nil {
		zap.L().Error("Register redis get error,", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.NoLegalCaptcha)
	}
	//3.校验业务逻辑(邮箱是否是注册 账号是否是是被注册 手机号是否被注册)
	exist, err := ls.memberRepo.GetMemberByEmail(c, req.Email)
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	//判断账号
	exist, err = ls.memberRepo.GetMemberByAccount(c, req.Name)
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}

	exist, err = ls.memberRepo.GetMemberByAccount(c, req.Name)
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}

	err = ls.tran.Action(func(conn database.DbConn) error {
		//4.执行业务 将数据存入member表中 生成一个数据
		pwd := encrypts.Md5(req.Password)
		mem := &member.Member{
			Account:       req.Name,
			Password:      pwd,
			Name:          req.Name,
			Mobile:        req.Mobile,
			Email:         req.Email,
			CreateTime:    time.Now().UnixMilli(),
			LastLoginTime: time.Now().UnixMilli(),
			Status:        model.Normal,
		}
		if err := ls.memberRepo.SaveMember(conn, c, mem); err != nil {
			return errs.GrpcError(model.DBError)
		}
		//存入组织
		org := &organization.Organization{
			Name:       mem.Name + "个人项目",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc -ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return model.DBError
		}
		return nil
	})
	return &login.RegisterResponse{}, err
}
