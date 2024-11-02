package models

// From is a generic interface that defines a conversion from one type to another
type From[T any, U any] interface {
	From(T) U
}

// To is a generic interface that defines a conversion to another type
type To[T any, U any] interface {
	To(U) T
}

// Convert is a generic function that converts from type T to type U
// using a converter that implements From[T, U]
func Convert[T any, U any](value T, converter From[T, U]) U {
	return converter.From(value)
}

// ConvertTo is a generic function that converts to type T from type U
// using a converter that implements To[T, U]
func ConvertTo[T any, U any](value U, converter To[T, U]) T {
	return converter.To(value)
}
