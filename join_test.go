package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestJoin_nil(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"err empty",
			args{[]error{}},
		},
		{
			"err nil",
			args{[]error{nil}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Join(tt.args.errs...)
			if got != nil {
				t.Errorf("Join() error = %v, wantErr nil", got)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr []error
	}{
		{
			"err1, err2",
			args{[]error{errors.New("err1"), errors.New("err2")}},
			[]error{errors.New("err1"), errors.New("err2")},
		},
		{
			"err1, err2, nil",
			args{[]error{errors.New("err1"), nil, errors.New("err2")}},
			[]error{errors.New("err1"), errors.New("err2")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Join(tt.args.errs...).(*withJoin).Unwrap()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_withJoin_Error(t *testing.T) {
	type fields struct {
		errs []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"errors.New: err1",
			fields{[]error{errors.New("err1")}},
			"err1",
		},
		{
			"errors.New: err1, err2",
			fields{[]error{errors.New("err1"), errors.New("err2")}},
			"err1\nerr2",
		},
		{
			"New: err1",
			fields{[]error{New("err1")}},
			"err1",
		},
		{
			"New: err1, err2",
			fields{[]error{New("err1"), New("err2")}},
			"err1\nerr2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &withJoin{
				errs: tt.fields.errs,
			}
			if got := w.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withJoin_Append(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
	jerr := NewWithJoin()
	errs := []error{
		New("err0"),
		New("err1"),
		New("err2"),
		New("err3"),
		New("err4"),
		New("err5"),
		New("err6"),
		New("err7"),
		New("err8"),
		New("err9"),
	}
	g := sync.WaitGroup{}
	for _, err := range errs {
		g.Add(1)
		go func(err error) {
			defer g.Done()
			jerr.Append(err)
		}(err)
	}
	g.Wait()
	if len(jerr.Unwrap()) != len(errs) {
		t.Errorf("Append() = %v, want %v", len(jerr.Unwrap()), len(errs))
	}
}

func Test_withJoin_Is(t *testing.T) {
	type fields struct {
		errs []error
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"err1",
			fields{[]error{New("err1")}},
			args{New("err1")},
			true,
		},
		{
			"err2",
			fields{[]error{New("err1"), New("err2")}},
			args{New("err2")},
			true,
		},
		{
			"err1-2",
			fields{[]error{New("err1"), New("err2")}},
			args{New("err1")},
			true,
		},
		{
			"err3",
			fields{[]error{New("err1"), New("err2")}},
			args{New("err3")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &withJoin{
				errs: tt.fields.errs,
			}
			if got := w.Is(tt.args.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withJoin_Format(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := New("error 2")
	err3 := NewWithStack("error 3")
	e := &withJoin{errs: []error{err1, err2, err3}}

	tests := []struct {
		format   string
		expected string
	}{
		{"%s", "error 1error 2error 3"},
		{"%q", "error 1error 2error 3"},
		{"%v", "error 1error 2error 3"},
		{"%+v", "error 1\nerror 2\nerror 3\ngithub.com/ace-zhaoy/errors.Test_withJoin_Format"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, e)
		got = strings.SplitN(got, "\n\t", 2)[0]
		if got != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, got)
		}
	}
}
