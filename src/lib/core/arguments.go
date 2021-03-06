package core

// Arguments represents a structured set of arguments passed to a predicate.
type Arguments struct {
	positionals   []*Thunk
	expandedList  *Thunk
	keywords      []KeywordArgument
	expandedDicts []*Thunk
}

// NewArguments creates a new Arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	expandedDicts []*Thunk) Arguments {
	if ks == nil {
		ks = []KeywordArgument{}
	}

	ts := make([]*Thunk, 0, len(ps))

	l := (*Thunk)(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergeRestPositionalArgs(ps[i:]...)
			break
		}

		ts = append(ts, p.value)
	}

	return Arguments{
		positionals:   ts,
		expandedList:  l,
		keywords:      ks,
		expandedDicts: expandedDicts,
	}
}

func mergeRestPositionalArgs(ps ...PositionalArgument) *Thunk {
	if !ps[0].expanded {
		panic("First PositionalArgument must be a list.")
	}

	t := ps[0].value

	for _, p := range ps[1:] {
		if p.expanded {
			t = PApp(Merge, t, p.value)
		} else {
			t = PApp(
				NewLazyFunction(appendFuncSignature, appendFunc), // Avoid initialization loop
				t, p.value)
		}
	}

	return t
}

func (args *Arguments) nextPositional() *Thunk {
	if len(args.positionals) != 0 {
		defer func() { args.positionals = args.positionals[1:] }()
		return args.positionals[0]
	}

	if args.expandedList == nil {
		return nil
	}

	defer func() {
		args.expandedList = PApp(Rest, args.expandedList)
	}()

	return PApp(First, args.expandedList)
}

func (args Arguments) restPositionals() *Thunk {
	if args.expandedList == nil {
		return NewList(args.positionals...)
	}

	if len(args.positionals) == 0 {
		return args.expandedList
	}

	return PApp(Merge, NewList(args.positionals...), args.expandedList)
}

func (args *Arguments) searchKeyword(s string) *Thunk {
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, t := range args.expandedDicts {
		o := t.Eval()
		d, ok := o.(DictionaryType)

		if !ok {
			return NotDictionaryError(o)
		}

		k := StringType(s)

		if v, ok := d.Search(k); ok {
			args.expandedDicts[i] = Normal(d.Remove(k))
			return v.(*Thunk)
		}
	}

	return nil
}

func (args Arguments) restKeywords() *Thunk {
	t := EmptyDictionary

	for _, k := range args.keywords {
		t = PApp(Set, t, NewString(k.name), k.value)
	}

	for _, tt := range args.expandedDicts {
		t = PApp(Merge, t, tt)
	}

	return t
}

func (args Arguments) Merge(merged Arguments) Arguments {
	var new Arguments

	if new.expandedList == nil {
		new.positionals = append(args.positionals, merged.positionals...)
		new.expandedList = merged.expandedList
	} else {
		new.positionals = args.positionals
		new.expandedList = PApp(
			Append,
			append([]*Thunk{args.expandedList}, merged.positionals...)...)

		if merged.expandedList != nil {
			new.expandedList = PApp(Merge, new.expandedList, merged.expandedList)
		}
	}

	new.keywords = append(args.keywords, merged.keywords...)
	new.expandedDicts = append(args.expandedDicts, merged.expandedDicts...)

	return new
}
