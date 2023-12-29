# Errors
Package errors Adds stack trace support for errors in go, along with support for error code and carry messages

# Install
```shell
go get github.com/ace-zhaoy/errors
```

# Usage
## stack trace
If err is nil, WithStack returns nil.
#### Example 1
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"os"
)

func OpenFile() error {
	_, err := os.Open("./test.err")
	return errors.WithStack(err)
}

func ReturnNil() error {
	return errors.WithStack(nil)
}

func main() {
	err := OpenFile()
	fmt.Printf("%+v\n", err)
	fmt.Println("------------")

	err = ReturnNil()
	fmt.Printf("%+v\n", err)
}

```
Output:
```shell
open ./test.err: no such file or directory
main.OpenFile
        /go/src/test/err/main.go:11
main.main
        /go/src/test/err/main.go:19
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1571
------------
<nil>

```

#### Example 2
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"os"
)

func A() error {
	_, err := os.Open("./test.err")
	return errors.WithStack(err)
}

func B() error {
	err := A()
	return errors.WithStack(err)
}

func C() error {
	err := B()
	return errors.WithStack(err)
}

func main() {
	err := C()
	fmt.Printf("%+v\n", err)
}

```

Output:
```
open ./test.err: no such file or directory
main.A
        /go/src/test/err/main.go:11
main.B
        /go/src/test/err/main.go:15
main.C
        /go/src/test/err/main.go:20
main.main
        /go/src/test/err/main.go:25
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:157
```

## Error code
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
)

const (
	CodeInvalidToken   = 401100000
	CodeInvalidParams  = 400102030
	CodeUserNotFound   = 500201010
	CodeOrderNotExists = 500302001
)

var (
	ErrCodeInvalidToken   = errors.NewWithCode(CodeInvalidToken, "invalid token")
	ErrCodeInvalidParams  = errors.NewWithCode(CodeInvalidParams, "invalid params")
	ErrCodeUserNotFound   = errors.NewWithCode(CodeUserNotFound, "user not found")
	ErrCodeOrderNotExists = errors.NewWithCode(CodeOrderNotExists, "order not exists")
)

func SqlFirst() error {
	// record not found
	// err := first()
	err := errors.New("record not found")

	return ErrCodeUserNotFound.Wrapf(err, "param: id: %d", 1)
}

func Server() error {
	err := SqlFirst()
	if err != nil {
		if errors.Is(err, ErrCodeUserNotFound) {
			err = ErrCodeInvalidParams.WrapStack(err)
		}
		return err
	}
	// other code
	return nil
}

func main() {
	err := Server()
	codeErr := errors.LatestCode(err)
	fmt.Printf("error code: %d\n", codeErr.Code())
	fmt.Printf("error msg: %s\n", codeErr.Message())
	fmt.Printf("error err: %s\n", codeErr)
	fmt.Println("---------------------")
	fmt.Printf("%+v\n", err)
}

```
Output:
```
error code: 400102030
error msg: invalid params
error err: [400102030: invalid params] -> {param: id: 1 -> {[500201010: user not found] -> {record not found}}}
---------------------
record not found
500201010: user not found
param: id: 1
main.SqlFirst
        /go/src/test/err/main.go:27
main.Server
        /go/src/test/err/main.go:31
main.main
        /go/src/test/err/main.go:43
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1571
400102030: invalid params
```
tips: 
> 500102030 \
> 500: http cede  \
> 10: level1 code \
> 20: level2 code \
> 30: level3 code


## Error Message
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
)

var (
	ErrInvalidParams = errors.NewWithMessage("invalid params")
)

func A() error {
	return errors.WithStack(ErrInvalidParams)
}

func B() error {
	err := A()
	return errors.Wrap(err, "param: %v", 11)
}

func main() {
	err := B()
	fmt.Printf("%s\n", err)
	fmt.Println("------------")
	fmt.Printf("%+v\n", err)
}

```
Output:
```shell
param: 11 -> {invalid params}
------------
invalid params
main.A
        /go/src/test/err/main.go:13
main.B
        /go/src/test/err/main.go:17
main.main
        /go/src/test/err/main.go:22
runtime.main
        /usr/local/go/src/runtime/proc.go:250
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1571
param: 11
```

## Check Error
It used to be that way
```go
func A() error {
	// ...
    if err != nil {
        return err
    }
}
```

Now
```go
func A() (err error) {
    defer errors.Recover(func(e error) {
		err = e
    })
	
	err := B()
    errors.Check(err)
    // If err is not nil, the code that follows will not execute
    // ...
	return nil
}
```