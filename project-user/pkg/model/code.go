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
	NoLegalMobile = errs.NewError(2001, "手机号不合法") //手机号不合法
)
