package opera

import "context"

// AwaitZip awaits 2 channels and returns their results as a tuple.
func AwaitZip[T1, T2 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
) Result[Tuple[T1, T2]] {
	res := AwaitAll[any](ctx, ch1, ch2)
	return Do(func() Tuple[T1, T2] {
		v := res.Yield()
		return Tup(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
		)
	})
}

// AwaitZip3 awaits 3 channels and returns their results as a tuple.
func AwaitZip3[T1, T2, T3 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
) Result[Tuple3[T1, T2, T3]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3)
	return Do(func() Tuple3[T1, T2, T3] {
		v := res.Yield()
		return Tup3(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
		)
	})
}

// AwaitZip4 awaits 4 channels and returns their results as a tuple.
func AwaitZip4[T1, T2, T3, T4 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
) Result[Tuple4[T1, T2, T3, T4]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4)
	return Do(func() Tuple4[T1, T2, T3, T4] {
		v := res.Yield()
		return Tup4(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
		)
	})
}

// AwaitZip5 awaits 5 channels and returns their results as a tuple.
func AwaitZip5[T1, T2, T3, T4, T5 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
	ch5 <-chan Result[T5],
) Result[Tuple5[T1, T2, T3, T4, T5]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4, ch5)
	return Do(func() Tuple5[T1, T2, T3, T4, T5] {
		v := res.Yield()
		return Tup5(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
			MayCast[T5](v[4]).OrEmpty(),
		)
	})
}

// AwaitZip6 awaits 6 channels and returns their results as a tuple.
func AwaitZip6[T1, T2, T3, T4, T5, T6 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
	ch5 <-chan Result[T5],
	ch6 <-chan Result[T6],
) Result[Tuple6[T1, T2, T3, T4, T5, T6]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4, ch5, ch6)
	return Do(func() Tuple6[T1, T2, T3, T4, T5, T6] {
		v := res.Yield()
		return Tup6(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
			MayCast[T5](v[4]).OrEmpty(),
			MayCast[T6](v[5]).OrEmpty(),
		)
	})
}

// AwaitZip7 awaits 7 channels and returns their results as a tuple.
func AwaitZip7[T1, T2, T3, T4, T5, T6, T7 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
	ch5 <-chan Result[T5],
	ch6 <-chan Result[T6],
	ch7 <-chan Result[T7],
) Result[Tuple7[T1, T2, T3, T4, T5, T6, T7]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4, ch5, ch6, ch7)
	return Do(func() Tuple7[T1, T2, T3, T4, T5, T6, T7] {
		v := res.Yield()
		return Tup7(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
			MayCast[T5](v[4]).OrEmpty(),
			MayCast[T6](v[5]).OrEmpty(),
			MayCast[T7](v[6]).OrEmpty(),
		)
	})
}

// AwaitZip8 awaits 8 channels and returns their results as a tuple.
func AwaitZip8[T1, T2, T3, T4, T5, T6, T7, T8 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
	ch5 <-chan Result[T5],
	ch6 <-chan Result[T6],
	ch7 <-chan Result[T7],
	ch8 <-chan Result[T8],
) Result[Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4, ch5, ch6, ch7, ch8)
	return Do(func() Tuple8[T1, T2, T3, T4, T5, T6, T7, T8] {
		v := res.Yield()
		return Tup8(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
			MayCast[T5](v[4]).OrEmpty(),
			MayCast[T6](v[5]).OrEmpty(),
			MayCast[T7](v[6]).OrEmpty(),
			MayCast[T8](v[7]).OrEmpty(),
		)
	})
}

// AwaitZip9 awaits 9 channels and returns their results as a tuple.
func AwaitZip9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
	ctx context.Context,
	ch1 <-chan Result[T1],
	ch2 <-chan Result[T2],
	ch3 <-chan Result[T3],
	ch4 <-chan Result[T4],
	ch5 <-chan Result[T5],
	ch6 <-chan Result[T6],
	ch7 <-chan Result[T7],
	ch8 <-chan Result[T8],
	ch9 <-chan Result[T9],
) Result[Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]] {
	res := AwaitAll[any](ctx, ch1, ch2, ch3, ch4, ch5, ch6, ch7, ch8, ch9)
	return Do(func() Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
		v := res.Yield()
		return Tup9(
			MayCast[T1](v[0]).OrEmpty(),
			MayCast[T2](v[1]).OrEmpty(),
			MayCast[T3](v[2]).OrEmpty(),
			MayCast[T4](v[3]).OrEmpty(),
			MayCast[T5](v[4]).OrEmpty(),
			MayCast[T6](v[5]).OrEmpty(),
			MayCast[T7](v[6]).OrEmpty(),
			MayCast[T8](v[7]).OrEmpty(),
			MayCast[T9](v[8]).OrEmpty(),
		)
	})
}
