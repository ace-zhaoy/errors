package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	CodeInvalidParams  = 400102030
	CodeUserNotFound   = 500201010
	CodeOrderNotExists = 500302001
)

var (
	ErrCodeInvalidParams  = NewWithCode(CodeInvalidParams, "invalid params")
	ErrCodeUserNotFound   = NewWithCode(CodeUserNotFound, "user not found")
	ErrCodeOrderNotExists = NewWithCode(CodeOrderNotExists, "order not exists")
)

var (
	ErrInvalidParams  = NewWithMessage("invalid params")
	ErrUserNotFound   = NewWithMessage("user not found")
	ErrOrderNotExists = NewWithMessage("order not exists")
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "new abc",
			input:   "abc",
			wantErr: "abc",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.input)
			if err == nil {
				t.Errorf("New() returned nil, want error")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("New() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestNewWithCode(t *testing.T) {
	type args struct {
		code   int
		format string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want ErrorCode
	}{
		{
			"ErrCodeUserNotFound",
			args{
				code:   CodeUserNotFound,
				format: "user not found",
				args:   nil,
			},
			ErrCodeUserNotFound,
		},
		{
			"error format",
			args{
				code:   100,
				format: "aaa: %s",
				args:   []any{"bbb"},
			},
			&withCode{
				code:    100,
				message: "aaa: bbb",
			},
		},
		{
			"empty string",
			args{
				code:   200,
				format: "",
			},
			&withCode{
				code:    200,
				message: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithCode(tt.args.code, tt.args.format, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithMessage(t *testing.T) {
	type args struct {
		format string
		args   []any
	}
	tests := []struct {
		name string
		args args
		want ErrorMessage
	}{
		{
			"ErrUserNotFound",
			args{
				format: "user not found",
			},
			ErrUserNotFound,
		},
		{
			"empty string",
			args{
				format: "",
			},
			&withMessage{
				message: "",
			},
		},
		{
			"abc",
			args{
				format: "abc",
			},
			&withMessage{
				message: "abc",
			},
		},
		{
			"format abc: 123",
			args{
				format: "abc: %v",
				args:   []any{123},
			},
			&withMessage{
				message: "abc: 123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithMessage(tt.args.format, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithStack(t *testing.T) {
	type args struct {
		format string
		args   []any
	}
	tests := []struct {
		name    string
		args    args
		wantStr []string
	}{
		{
			"aaa",
			args{
				format: "aaa",
			},
			[]string{"aaa", "TestNewWithStack.func1", "error_test.go", "testing.tRunner", "testing.go:", "runtime.goexit", "runtime"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewWithStack(tt.args.format, tt.args.args...)
			str := fmt.Sprintf("%+v", err)
			strArr := strings.Split(str, "\n")
			if len(strArr) != len(tt.wantStr) {
				t.Errorf("NewWithStack() = %v, want %v", strArr, tt.wantStr)
				return
			}
			for i, v := range strArr {
				if !strings.Contains(v, tt.wantStr[i]) {
					t.Errorf("NewWithStack() %v = %v, want %v", i, v, tt.wantStr[i])
				}
			}
		})
	}
}

func TestCause(t *testing.T) {
	err1 := Wrap(ErrCodeUserNotFound, "err1")
	err2 := Wrap(err1, "err2")
	err3 := Wrap(ErrOrderNotExists, "err3")
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			"nil",
			args{nil},
			nil,
		},
		{
			"err1 - ErrCodeUserNotFound",
			args{err1},
			ErrCodeUserNotFound,
		},
		{
			"err2 - ErrCodeUserNotFound",
			args{err2},
			ErrCodeUserNotFound,
		},
		{
			"err3 - ErrOrderNotExists",
			args{err3},
			ErrOrderNotExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Cause(tt.args.err); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Cause() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIs(t *testing.T) {
	err1 := Wrap(ErrCodeUserNotFound, "err1")
	err2 := ErrCodeInvalidParams.Wrapf(err1, "err2")
	err3 := Wrap(err2, "err3")

	err11 := Wrap(ErrUserNotFound, "err11")
	err22 := ErrOrderNotExists.Wrapf(err11, "err22")
	err33 := Wrap(err22, "err33")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"err1 is ErrCodeUserNotFound",
			args{err1, ErrCodeUserNotFound},
			true,
		},
		{
			"err3 is ErrCodeUserNotFound",
			args{err3, ErrCodeUserNotFound},
			true,
		},
		{
			"err3 is not ErrCodeOrderNotExists",
			args{err3, ErrCodeOrderNotExists},
			false,
		},
		{
			"err11 is ErrUserNotFound",
			args{err11, ErrUserNotFound},
			true,
		},
		{
			"err33 is ErrOrderNotExists",
			args{err33, ErrOrderNotExists},
			true,
		},
		{
			"err33 is not ErrInvalidParams",
			args{err33, ErrInvalidParams},
			false,
		},
		{
			"ErrUserNotFound is not ErrCodeUserNotFound",
			args{ErrUserNotFound, ErrCodeUserNotFound},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatestCode(t *testing.T) {
	err1 := Wrap(ErrCodeUserNotFound, "err1")
	err2 := ErrCodeInvalidParams.Wrapf(err1, "err2")
	err3 := Wrap(err2, "err3")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want ErrorCode
	}{
		{
			"ErrCodeInvalidParams",
			args{err3},
			ErrCodeInvalidParams,
		},
		{
			"ErrCodeUserNotFound",
			args{err1},
			ErrCodeUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LatestCode(tt.args.err); got.Code() != tt.want.Code() {
				t.Errorf("LatestCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatestMessage(t *testing.T) {
	err11 := Wrap(ErrUserNotFound, "err11")
	err22 := ErrOrderNotExists.Wrapf(err11, "err22")
	err33 := Wrap(err22, "err33")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"ErrOrderNotExists",
			args{err33},
			"err33",
		},
		{
			"ErrUserNotFound",
			args{err11},
			"err11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LatestMessage(tt.args.err); got.Message() != tt.want {
				t.Errorf("LatestMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithStack(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			"go err",
			args{errors.New("go err")},
			errors.New("go err"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WithStack(tt.args.err); err.Error() != tt.wantErr.Error() {
				t.Errorf("WithStack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err    error
		format string
		args   []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			"ErrCodeUserNotFound - err1",
			args{
				err:    ErrCodeUserNotFound,
				format: "err1",
			},
			"err1 -> {[500201010: user not found]}",
		},
		{
			"go err - err2",
			args{
				err:    errors.New("go err"),
				format: "err2",
			},
			"err2 -> {go err}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.err, tt.args.format, tt.args.args...); err.Error() != tt.wantErr {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckWithStack(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr string
	}{
		{
			"nil error",
			nil,
			"",
		},
		{
			"regular error",
			errors.New("regular error"),
			"regular error",
		},
		{
			"error with stack",
			WithStack(errors.New("with stack")),
			"with stack",
		},
		{
			"error code",
			ErrCodeUserNotFound,
			"[500201010: user not found]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				CheckWithStack(tt.err)
				return
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						t.Errorf("Expected panic with error, got %v", r)
						return
					}

					if !strings.Contains(err.Error(), tt.wantErr) {
						t.Errorf("CheckWithStack() panicked with = %v, want containing %v", err, tt.wantErr)
					}

					_, hasStack := r.(*withStack)
					if !hasStack && !stackExists(err) {
						t.Errorf("CheckWithStack() error should have stack information")
					}
				}
			}()

			CheckWithStack(tt.err)
			t.Errorf("CheckWithStack() did not panic as expected")
		})
	}
}

func TestCheckWithWrap(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		format  string
		args    []any
		wantErr string
	}{
		{
			"nil error",
			nil,
			"format",
			nil,
			"",
		},
		{
			"regular error",
			errors.New("regular error"),
			"wrapped error: %s",
			[]any{"additional info"},
			"wrapped error: additional info -> {regular error}",
		},
		{
			"error with stack",
			WithStack(errors.New("with stack")),
			"wrapped error: %s",
			[]any{"additional info"},
			"wrapped error: additional info -> {with stack}",
		},
		{
			"error code",
			ErrCodeUserNotFound,
			"wrapped error: %s",
			[]any{"additional info"},
			"wrapped error: additional info -> {[500201010: user not found]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				CheckWithWrap(tt.err, tt.format, tt.args...)
				return
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						t.Errorf("Expected panic with error, got %v", r)
						return
					}

					if !strings.Contains(err.Error(), tt.wantErr) {
						t.Errorf("CheckWithWrap() panicked with = %v, want containing %v", err, tt.wantErr)
					}

					_, hasStack := r.(*withStack)
					if !hasStack && !stackExists(err) {
						t.Errorf("CheckWithWrap() error should have stack information")
					}
				}
			}()

			CheckWithWrap(tt.err, tt.format, tt.args...)
			t.Errorf("CheckWithWrap() did not panic as expected")
		})
	}
}
