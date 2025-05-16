package diff

import "github.com/jinzhu/copier"

type Traceable[T any] struct {
	Entity T
	Props  map[string]PropValue
}

func Trace[T any](entity T, fn func(entity *T)) (e T, props map[string]PropValue) {
	fromEntity := entity
	var toEntity T
	copier.Copy(&toEntity, &fromEntity)
	fn(&toEntity)
	props = getChangedFields(fromEntity, toEntity)
	return toEntity, props
}

func TraceValue[T any](entity T, fn func(entity *T)) (t Traceable[T]) {
	newEntity, props := Trace(entity, fn)
	t.Entity = newEntity
	t.Props = props
	return
}
