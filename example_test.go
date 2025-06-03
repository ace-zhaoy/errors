package errors_test

import (
	errors2 "errors"
	"fmt"
	"github.com/ace-zhaoy/errors"
)

var (
	causeErr = errors.New("go err")
)

func ExampleNew() {
	err := errors.New("go error")
	fmt.Println(err)
	fmt.Printf("%+v\n", err)

	// Output:
	// go error
	// go error
}

func ExampleNewWithCode() {
	err := errors.NewWithCode(500102030, "go error")
	fmt.Println(err)
	fmt.Printf("%+v", err)

	// Output:
	// [500102030: go error]
	// 500102030: go error
}

func ExampleNewWithStack() {
	err := errors.NewWithStack("go err, param: %v", 123)
	fmt.Println(err)
	fmt.Println("-----------")
	fmt.Printf("%+v", err)

	// Example output:
	// go err, param: 123
	// -----------
	// go err, param: 123
	// github.com/ace-zhaoy/errors_test.ExampleNewWithStack
	//         /go/src/errors/example_test.go:35
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:55
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}

func ExampleWithStack() {
	errNil := errors.WithStack(nil)
	err := func() error {
		return errors.WithStack(causeErr)
	}()
	fmt.Printf("%+v\n", errNil)
	fmt.Println("--------------------")
	fmt.Printf("%+v", err)

	// Example output:
	// <nil>
	// --------------------
	// go err
	// github.com/ace-zhaoy/errors_test.ExampleWithStack.func1
	//         /go/src/errors/example_test.go:64
	// github.com/ace-zhaoy/errors_test.ExampleWithStack
	//         /go/src/errors/example_test.go:65
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:55
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}

func ExampleWithStack_multiple_calls() {
	// multiple calls do not duplicate the trace stack
	err := func() error {
		return errors.WithStack(causeErr)
	}()

	err = func() error {
		return errors.WithStack(err)
	}()

	err = errors.WithStack(err)
	fmt.Printf("%+v\n", err)

	// Example output:
	// go err
	// github.com/ace-zhaoy/errors_test.ExampleWithStack_multiple_calls.func1
	//         /go/src/errors/example_test.go:99
	// github.com/ace-zhaoy/errors_test.ExampleWithStack_multiple_calls
	//         /go/src/errors/example_test.go:100
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:53
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}

func ExampleWrap() {
	err := func() error {
		err := errors2.New("abc")
		return errors.Wrap(err, "param: %v", 123)
	}()

	err = func() error {
		return errors.Wrap(err, "param: %v", 456)
	}()
	fmt.Printf("%+v", err)

	// Example Output:
	// abc
	// param: 123
	// github.com/ace-zhaoy/errors_test.ExampleWrap.func1
	//         /go/src/errors/example_test.go:132
	// github.com/ace-zhaoy/errors_test.ExampleWrap
	//         /go/src/errors/example_test.go:133
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:81
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
	// param: 456
}

func ExampleWrapForce() {
	err := func() error {
		err := errors2.New("abc")
		return errors.WrapForce(err, "param: %v", 123)
	}()

	err = func() error {
		return errors.WrapForce(err, "param: %v", 456)
	}()
	fmt.Printf("%+v", err)

	// Example Output:
	// abc
	// param: 123
	// github.com/ace-zhaoy/errors_test.ExampleWrapForce.func1
	//         /go/src/errors/example_test.go:165
	// github.com/ace-zhaoy/errors_test.ExampleWrapForce
	//         /go/src/errors/example_test.go:166
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:53
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
	// param: 456
	// github.com/ace-zhaoy/errors_test.ExampleWrapForce.func2
	//         /go/src/errors/example_test.go:169
	// github.com/ace-zhaoy/errors_test.ExampleWrapForce
	//         /go/src/errors/example_test.go:170
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:53
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}

