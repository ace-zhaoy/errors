# Errors

Package errors 添加了对go中的错误的堆栈跟踪支持，以及对 error code 和 message 的支持

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

# 最佳实践
### 1. 使用官方的 `error` 接口
返回的错误应该使用官方的 `error` 接口
```go
type error interface {
	Error() string
}
```

### 2. 对外部包返回的 error 增加堆栈跟踪
如果不清楚接收到的 error 是否有堆栈跟踪，可以使用 `errors.WithStack()` 进行包装<br>
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

func Service() error {
	err := OpenFile()
	if err != nil {
		return err
	}
	// ...
	return nil
}

func main() {
	err := Service()
	// 打印堆栈信息
	fmt.Printf("%+v\n", err)
}
```

### 3. 屏蔽项目外部 error
比如，`github.com/jinzhu/gorm` 未找到记录时会返回 `gorm.ErrRecordNotFound`， <br>
`github.com/redis/go-redis/v9` 未找到记录会返回 `redis.Nil`。<br>
可以在项目中预定义一个未找到记录的 error，将这种错误统一包裹
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFount = errors.NewWithMessage("not found")
)

func sqlFirst() error {
	// first 未找到记录
	err := gorm.ErrRecordNotFound
	if err != nil {
		// 统一用项目预定义错误包裹
		return ErrNotFount.Wrapf(err, "id: %d", 123)
	}
	// ...
	return nil
}

func redisFirst() error {
	// Get 未找到记录
	var err error = redis.Nil
	if err != nil {
		// 统一用项目预定义错误包裹
		return ErrNotFount.Wrapf(err, "key: %d", 456)
	}
	// ...
	return nil
}

func main() {
	err1 := sqlFirst()
	if err1 != nil {
		if errors.Is(err1, ErrNotFount) {
			// 打印堆栈信息
			fmt.Printf("%+v\n", err1)
		}
		// other code
		//return
	}

	err2 := redisFirst()
	if err2 != nil {
		if errors.Is(err2, ErrNotFount) {
			// 打印堆栈信息
			fmt.Printf("%+v\n", err2)
		}
		// other code
		return
	}
}
```

### 4. 给 error 携带信息
子函数（方法）参数 不包含于 父函数（方法）参数时，应将子函数（方法）的参数作为 message 返回，反之则不用
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
)

func A(a int) error {
	return errors.NewWithStack("a err")
}

func B(b, c int) error {
	return errors.NewWithStack("a err")
}

func C() int {
	return 1
}

func service1(a int) error {
	err := A(a)
	if err != nil {
		// A 的参数包含于 service1 的参数
		return err
	}
	// other code
	return nil
}

func service2(b int) error {
	c := C()
	err := B(b, c)
	if err != nil {
		// B 的参数多了 c
		return errors.Wrap(err, "b: %v, c: %v", b, c)
	}
	// other code
	return nil
}

func main() {
	err1 := service1(1)
	fmt.Printf("%+v\n", err1)

	err2 := service2(2)
	fmt.Printf("%+v\n", err2)
}

```


### 5. 吞没 error 时，应增加日志
无需返回 error 时，应记录到日志中，反之不用每个 error 都记录日志
```go
package main

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"log"
)

func A() error {
	return errors.NewWithStack("a err")
}

func service() error {
	err := A()
	if err != nil {
		// 不返回错误时，需要将错误记录到日志中，反之不用每个 error 都记录日志
		log.Printf("%+v", err)
		err = nil
		return err
	}
	// other code
	return nil
}

func main() {
	err1 := service()
	fmt.Printf("%+v\n", err1)
}

```

## 检查 Error
之前判断错误通常使用下面形式
```
func A() error {
	// ...
	err := B()
	if err != nil {
		return err
	}

	// ...
	err = C()
	if err != nil {
		return err
	}
}
```

现在可以这么用
```
func A() (err error) {
	defer errors.Recover(func(e error) {
		err = e
	})

	err := B()
	errors.Check(err)
	// If err is not nil, the code that follows will not execute
	// ...

	err = C()
	errors.Check(err)
	return nil
}
```