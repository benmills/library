package quiz

import (
	"testing"
	"fmt"
	"runtime"
	"strings"
)

func Test(t *testing.T) *tester {
	return &tester{t}
}

type tester struct {
	*testing.T
}

func (t *tester) Expect(target interface{}) *expectation {
	return &expectation{t: t, target: target}
}

type expectation struct {
	t *tester
	target interface{}
}

type assertion struct {
	failure bool
	failureMessage string
	messageParts []interface{}
	expect *expectation
}

func (a assertion) eval(expect *expectation) {
	if a.failure {
		_, file, line, _ := runtime.Caller(2)
		expect.t.Fail()
		fmt.Printf(a.failureMessage+"\n  %s:%d\n", append(a.messageParts, file, line)...)
	}
}

func (expect *expectation) ToEqual(value interface{}) {
	assertion{
		failure: expect.target != value,
		failureMessage: "Expected %s to equal %s.",
		messageParts: []interface{}{value, expect.target},
	}.eval(expect)
}

func (expect *expectation) ToBeTrue() {
	assertion{
		failure: expect.target != true,
		failureMessage: "Expected %s to be true.",
		messageParts: []interface{}{expect.target},
	}.eval(expect)
}

func (expect *expectation) ToBeFalse() {
	assertion{
		failure: expect.target != false,
		failureMessage: "Expected %s to be false.",
		messageParts: []interface{}{expect.target},
	}.eval(expect)
}

func (expect *expectation) ToBeLessThan(value int) {
	intTarget := expect.target.(int)
	assertion{
		failure: intTarget > value,
		failureMessage: "Expected %s to be less tahn %s.",
		messageParts: []interface{}{expect.target, value},
	}.eval(expect)
}

func (expect *expectation) ToBeGreaterThan(value int) {
	intTarget := expect.target.(int)
	assertion{
		failure: intTarget < value,
		failureMessage: "Expected %s to be greater than %s.",
		messageParts: []interface{}{expect.target, value},
	}.eval(expect)
}

func (expect *expectation) ToContain(value string) {
	stringTarget := expect.target.(string)
	assertion{
		failure: !strings.Contains(stringTarget, value),
		failureMessage: "Expected %s to contain %s.",
		messageParts: []interface{}{expect.target, value},
	}.eval(expect)
}