func ExampleCause() {
	err := func() error {
		return errors.Wrap(causeErr, "err1")
	}()

	err = func() error {
		return errors.Wrap(err, "err2")
	}()

	causeResErr := errors.Cause(err)

	fmt.Printf("causeErr: %s\n", causeErr)
	fmt.Printf("causeResErr: %s\n", causeResErr)
	fmt.Printf("err: %s\n", err)

	// Output:
	//causeErr:  go err
	//causeResErr:  go err
	//err:  err2 -> {err1 -> {go err}}

	fmt.Printf("%+v ", err)
	// Example Output:
	// go err
	// err1
	// github.com/ace-zhaoy/errors_test.ExampleCause.func1
	//         /go/src/errors/example_test.go:162
	// github.com/ace-zhaoy/errors_test.ExampleCause
	//         /go/src/errors/example_test.go:163
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:53
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
	// err2
}

const (
	CodeInvalidParams  = 400102030
	CodeUserNotFound   = 500201010
	CodeOrderNotExists = 500302001
)

var (
	ErrCodeInvalidParams  = errors.NewWithCode(CodeInvalidParams, "invalid params")
	ErrCodeUserNotFound   = errors.NewWithCode(CodeUserNotFound, "user not found")
	ErrCodeOrderNotExists = errors.NewWithCode(CodeOrderNotExists, "order not exists")
)

var (
	ErrInvalidParams  = errors.NewWithMessage("invalid params")
	ErrUserNotFound   = errors.NewWithMessage("user not found")
	ErrOrderNotExists = errors.NewWithMessage("order not exists")
)

func ExampleLatestCode() {
	sqlQuery := func() error {
		// sql query
		return errors.WithStack(ErrCodeUserNotFound)
	}
	serviceErr := func() error {
		err := sqlQuery()
		return ErrCodeInvalidParams.WrapStack(err)
	}()

	otherErr := errors.Wrap(serviceErr, "service err")

	latestErr := errors.LatestCode(serviceErr)
	fmt.Printf("ErrCodeUserNotFound: %s\n", ErrCodeUserNotFound)
	fmt.Printf("ErrCodeInvalidParams: %s\n", ErrCodeInvalidParams)
	fmt.Printf("latestErr: %s\n", latestErr)
	fmt.Printf("otherErr: %s\n", otherErr)

	// Output:
	//ErrCodeUserNotFound: [500201010: user not found]
	//ErrCodeInvalidParams: [400102030: invalid params]
	//latestErr: [400102030: invalid params] -> {[500201010: user not found]}
	//otherErr: service err -> {[400102030: invalid params] -> {[500201010: user not found]}}
}

func ExampleLatestMessage() {
	sqlQuery := func() error {
		// sql query
		return errors.WithStack(ErrUserNotFound)
	}
	serviceErr := func() error {
		err := sqlQuery()
		return ErrInvalidParams.Wrapf(err, "param <%v>", 123)
	}()

	otherErr := errors.Wrap(serviceErr, "service err")

	latestErr := errors.LatestMessage(serviceErr)
	fmt.Printf("ErrUserNotFound: %s\n", ErrUserNotFound)
	fmt.Printf("ErrInvalidParams: %s\n", ErrInvalidParams)
	fmt.Printf("latestErr: %s\n", latestErr)
	fmt.Printf("otherErr: %s\n", otherErr)

	// Output:
	//ErrUserNotFound: user not found
	//ErrInvalidParams: invalid params
	//latestErr: invalid params -> {param <123> -> {user not found}}
	//otherErr: service err -> {invalid params -> {param <123> -> {user not found}}}
}

func ExampleIs_code() {
	sqlQuery := func() error {
		// sql query
		return ErrCodeUserNotFound.Wrapf(causeErr, "sql err")
	}

	serviceErr := func() error {
		err := sqlQuery()
		return ErrCodeInvalidParams.WrapStack(err)
	}()

	otherErr := errors.Wrap(serviceErr, "service err")

	isErrCodeUserNotFound := errors.Is(otherErr, ErrCodeUserNotFound)
	isErrCodeInvalidParams := errors.Is(otherErr, ErrCodeInvalidParams)
	isErrCodeOrderNotExists := errors.Is(otherErr, ErrCodeOrderNotExists)
	isErrUserNotFound := errors.Is(otherErr, ErrUserNotFound)

	fmt.Printf("isErrCodeUserNotFound: %v\n", isErrCodeUserNotFound)
	fmt.Printf("isErrCodeInvalidParams: %v\n", isErrCodeInvalidParams)
	fmt.Printf("isErrCodeOrderNotExists: %v\n", isErrCodeOrderNotExists)
	fmt.Printf("isErrUserNotFound: %v\n", isErrUserNotFound)

	// Output:
	//isErrCodeUserNotFound: true
	//isErrCodeInvalidParams: true
	//isErrCodeOrderNotExists: false
	//isErrUserNotFound: false
}

