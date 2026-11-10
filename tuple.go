package opera

// Tuple holds two typed values in order. Its zero value contains the zero value
// of each element type.
type Tuple[T0, T1 any] struct {
	V0 T0
	V1 T1
}

// Tup creates a tuple with two elements in argument order.
func Tup[T0, T1 any](v0 T0, v1 T1) Tuple[T0, T1] {
	return Tuple[T0, T1]{V0: v0, V1: v1}
}

// Get returns the tuple's elements in order, preserving their types.
// It supports multiple assignment, for example:
//
//	name, age := opera.Tup("Alice", 30).Get()
func (t Tuple[T0, T1]) Get() (T0, T1) {
	return t.V0, t.V1
}

// Tuple3 holds three typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple3[T0, T1, T2 any] struct {
	V0 T0
	V1 T1
	V2 T2
}

// Tup3 creates a tuple with three elements in argument order.
func Tup3[T0, T1, T2 any](v0 T0, v1 T1, v2 T2) Tuple3[T0, T1, T2] {
	return Tuple3[T0, T1, T2]{V0: v0, V1: v1, V2: v2}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple3[T0, T1, T2]) Get() (T0, T1, T2) {
	return t.V0, t.V1, t.V2
}

// Tuple4 holds four typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple4[T0, T1, T2, T3 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
}

// Tup4 creates a tuple with four elements in argument order.
func Tup4[T0, T1, T2, T3 any](v0 T0, v1 T1, v2 T2, v3 T3) Tuple4[T0, T1, T2, T3] {
	return Tuple4[T0, T1, T2, T3]{V0: v0, V1: v1, V2: v2, V3: v3}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple4[T0, T1, T2, T3]) Get() (T0, T1, T2, T3) {
	return t.V0, t.V1, t.V2, t.V3
}

// Tuple5 holds five typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple5[T0, T1, T2, T3, T4 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

// Tup5 creates a tuple with five elements in argument order.
func Tup5[T0, T1, T2, T3, T4 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
) Tuple5[T0, T1, T2, T3, T4] {
	return Tuple5[T0, T1, T2, T3, T4]{V0: v0, V1: v1, V2: v2, V3: v3, V4: v4}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple5[T0, T1, T2, T3, T4]) Get() (T0, T1, T2, T3, T4) {
	return t.V0, t.V1, t.V2, t.V3, t.V4
}

// Tuple6 holds six typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple6[T0, T1, T2, T3, T4, T5 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

// Tup6 creates a tuple with six elements in argument order.
func Tup6[T0, T1, T2, T3, T4, T5 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
) Tuple6[T0, T1, T2, T3, T4, T5] {
	return Tuple6[T0, T1, T2, T3, T4, T5]{V0: v0, V1: v1, V2: v2, V3: v3, V4: v4, V5: v5}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Get() (T0, T1, T2, T3, T4, T5) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5
}

// Tuple7 holds seven typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple7[T0, T1, T2, T3, T4, T5, T6 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
}

// Tup7 creates a tuple with seven elements in argument order.
func Tup7[T0, T1, T2, T3, T4, T5, T6 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
) Tuple7[T0, T1, T2, T3, T4, T5, T6] {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{
		V0: v0,
		V1: v1,
		V2: v2,
		V3: v3,
		V4: v4,
		V5: v5,
		V6: v6,
	}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Get() (T0, T1, T2, T3, T4, T5, T6) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6
}

// Tuple8 holds eight typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple8[T0, T1, T2, T3, T4, T5, T6, T7 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
}

// Tup8 creates a tuple with eight elements in argument order.
func Tup8[T0, T1, T2, T3, T4, T5, T6, T7 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
	v7 T7,
) Tuple8[T0, T1, T2, T3, T4, T5, T6, T7] {
	return Tuple8[T0, T1, T2, T3, T4, T5, T6, T7]{
		V0: v0,
		V1: v1,
		V2: v2,
		V3: v3,
		V4: v4,
		V5: v5,
		V6: v6,
		V7: v7,
	}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple8[T0, T1, T2, T3, T4, T5, T6, T7]) Get() (T0, T1, T2, T3, T4, T5, T6, T7) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6, t.V7
}

// Tuple9 holds nine typed values in order. Its zero value contains the zero
// value of each element type.
type Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
	V8 T8
}

// Tup9 creates a tuple with nine elements in argument order.
func Tup9[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
	v5 T5,
	v6 T6,
	v7 T7,
	v8 T8,
) Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8] {
	return Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8]{
		V0: v0,
		V1: v1,
		V2: v2,
		V3: v3,
		V4: v4,
		V5: v5,
		V6: v6,
		V7: v7,
		V8: v8,
	}
}

// Get returns the tuple's elements in order, preserving their types.
func (t Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8]) Get() (T0, T1, T2, T3, T4, T5, T6, T7, T8) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6, t.V7, t.V8
}
