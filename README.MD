# diff

对某个struct要做修改时，要更新到数据库。通过diff包，可以按需生成sql，只包含需要修改的字段。

```go

func TestTraceV2(t *testing.T) {
	now := time.Now().Add(-time.Hour * 10).UTC()
	user := User{
		UserId:       54746881978,
		UserName:     "test",
		UserPassword: "test",
		UserStatus:   1,
		UpdatedAt:    &now,
	}
	diff := TraceValue(user, func(user *User) {
		now := time.Now().UTC()
		user.UserName = "test2"
		user.UserStatus = 2
		user.UpdatedAt = &now
	})
	sql := buildUpdateSql(diff.Props, diff.Entity)
    fmt.Println(sql)
}
```

输出SQL：

```sql
UPDATE "sys_user" SET "user_name" = 'test2',"user_status" = 2,"updated_at" = '2025-05-16T05:29:28Z' WHERE "user_id" = 54746881978
```