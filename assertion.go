package mapq

import (
	"strconv"
	"strings"
)

type Operation func(left any, expected ...any) string

type assertion struct {
	property  string
	operation Operation
	expected  []any
}

func (a assertion) And(a2 any) joiner {
	return joiner{andJoin, []any{a, a2}}
}

func (a assertion) Or(a2 any) joiner {
	return joiner{orJoin, []any{a, a2}}
}

func (a assertion) XOr(a2 any) joiner {
	return joiner{xOrJoin, []any{a, a2}}
}

func (a assertion) compute(m map[string]any) bool {
	props := strings.Split(a.property, ".")

	var obj any = m

	getProp := func(p string) any {
		switch v := obj.(type) {
		case map[string]any:
			return v[p]
		case []any:
			pAsInt, err := strconv.ParseInt(p, 10, 64)

			if err != nil {
				return nil
			}

			return v[pAsInt]
		}
		return nil
	}

	var i int
	var p string
	for i, p = range props {
		if i == len(props)-1 {
			continue
		}

		obj = getProp(p)
		if obj == nil {
			return false
		}
	}

	return a.operation(getProp(p), a.expected...) == ""
}
