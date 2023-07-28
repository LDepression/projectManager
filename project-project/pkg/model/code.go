/**
 * @Author: lenovo
 * @Description:
 * @File:  code
 * @Version: 1.0.0
 * @Date: 2023/07/15 18:55
 */

package model

import (
	"projectManager/project-common/errs"
)

var (
	RedisError     = errs.NewError(-100, "redis错误") //redis不正确
	DBError        = errs.NewError(-101, "数据库错误")
	NoLegalMobile  = errs.NewError(2001, "手机号不合法")   //手机号不合法
	NoLegalCaptcha = errs.NewError(2002, "验证码不正确")   //验证码不正确
	EmailExist     = errs.NewError(2003, "邮箱已经存在了")  //验证码不正确
	AccountExist   = errs.NewError(2004, "账号已经存在了")  //账号已经存在了
	MobileExist    = errs.NewError(2005, "手机号已经存在了") //手机号已经存在了
	ErrCaptcha     = errs.NewError(2006, "验证码不正确")   //验证码不存在或者过期
	PwdError       = errs.NewError(2007, "密码不正确")    //密码不正确

)
