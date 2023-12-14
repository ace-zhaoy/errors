package errors

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
)

type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				pc = pc - 1
				fn := runtime.FuncForPC(pc)
				name, line, file := "unknown", 0, "unknown"
				if fn != nil {
					name = fn.Name()
					file, line = fn.FileLine(pc)
				}
				io.WriteString(st, "\n")
				io.WriteString(st, name)
				io.WriteString(st, "\n\t")
				io.WriteString(st, file)
				io.WriteString(st, ":")
				io.WriteString(st, strconv.Itoa(line))
			}
		}
	}
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
