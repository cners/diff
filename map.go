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
		if !field.IsNil() { // 只处理非nil指针字段
			fieldName := updatesValue.Type().Field(i).Name
			destField := destValue.FieldByName(fieldName)
			if destField.IsValid() && destField.CanSet() {
				destField.Set(field.Elem()) // 解引用指针
			}
		}
	}
	return nil
}
