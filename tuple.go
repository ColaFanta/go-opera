package opera

type Tuple[T0, T1 any] = struct {
	V0 T0
	V1 T1
}

// Tup creates a tuple of with 2 elements.
func Tup[T0, T1 any](v0 T0, v1 T1) Tuple[T0, T1] {
	return Tuple[T0, T1]{V0: v0, V1: v1}
}

type Tuple3[T0, T1, T2 any] = struct {
	V0 T0
	V1 T1
	V2 T2
}

// Tup3 creates a tuple of with 3 elements.
func Tup3[T0, T1, T2 any](v0 T0, v1 T1, v2 T2) Tuple3[T0, T1, T2] {
	return Tuple3[T0, T1, T2]{V0: v0, V1: v1, V2: v2}
}

type Tuple4[T0, T1, T2, T3 any] = struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
}

// Tup4 creates a tuple of with 4 elements.
func Tup4[T0, T1, T2, T3 any](v0 T0, v1 T1, v2 T2, v3 T3) Tuple4[T0, T1, T2, T3] {
	return Tuple4[T0, T1, T2, T3]{V0: v0, V1: v1, V2: v2, V3: v3}
}

type Tuple5[T0, T1, T2, T3, T4 any] = struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

// Tup5 creates a tuple of with 5 elements.
func Tup5[T0, T1, T2, T3, T4 any](
	v0 T0,
	v1 T1,
	v2 T2,
	v3 T3,
	v4 T4,
) Tuple5[T0, T1, T2, T3, T4] {
	return Tuple5[T0, T1, T2, T3, T4]{V0: v0, V1: v1, V2: v2, V3: v3, V4: v4}
}

type Tuple6[T0, T1, T2, T3, T4, T5 any] = struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

// Tup6 creates a tuple of with 6 elements.
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

type Tuple7[T0, T1, T2, T3, T4, T5, T6 any] = struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
}

// Tup7 creates a tuple of with 7 elements.
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

type Tuple8[T0, T1, T2, T3, T4, T5, T6, T7 any] = struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
	V7 T7
}

// Tup8 creates a tuple of with 8 elements.
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

type Tuple9[T0, T1, T2, T3, T4, T5, T6, T7, T8 any] = struct {
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

// Tup9 creates a tuple of with 9 elements.
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
