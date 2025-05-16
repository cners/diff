package diff

import "github.com/jinzhu/copier"

type Traceable[T any] struct {
	Entity    T                    // Entity.
	Props     map[string]PropValue // Changed fields.
	UpdateSql string               // Update sql statement.
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
	return
}

func TraceUpdate[T any](entity T, fn func(entity *T)) (t Traceable[T]) {
	newEntity, props := Trace(entity, fn)
	t.Entity = newEntity
	t.Props = props
	t.UpdateSql = BuildUpdateSql(t)
	return
}
