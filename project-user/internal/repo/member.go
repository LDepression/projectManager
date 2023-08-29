/**
 * @Author: lenovo
 * @Description:
 * @File:  member
 * @Version: 1.0.0
 * @Date: 2023/07/19 16:57
 */

package repo

import (
	"context"
	"projectManager/project-user/internal/data/member"
	"projectManager/project-user/internal/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, name string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (mem *member.Member, err error)
	FindMemberByID(ctx context.Context, id int64) (mem *member.Member, err error)
	FindMemberByIds(background context.Context, ids []int64) (list []*member.Member, err error)
}
