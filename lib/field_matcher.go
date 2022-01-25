package quamina

// fieldMatchState has a map which is keyed by the field pathSegments values that can start transitions from this
//  state; for eacrh such field, there is a valueMatcher which, given the field's value, determines whether
//  the automaton progresses to another fieldMatchState
// matches contains the X values that arrival at this state implies have matched
// existsFalseFailures reports the condition that traversal has occurred by matching a field which is named in an
//  exists:false pattern, and the named X's should be subtracted from the matches list being built up by a match project
// The matches field contains a list of the patterns that have been matched if traversal arrives at this state
type fieldMatchState struct {
	transitions         map[string]*valueMatchState
	matches             []X
	existsFalseFailures *matchSet
}

func newFieldMatchState() *fieldMatchState {
	return &fieldMatchState{transitions: make(map[string]*valueMatchState), existsFalseFailures: newMatchSet()}
}

func (m *fieldMatchState) addTransition(field *patternField) []*fieldMatchState {

	// transition from a fieldMatchstate might already be present; create a new empty one if not
	valueMatcher, ok := m.transitions[field.path]
	if !ok {
		valueMatcher = newValueMatchState()
		m.transitions[field.path] = valueMatcher
	}

	// suppose I'm adding the first pattern to a matcher and it has "x": [1, 2]. In principle the branches on
	//  "x": 1 and "x": 2 could go to tne same next state. But we have to make a unique next state for each of them
	//  because some future other pattern might have "x": [2, 3] and thus we need a separate branch to potentially
	//  match two patterns on "x": 2 but not "x": 1. If you were optimizing the automaton for size you might detect
	//  cases where this doesn't happen and reduce the number of fieldMatchstates

	var nextFieldMatchers []*fieldMatchState
	for _, val := range field.vals {
		nextFieldMatchers = append(nextFieldMatchers, valueMatcher.addTransition(val))
	}
	return nextFieldMatchers
}

// transitionOn returns one or more fieldMatchStates you can transition to on a field's name/value combination,
//  or nil if no trnasitions are possible.
func (m *fieldMatchState) transitionOn(field *Field) []*fieldMatchState {

	// are there transitions on this field name?
	valMatcher, ok := m.transitions[string(field.Path)]
	if !ok {
		return nil
	}

	return valMatcher.transitionOn(field.Val)
}

/* for debugging
func (m *fieldMatchState) String() string {
	var keys []string
	for k := range m.transitions {
		p := strings.ReplaceAll(k, "\n", "**")
		keys = append(keys, p)
	}
	keys = append(keys, fmt.Sprintf(" Matches: %d", len(m.matches)))
	return strings.Join(keys, " / ")
}
*/