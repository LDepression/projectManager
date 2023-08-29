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
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	common "projectManager/project-common"
	"projectManager/project-common/encrypts"
	"projectManager/project-common/errs"
	"projectManager/project-common/jwts"
	"projectManager/project-common/logs"
	"projectManager/project-common/tms"
	"projectManager/project-grpc/user/login"
	"projectManager/project-user/config"
	"projectManager/project-user/internal/dao"
	"projectManager/project-user/internal/data/member"
	"projectManager/project-user/internal/data/organization"
	"projectManager/project-user/internal/database"
	"projectManager/project-user/internal/database/tran"
	"projectManager/project-user/internal/repo"
	"projectManager/project-user/pkg/model"
	"strconv"
	"strings"
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

func (ls LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	//去数据库中查询
	pwd := encrypts.Md5(msg.Password)
	mem, err := ls.memberRepo.FindMember(ctx, msg.Account, pwd)
	if err != nil {
		return &login.LoginResponse{}, errs.GrpcError(model.PwdError)
	}
	if mem == nil {
		return &login.LoginResponse{}, errs.GrpcError(model.PwdError)
	}
	memMsg := &login.MemberMessage{}
	err = copier.Copy(memMsg, mem)
	memMsg.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKEY)
	memMsg.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	memMsg.CreateTime = tms.FormatByMill(mem.CreateTime)
	//查询组织
	orgs, err := ls.organizationRepo.FindOrganizationsByMemID(ctx, mem.Id)
	if err != nil {
		zap.L().Error("login db findMember error", zap.Error(err))
		return &login.LoginResponse{}, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)

	for _, v := range orgsMessage {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKEY)
		v.OwnerCode = memMsg.Code
		o := organization.ToMap(orgs)[v.Id]
		v.CreateTime = tms.FormatByMill(o.CreateTime)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKEY)
	}
	//使用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.C.Jc.AccessExp*3600*24) * time.Second
	rExp := time.Duration(config.C.Jc.RefreshExp*3600*24) * time.Second
	token := jwts.CreateToken(memIdStr, exp, config.C.Jc.AccessSecret, rExp, config.C.Jc.RefreshSecret)
	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}

	//放入缓存中
	mar, _ := json.Marshal(mem)
	ls.cache.Put(ctx, model.Member+memIdStr, string(mar), exp)
	orgJson, _ := json.Marshal(orgs)
	ls.cache.Put(ctx, model.MeberOrganization+memIdStr, string(orgJson), exp)
	return &login.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}
func (ls LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	parseToken, err := jwts.ParseToken(token, config.C.Jc.AccessSecret)
	if err != nil {
		zap.L().Error("Login TokenVerify failed ", zap.Error(err))
		return &login.LoginResponse{}, errs.GrpcError(model.NoLogin)
	}
	//数据库查询 优化点:登录之后 应该把用户信息缓存起来
	//从缓存中查询 如果没有 直接返回认证失败
	memJson, err := ls.cache.Get(context.Background(), model.Member+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get member error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	if memJson == "" {
		zap.L().Error("TokenVerify cache get member expire")
		return nil, errs.GrpcError(model.NoLogin)
	}
	memberById := &member.Member{}
	json.Unmarshal([]byte(memJson), memberById)
	//数据库查询 优化点 登录之后 应该把用户信息缓存起来
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKEY)

	orgsJson, err := ls.cache.Get(context.Background(), model.MeberOrganization+parseToken)
	if err != nil {
		zap.L().Error("TokenVerify cache get organization error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	if orgsJson == "" {
		zap.L().Error("TokenVerify cache get organization expire")
		return nil, errs.GrpcError(model.NoLogin)
	}
	var orgs []*organization.Organization
	json.Unmarshal([]byte(orgsJson), &orgs)
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKEY)
	}
	return &login.LoginResponse{Member: memMsg}, nil
}
func (ls *LoginService) FindMemInfoById(ctx context.Context, msg *login.UserMessage) (*login.MemberMessage, error) {
	memberById, err := ls.memberRepo.FindMemberByID(context.Background(), msg.MemId)
	if err != nil {
		zap.L().Error("TokenVerify db FindMemberById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMsg := &login.MemberMessage{}
	copier.Copy(memMsg, memberById)
	memMsg.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKEY)
	orgs, err := ls.organizationRepo.FindOrganizationsByMemID(context.Background(), memberById.Id)
	if err != nil {
		zap.L().Error("TokenVerify db FindMember error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, model.AESKEY)
	}
	memMsg.CreateTime = tms.FormatByMill(memberById.CreateTime)
	return memMsg, nil
}
func (ls *LoginService) FindMemInfoByIds(ctx context.Context, msg *login.UserMessage) (*login.MemberMessageList, error) {
	memberList, err := ls.memberRepo.FindMemberByIds(context.Background(), msg.MemIDs)
	if err != nil {
		zap.L().Error("FindMemInfoByIds db memberRepo.FindMemberByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if memberList == nil || len(memberList) <= 0 {
		return &login.MemberMessageList{List: nil}, nil
	}
	mMap := make(map[int64]*member.Member)
	for _, v := range memberList {
		mMap[v.Id] = v
	}
	var memMsgs []*login.MemberMessage
	copier.Copy(&memMsgs, memberList)
	for _, v := range memMsgs {
		m := mMap[v.Id]
		v.CreateTime = tms.FormatByMill(m.CreateTime)
		v.Code = encrypts.EncryptNoErr(v.Id)
	}

	return &login.MemberMessageList{List: memMsgs}, nil
}
