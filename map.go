package diff

import (
	"reflect"
)

/*
为目标结构体批量更新字段值。
注意：推荐传入的字段都是指针类型，空指针会被忽略更新。
例如：

	type UpdateUser struct {
		UserName   *string
		UserStatus *int
	}
*/
func UpdateMap(dest interface{}, updates interface{}) error {
	destValue := reflect.ValueOf(dest).Elem()
	updatesValue := reflect.ValueOf(updates).Elem()
	for i := 0; i < updatesValue.NumField(); i++ {
		field := updatesValue.Field(i)
		fieldName := updatesValue.Type().Field(i).Name
		destField := destValue.FieldByName(fieldName)
		if destField.IsValid() && destField.CanSet() {
			// 检查字段类型是否匹配
			if field.Kind() == reflect.Ptr {
				if !field.IsNil() && field.Elem().Type().AssignableTo(destField.Type()) {
					destField.Set(field.Elem())
				}
			} else if field.Type().AssignableTo(destField.Type()) {
				destField.Set(field)
			}
		}
	}
	return nil
}
