# Errors

Package errors 添加了对go中的错误的堆栈跟踪支持，以及对 error code 和 message 的支持

# 最佳实践
1. 返回的错误应该使用官方的 `error` 接口
```go
type error interface {
	Error() string
}
```

2. 如果不清楚接收到的 error 是否有堆栈跟踪，可以使用 `errors.WithStack()` 进行包装<br>
该方法会检测是否有堆栈跟踪，且不会重复增加堆栈跟踪<br>
也可以使用 `errors.Wrap()`，在增加堆栈跟踪的同时增加 message 消息
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"os"
)

func OpenFile() error {
	_, err := os.Open("./test.err")
	// 无需判断 nil，自动添加堆栈
	return errors.WithStack(err)
}

func main() {
	err := OpenFile()
	// 打印堆栈信息
	fmt.Printf("%+v\n", err)
}
```

## ErrorMessage
不返回给客户端的错误，通常用于日志记录，可以使用 `ErrorMessage` 
#### 方式一
使用`errors.New()`、`errors.NewWithMessage()` 预定义错误，<br>
在调用处使用`errors.WithStack()`调用
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"os"
)

var (
	ErrInvalidParams = errors.NewWithMessage("invalid params")
)

func OpenFile() error {
	filePath := "./test.err"
	_, err := os.Open(filePath)
	if err != nil {
		err = ErrInvalidParams.Wrapf(err, "param: %s", filePath)
	}
	return err
}

func main() {
	err := OpenFile()
	// 打印堆栈信息
	fmt.Printf("%+v\n", err)
}

// Output:
// open ./test.err: no such file or directory
// invalid params
// param: ./test.err
// main.OpenFile
//         /go/src/test/err/main.go:17
// main.main
//         /go/src/test/err/main.go:23
// runtime.main
//         /usr/local/go/src/runtime/proc.go:250
// runtime.goexit
//         /usr/local/go/src/runtime/asm_amd64.s:1571
```

#### 方式二
直接在调用处使用 `errors.NewWithStack` 
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
)

func Service() error {
	err := errors.NewWithStack("sql err: %s", "sql 123")
	return err
}

func main() {
	err := Service()
	// 打印堆栈信息
	fmt.Printf("%+v\n", err)
}

// Output:
// sql err: sql 123
// main.Service
//         /go/src/test/err/main.go:9
// main.main
//         /go/src/test/err/main.go:14
// runtime.main
//         /usr/local/go/src/runtime/proc.go:250
// runtime.goexit
//         /usr/local/go/src/runtime/asm_amd64.s:1571
```

## ErrorCode
通常用于识别特殊错误，或者给客户端返回特殊错误码 <br>
一般先预定义，在需要的地方调用。需要注意，预定义没有堆栈跟踪，需要使用
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

// Output:
// error code: 400102030
// error msg: invalid params
// error err: [400102030: invalid params] -> {param: id: 1 -> {[500201010: user not found] -> {record not found}}}
// ---------------------
// record not found
// 500201010: user not found
// param: id: 1
// main.SqlFirst
//         /go/src/test/err/main.go:27
// main.Server
//         /go/src/test/err/main.go:31
// main.main
//         /go/src/test/err/main.go:43
// runtime.main
//         /usr/local/go/src/runtime/proc.go:250
// runtime.goexit
//         /usr/local/go/src/runtime/asm_amd64.s:1571
// 400102030: invalid params
```