func ExampleIs_message() {
	sqlQuery := func() error {
		// sql query
		return ErrUserNotFound.Wrapf(causeErr, "sql err")
	}

	serviceErr := func() error {
		err := sqlQuery()
		return ErrInvalidParams.WrapStack(err)
	}()

	otherErr := errors.Wrap(serviceErr, "service err")

	isErrUserNotFound := errors.Is(otherErr, ErrUserNotFound)
	isErrInvalidParams := errors.Is(otherErr, ErrInvalidParams)
	isErrOrderNotExists := errors.Is(otherErr, ErrOrderNotExists)
	isErrCodeUserNotFound := errors.Is(otherErr, ErrCodeUserNotFound)

	fmt.Printf("isErrUserNotFound: %v\n", isErrUserNotFound)
	fmt.Printf("isErrInvalidParams: %v\n", isErrInvalidParams)
	fmt.Printf("isErrOrderNotExists: %v\n", isErrOrderNotExists)
	fmt.Printf("isErrCodeUserNotFound: %v\n", isErrCodeUserNotFound)

	// Output:
	//isErrUserNotFound: true
	//isErrInvalidParams: true
	//isErrOrderNotExists: false
	//isErrCodeUserNotFound: false
}

func checkA1() error {
	return nil
}

func checkRecover1() (err error) {
	defer errors.Recover(func(e error) {
		err = e
	})

	err = checkA1()

	errors.Check(err)
	// output
	fmt.Println("123")
	return nil
}

func ExampleCheck_nil() {
	err := checkRecover1()
	fmt.Printf("%+v", err)

	// Output:
	//123
	//<nil>
}

func checkA2() error {
	return errors.WithStack(ErrCodeUserNotFound)
}

func checkRecover2() (err error) {
	defer errors.Recover(func(e error) {
		err = e
	})
	err = checkA2()
	errors.Check(err)
	// no output
	fmt.Println("123")
	return nil
}

func ExampleCheck_error() {
	err := checkRecover2()
	fmt.Printf("%+v", err)

	// Example Output:
	// 500201010: user not found
	// github.com/ace-zhaoy/errors_test.checkA2
	//         /go/src/errors/example_test.go:403
	// github.com/ace-zhaoy/errors_test.checkRecover2
	//         /go/src/errors/example_test.go:410
	// github.com/ace-zhaoy/errors_test.ExampleCheck_error
	//         /go/src/errors/example_test.go:418
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:85
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}

func checkRecover3() (err error) {
	defer errors.Recover(func(e error) {
		err = e
	})
	panic("test panic")

	return nil
}

func ExampleCheck_panic() {
	err := checkRecover3()
	fmt.Printf("%+v", err)

	// Example Output:
	// test panic
	// runtime.gopanic
	//         /usr/local/go/src/runtime/panic.go:838
	// github.com/ace-zhaoy/errors_test.checkRecover3
	//         /go/src/errors/example_test.go:447
	// github.com/ace-zhaoy/errors_test.ExampleCheck_panic
	//         /go/src/errors/example_test.go:453
	// testing.runExample
	//         /usr/local/go/src/testing/run_example.go:63
	// testing.runExamples
	//         /usr/local/go/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/go/src/testing/testing.go:1721
	// main.main
	//         _testmain.go:83
	// runtime.main
	//         /usr/local/go/src/runtime/proc.go:250
	// runtime.goexit
	//         /usr/local/go/src/runtime/asm_amd64.s:1571
}
