package domain

// Nullable описывает опциональное поле в PATCH-запросе.
//
// Set=false — поле не передано, значение не меняем.
// Set=true, Value=nil — поле явно null, сбрасываем (например, phone_number).
// Set=true, Value=&v — поле передано, обновляем до v.
//
// Обычный *T не различает «не передано» и «передано null».
type Nullable[T any] struct {
	Value *T
	Set   bool
}