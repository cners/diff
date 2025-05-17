package diff

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func sqlChecking(sql string) bool {
	return true
}

// 构建更新sql,如果sql不合法则返回空字符串
func buildUpdateSql[T any](changedFields map[string]PropValue, entity T) (sql string) {
	if len(changedFields) == 0 {
		return ""
	}
	// 构建sql
	sql = "UPDATE \"{tableName}\" SET "
	values := []string{}
	for k, v := range changedFields {
		columnName := getGormColumnName(entity, k)
		values = append(values, fmt.Sprintf("\"%s\" = %v", columnName, getVal(v)))
	}
	sql += strings.Join(values, ",")

	// 替换表明
	sql = strings.Replace(sql, "{tableName}", getTableName(&entity), 1)
	sql += buildWherePrimary(entity)

	// 检查sql是否合法
	if !sqlChecking(sql) {
		sql = ""
	}
	return sql
}

func getVal(v PropValue) string {
	switch v.Type.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%v", v.Value.Interface())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", v.Value.Interface())
	case reflect.Ptr:
		if v.Type.Elem().String() == "time.Time" {
			if v.Value.IsNil() {
				return "NULL"
			}
			t := v.Value.Elem().Interface().(time.Time)
			return fmt.Sprintf("'%s'", t.Format(time.RFC3339))
		}
		return fmt.Sprintf("'%v'", v.Value.Interface())
	case reflect.String:
		return fmt.Sprintf("'%v'", v.Value.Interface())
	default:
		return fmt.Sprintf("'%v'", v.Value.Interface())
	}
}

func buildWherePrimary(entity any) (sql string) {
	primaryKey, val := getPrimaryKey(entity)
	if primaryKey == "" {
		return ""
	}
	return fmt.Sprintf(" WHERE \"%s\" = %v", primaryKey, val)
}

func BuildUpdateSql[T any](t Traceable[T]) (sql string) {
	return buildUpdateSql(t.Props, t.Entity)
}

func getTableName[T any](model T) (tableName string) {
	// 获取值的反射对象
	val := reflect.ValueOf(model)

	// 处理指针类型
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 查找TableName方法
	method := val.MethodByName("TableName")
	if !method.IsValid() {
		// 如果没有找到方法，尝试获取类型的方法
		t := reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		method = reflect.New(t).MethodByName("TableName")
	}

	if !method.IsValid() {
		return "" // 没有定义TableName方法
	}

	// 调用方法并获取结果
	result := method.Call(nil) // 无参数调用
	if len(result) > 0 {
		return result[0].String() // 返回第一个结果的字符串形式
	}

	return ""
}
