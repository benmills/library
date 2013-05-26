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

func (expect *expectation) ToEqual(value interface{}) {
	if expect.target != value {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to equal %s.\n  %s:%d\n", value, expect.target, file, line)
	}
}

func (expect *expectation) ToBeTrue() {
	if expect.target != true {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to be true.\n  %s:%d\n", expect.target, file, line)
	}
}

func (expect *expectation) ToBeFalse() {
	if expect.target != false {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to be false.\n  %s:%d\n", expect.target, file, line)
	}
}

func (expect *expectation) ToBeLessThan(value int) {
	intTarget := expect.target.(int)
	if intTarget > value {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to be less than %s.\n  %s:%d\n", expect.target, value, file, line)
	}
}

func (expect *expectation) ToBeGreaterThan(value int) {
	intTarget := expect.target.(int)
	if intTarget < value {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to be greater than %s.\n  %s:%d\n", expect.target, value, file, line)
	}
}

func (expect *expectation) ToContain(value string) {
	stringTarget := expect.target.(string)
	if !strings.Contains(stringTarget, value) {
		_, file, line, _ := runtime.Caller(1)
		expect.t.Fail()
		fmt.Printf("Expected %s to contain %s.\n  %s:%d\n", expect.target, value, file, line)
	}
}
