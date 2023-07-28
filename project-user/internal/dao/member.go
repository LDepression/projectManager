/**
 * @Author: lenovo
 * @Description:
 * @File:  member
 * @Version: 1.0.0
 * @Date: 2023/07/19 17:11
 */

package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"projectManager/project-user/internal/data/member"
	"projectManager/project-user/internal/database"
	"projectManager/project-user/internal/database/gorms"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func NewMemberDao() *MemberDao {
	return &MemberDao{
		conn: gorms.New(),
	}
}

func (m MemberDao) SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error {
	m.conn = conn.(*gorms.GormConn)
	return m.conn.Tx(ctx).Create(mem).Error
}

func (m MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("email=?", email).Count(&count).Error
	return count > 0, err

}

func (m MemberDao) GetMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}

func (m MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}

func (m MemberDao) FindMember(ctx context.Context, account, password string) (*member.Member, error) {
	var mem member.Member
	err := m.conn.Session(ctx).Where("account = ? and password =?", account, password).First(&mem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &mem, err
}

func (m MemberDao) FindMemberByID(ctx context.Context, id int64) (*member.Member, error) {
	var mem member.Member
	err := m.conn.Session(ctx).Where("id =?", id).Find(&mem).Error
	return &mem, err
}
