package diff

import (
	"slices"
	"time"

	"github.com/jinzhu/copier"
)

type Traceable[T any] struct {
	Entity    T                    // Entity.
	Props     map[string]PropValue // Changed fields.
	Columns   map[string]interface{}
	IsChanged bool   // Whether the field is changed.
	UpdateSql string // Update sql statement.
}

func Trace[T any](entity T, fn func(entity *T)) (e T, props map[string]PropValue) {
	fromEntity := entity
	var toEntity T
	copier.Copy(&toEntity, &fromEntity)
	fn(&toEntity)
	props = getChangedFields(fromEntity, toEntity)
	return toEntity, props
}

func TraceProps[T any](entity T, fn func(entity *T)) (t Traceable[T]) {
	newEntity, props := Trace(entity, fn)
	t.Entity = newEntity
	t.Props = props
	t.Columns = parsePropsToPostgresColumns(props)
	t.IsChanged = IsChanged(props)
	return
}

func TraceUpdate[T any](entity T, fn func(entity *T)) (t Traceable[T]) {
	t = TraceProps(entity, fn)
	if !t.IsChanged {
		return
	}
	t.UpdateSql = BuildUpdateSql(t)
	return
}

func CopyValues[T any](fromEntity T, toEntity *T) {
	copier.Copy(toEntity, &fromEntity)
}

func UTC() *time.Time {
	now := time.Now().UTC()
	return &now
}

func parsePropsToPostgresColumns(props map[string]PropValue) (columns map[string]interface{}) {
	columns = make(map[string]interface{})
	for _, v := range props {
		columns[v.ColumnName] = v.Value.Interface()
	}
	return
}

var IgnoreName = []string{
	"updated_at",
	"updated_by",
	"created_at",
	"created_by",
}

// 结构体中是否有字段被修改。这里会排除 IgnoreName 中的字段
func IsChanged(props map[string]PropValue) (isChanged bool) {
	changedCount := 0
	for _, v := range props {
		if !slices.Contains(IgnoreName, v.ColumnName) {
			changedCount++
		}
	}
	return changedCount > 0
}
