/**
 * @Author: lenovo
 * @Description:
 * @File:  organization
 * @Version: 1.0.0
 * @Date: 2023/07/19 15:38
 */

package organization

type Organization struct {
	Id          int64
	Name        string
	Avatar      string
	Description string
	MemberId    int64
	CreateTime  int64
	Personal    int32
	Address     string
	Province    int32
	City        int32
	Area        int32
}

func (*Organization) TableName() string {
	return "ms_organization"
}