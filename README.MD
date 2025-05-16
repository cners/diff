# diff

对某个struct要做修改时，要更新到数据库。通过diff包，可以按需生成sql，只包含需要修改的字段。

## 使用

```
go get github.com/cners/diff
```

```go

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

```

输出SQL：

```sql
UPDATE "sys_user" SET "updated_at" = '2025-05-16T07:16:40Z',"user_name" = 'test2',"user_status" = 2,"user_gender" = 2,"is_deleted" = true WHERE "user_id" = 54746881978
```