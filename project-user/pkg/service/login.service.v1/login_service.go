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
	"go.uber.org/zap"
	"log"
	common "projectManager/project-common"
	"projectManager/project-common/errs"
	"projectManager/project-common/logs"
	"projectManager/project-user/pkg/dao"
	"projectManager/project-user/pkg/model"
	"projectManager/project-user/pkg/repo"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	cache repo.Cache
}

func New() *LoginService {
	return &LoginService{
		cache: dao.Rc,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, req *CaptchaMessage) (*CaptchaResponse, error) {
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
		if err := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Second); err != nil {
			log.Printf("验证码存入redis出错,cause by: %v\n", err)
		}
		log.Printf("将手机号和验证码存入redis成功: REGISTER_%s : %s", mobile, code)
	}()
	return &CaptchaResponse{Code: code}, nil
}
