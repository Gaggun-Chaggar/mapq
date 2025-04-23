package mapq

type join string

const (
	andJoin join = "and"
	orJoin  join = "or"
	xOrJoin join = "xor"
)

type joiner struct {
	joinType join
	stack    []any
}

// Shorthand to chain two joiners with "And"
func (j joiner) And(a2 any) joiner {
	return joiner{andJoin, []any{j, a2}}
}

// Shorthand to chain two joiners with "Or"
func (j joiner) Or(a2 any) joiner {
	return joiner{orJoin, []any{j, a2}}
}

// Shorthand to chain two joiners with "XOr"
func (j joiner) XOr(a2 any) joiner {
	return joiner{xOrJoin, []any{j, a2}}
}

func (j joiner) collect(m map[string]any) bool {

	var assertionResults []bool

	for _, a := range j.stack {
		switch v := a.(type) {
		case assertion:
			assertionResults = append(assertionResults, v.compute(m))
		case joiner:
			assertionResults = append(assertionResults, v.collect(m))
		}
	}

	return collectFuncs[j.joinType](assertionResults)
}

type collectResult func(results []bool) bool

var collectFuncs = map[join]collectResult{
	andJoin: collectAndResult,
	orJoin:  collectOrResult,
	xOrJoin: collectXOrResult,
}

func collectAndResult(results []bool) bool {
	for _, r := range results {
		if !r {
			return false
		}
	}
	return true
}

func collectOrResult(results []bool) bool {
	for _, r := range results {
		if r {
			return true
		}
	}
	return false
}

func collectXOrResult(results []bool) bool {
	trueCount := 0
	for _, r := range results {
		if r {
			trueCount++
		}
	}
	return trueCount == 1
}
