package mock

// Expectations are a mapping from function name to Expectations
type Expectations map[string]*Expectation

// Expectation depicts expectations for multiple calls of a function
type Expectation struct {
	CallCount     int
	DefaultReturn interface{}
	ReturnValues  []interface{}
}

// Add adds an expectation
func (e Expectations) Add(fn string, def interface{}, retValues ...interface{}) {
	exp := e[fn]
	if exp == nil {
		exp = &Expectation{}
	}
	exp.DefaultReturn = def
	exp.ReturnValues = append(exp.ReturnValues, retValues...)
	e[fn] = exp
}

// Return returns an expecatation's Return function result
func (e Expectations) Return(fn string) (interface{}, error) {
	return e[fn].Return()
}

// Return returns a return value and increases the callcount
func (e *Expectation) Return() (interface{}, error) {
	res := e.ReturnValues[e.CallCount]
	var err error
	errTyped, ok := res.(error)
	if ok {
		err = errTyped
		res = e.DefaultReturn
	}
	e.CallCount = e.CallCount + 1
	return res, err
}
