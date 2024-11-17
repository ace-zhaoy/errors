package errors

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	e := &withJoin{errs: make([]error, 0, len(errs))}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	if len(e.errs) == 0 {
		return nil
	}
	return e
}

func NewWithJoin(errs ...error) ErrorJoin {
	e := &withJoin{errs: make([]error, 0, len(errs))}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type ErrorJoin interface {
	error
	Append(error)
	Unwrap() []error
	Len() int
	ToError() error
}

type withJoin struct {
	errs []error
	mu   sync.Mutex
}

func (w *withJoin) Error() string {
	if len(w.errs) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(w.errs[0].Error())
	for _, err := range w.errs[1:] {
		builder.WriteByte('\n')
		builder.WriteString(err.Error())
	}
	return builder.String()
}

func (w *withJoin) Append(err error) {
	if err == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.errs = append(w.errs, err)
}

func (w *withJoin) Is(err error) bool {
	for _, e := range w.errs {
		if Is(e, err) {
			return true
		}
	}
	return false
}

func (w *withJoin) Unwrap() []error {
	return w.errs
}

func (w *withJoin) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, err := range w.errs {
				fmt.Fprintf(s, "%+v\n", err)
			}
			return
		}
		fallthrough
	case 's', 'q':
		for _, err := range w.errs {
			io.WriteString(s, err.Error())
		}
	}
}

func (w *withJoin) Len() int {
	return len(w.errs)
}

func (w *withJoin) ToError() error {
	if len(w.errs) == 0 {
		return nil
	}
	return w
}
