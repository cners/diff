package diff

import (
	"reflect"
	"strings"
	"time"
)

type PropValue struct {
	Type  reflect.Type
	Value reflect.Value
}

func getGormColumnName(s interface{}, fieldName string) (columnName string) {
	t := reflect.TypeOf(s)
	field, _ := t.FieldByName(fieldName)
	tag := field.Tag.Get("gorm")
	columnName = strings.Split(tag, ":")[1]
	columnName = strings.Split(columnName, ";")[0]
	return columnName
}

func getPrimaryKey(s any) (primaryKey string, value reflect.Value) {
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.Contains(field.Tag.Get("gorm"), "primary_key") {
			return getGormColumnName(s, field.Name), reflect.ValueOf(s).Field(i)
		}
		if strings.Contains(field.Tag.Get("gorm"), "primaryKey") {
			return getGormColumnName(s, field.Name), reflect.ValueOf(s).Field(i)
		}
	}
	return "", reflect.Value{}
}

// 获取变化的字段和值
func getChangedFields(old, new any) (changedFields map[string]PropValue) {
	changedFields = make(map[string]PropValue)
	oldValue := reflect.ValueOf(old)
	newValue := reflect.ValueOf(new)
	oldType := oldValue.Type()

	newType := newValue.Type()

	if oldType.Kind() != reflect.Struct || newType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < oldType.NumField(); i++ {
		oldFieldName := oldType.Field(i).Name
		newFieldName := newType.Field(i).Name

		if oldFieldName != newFieldName {
			continue
		}

		if oldFieldName == "Struct" {
			continue
		}

		oldValueField := oldValue.Field(i)
		newValueField := newValue.Field(i)

		// 检查是否为 *time.Time 类型
		if oldValueField.Type().String() == "*time.Time" {
			// 如果两个值都是 nil，则认为相等
			if oldValueField.IsNil() && newValueField.IsNil() {
				continue
			}
			// 如果只有一个值是 nil，则认为不相等
			if oldValueField.IsNil() || newValueField.IsNil() {
				changedFields[oldFieldName] = PropValue{
					Type:  newValueField.Type(),
					Value: newValueField,
				}
				continue
			}
			// 比较时间戳
			oldTime := oldValueField.Elem().Interface().(time.Time)
			newTime := newValueField.Elem().Interface().(time.Time)
			if oldTime.Unix() != newTime.Unix() {
				changedFields[oldFieldName] = PropValue{
					Type:  newValueField.Type(),
					Value: newValueField,
				}
			}
		} else if !reflect.DeepEqual(oldValueField.Interface(), newValueField.Interface()) {
			changedFields[oldFieldName] = PropValue{
				Type:  newValueField.Type(),
				Value: newValueField,
			}
		}
	}

	return changedFields
}
