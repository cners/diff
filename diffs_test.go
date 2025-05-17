package diff

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/copier"
)

type User struct {
	UserId     int64      `json:"user_id" gorm:"column:user_id;primaryKey;"`                                               // 用户id
	UserName   string     `json:"user_name" gorm:"column:user_name;type:varchar;not null;default:''"`                      // 用户名
	UserStatus int        `json:"user_status" gorm:"column:user_status;type:int;not null;default:1"`                       // 状态1正常2停用
	UserGender UserGender `json:"user_gender" gorm:"column:user_gender;type:int;not null;default:0"`                       // 性别0未设置1男2女
	UserRemark string     `json:"user_remark" gorm:"column:user_remark;type:varchar;not null;default:''"`                  // 备注
	IsDeleted  bool       `json:"is_deleted" gorm:"column:is_deleted;type:boolean;not null;default:false"`                 // 是否删除
	CreatedBy  int64      `json:"created_by" gorm:"column:created_by;type:bigint;not null;default:0"`                      // 创建人
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP"` // 创建时间
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamptz"`                                    // 修改时间
	UpdatedBy  int64      `json:"updated_by" gorm:"column:updated_by;type:bigint;not null;default:0"`                      // 修改人
}

func (*User) TableName() string {
	return "sys_user"
}

// 定义UserGender的枚举 0未知 1男 2女
type UserGender int

const (
	UserGenderUnknown UserGender = iota
	UserGenderMale
	UserGenderFemale
)

func (UserGender) Values() []UserGender {
	return []UserGender{
		UserGenderUnknown,
		UserGenderMale,
		UserGenderFemale,
	}
}

func TestDiff(t *testing.T) {
	old := &User{
		UserId:     1,
		UserName:   "test",
		UserStatus: 1,
	}
	var new User
	copier.Copy(&new, &old)
	new.UserName = "test2"
	changedFields := getChangedFields(*old, new)
	sql := buildUpdateSql(changedFields, new)
	t.Logf("\033[32mSQL: %s\033[0m", sql)
}

func TestTrace(t *testing.T) {
	user := User{
		UserId:     54746881978,
		UserName:   "test",
		UserStatus: 1,
	}
	now := time.Now().UTC()
	user, props := Trace(user, func(user *User) {
		user.UserName = "test2"
		user.UserStatus = 2
		user.UpdatedAt = &now
	})
	sql := buildUpdateSql(props, user)
	t.Log(sql)
}

func TestTraceV2(t *testing.T) {
	now := time.Now().Add(-time.Hour * 10).UTC()
	user := User{
		UserId:     54746881978,
		UserName:   "test",
		UserStatus: 1,
		UpdatedAt:  &now,
	}
	result := TraceProps(user, func(user *User) {
		now := time.Now().UTC()
		user.UserName = "test2"
		user.UserStatus = 2
		user.UserGender = UserGenderFemale
		user.UpdatedAt = &now
		user.IsDeleted = true
	})
	t.Log("\033[32m" + fmt.Sprintf("%+v", result) + "\033[0m")
	sql := BuildUpdateSql(result)
	t.Log("\033[32m" + sql + "\033[0m")
}

func TestTraceUpdate(t *testing.T) {
	now := time.Now().Add(-time.Hour * 10).UTC()
	user := User{
		UserId:     54746881978,
		UserName:   "test",
		UserStatus: 1,
		UpdatedAt:  &now,
	}
	result := TraceUpdate(user, func(user *User) {
		now := time.Now().UTC()
		user.UserName = "test2"
		user.UserStatus = 2
		user.UserGender = UserGenderFemale
		user.UpdatedAt = &now
		user.IsDeleted = true
	})
	t.Log("\033[32m" + result.UpdateSql + "\033[0m")
}

func TestTraceUpdateV2(t *testing.T) {
	type UpdateUser struct {
		UserName   *string
		UserStatus *int
		UserGender *UserGender
	}
	// mock the request the user info
	userName := "test3"
	userStatus := 2
	params_user := UpdateUser{
		UserName:   &userName,
		UserStatus: &userStatus,
	}
	now := time.Now().Add(-time.Hour * 10).UTC()
	user := User{
		UserId:     54746881978,
		UserName:   "test",
		UserStatus: 1,
		UpdatedAt:  &now,
	}

	result := TraceUpdate(user, func(user *User) {
		UpdateMap(user, &params_user)
		user.UpdatedAt = UTC
	})
	t.Log("\033[32m" + result.UpdateSql + "\033[0m")
}
