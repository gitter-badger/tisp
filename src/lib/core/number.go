package core

import "math"

type NumberType float64

func NewNumber(n float64) *Thunk {
	return Normal(NumberType(n))
}

func (n NumberType) equal(e equalable) Object {
	return rawBool(n == e.(NumberType))
}

var Add = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "nums",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		l, ok := o.(ListType)

		if !ok {
			return notListError(o)
		}

		os, err := l.ToObjects()

		if err != nil {
			return err
		}

		sum := NumberType(0)

		for _, o := range os {
			n, ok := o.(NumberType)

			if !ok {
				return notNumberError(o)
			}

			sum += n
		}

		return sum
	})

var Sub = NewLazyFunction(
	NewSignature(
		[]string{"minuend"}, []OptionalArgument{}, "subtrahends",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		n0, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		o = ts[1].Eval()
		l, ok := o.(ListType)

		if !ok {
			return notListError(o)
		}

		os, err := l.ToObjects()

		if err != nil {
			return err
		}

		if len(os) == 0 {
			return NumArgsError("sub", ">= 1")
		}

		for _, o := range os {
			n, ok := o.(NumberType)

			if !ok {
				return notNumberError(o)
			}

			n0 -= n
		}

		return n0
	})

var Mul = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "nums",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		l, ok := o.(ListType)

		if !ok {
			return notListError(o)
		}

		os, err := l.ToObjects()

		if err != nil {
			return err
		}

		prod := NumberType(1)

		for _, o := range os {
			n, ok := o.(NumberType)

			if !ok {
				return notNumberError(o)
			}

			prod *= n
		}

		return prod
	})

var Div = NewLazyFunction(
	NewSignature(
		[]string{"dividend"}, []OptionalArgument{}, "divisors",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Object {
		o := ts[0].Eval()
		n0, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		o = ts[1].Eval()
		l, ok := o.(ListType)

		if !ok {
			return notListError(o)
		}

		os, err := l.ToObjects()

		if err != nil {
			return err
		}

		if len(os) == 0 {
			return NumArgsError("div", ">= 1")
		}

		for _, o := range os {
			n, ok := o.(NumberType)

			if !ok {
				return notNumberError(o)
			}

			n0 /= n
		}

		return n0
	})

var Mod = NewStrictFunction(
	NewSignature(
		[]string{"dividend", "divisor"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		if len(os) != 2 {
			return NumArgsError("mod", "2")
		}

		o := os[0]
		n1, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		o = os[1]
		n2, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		return NewNumber(math.Mod(float64(n1), float64(n2)))
	})

var Pow = NewStrictFunction(
	NewSignature(
		[]string{"base", "exponent"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(os ...Object) Object {
		o := os[0]
		n1, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		o = os[1]
		n2, ok := o.(NumberType)

		if !ok {
			return notNumberError(o)
		}

		return NewNumber(math.Pow(float64(n1), float64(n2)))
	})

func notNumberError(o Object) *Thunk {
	return TypeError(o, "Number")
}

// ordered

func (n NumberType) less(o ordered) bool {
	return n < o.(NumberType)
}
