/**
 * @Author: lenovo
 * @Description:
 * @File:  page
 * @Version: 1.0.0
 * @Date: 2023/07/21 21:27
 */

package model

import "github.com/gin-gonic/gin"

type Page struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"pageSize" form:"pageSize"`
}

func (p *Page) Bind(c *gin.Context) {
	c.ShouldBind(&p)
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}
}
