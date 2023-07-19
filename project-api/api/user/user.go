/**
 * @Author: lenovo
 * @Description:
 * @File:  user
 * @Version: 1.0.0
 * @Date: 2023/07/16 22:48
 */

package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"net/http"
	"projectManager/project-api/pkg/model/user"
	common "projectManager/project-common"
	"projectManager/project-common/errs"
	"projectManager/project-grpc/user/login"
	"time"
)

type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}
func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	mobile := ctx.PostForm("mobile")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rsp, err := LoginServiceClient.GetCaptcha(c, &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.HandleGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(rsp.Code))
}

func (u *HandlerUser) register(c *gin.Context) {
	//1.接收参数 参数模型
	result := common.Result{}
	var req user.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "参数错误"))
		return
	}

	//2.校验参数是否合法
	if err := req.Verify(); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	//3.调用user的grpc的服务
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &login.RegisterMessage{}
	if err := copier.Copy(msg, req); err != nil {
		c.JSON(http.StatusOK, result.Fail(http.StatusBadRequest, "copy有误"))
		return
	}

	rsp, err := LoginServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.HandleGrpcError(err)
		zap.S().Info(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	//4.返回结果
	c.JSON(http.StatusOK, result.Success(rsp))
}
