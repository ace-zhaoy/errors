package errors

import (
	"errors"
	"fmt"
	"io"
)

type ErrorMessage interface {
	Message() string
	WrapStack(err error) error
	Wrapf(err error, format string, args ...any) error
	error
}

type withMessage struct {
	message string
	cause   error
}

func (w *withMessage) Message() string {
	return w.message
}

func (w *withMessage) Error() string {
	if w.cause == nil {
		return w.message
	}
	return fmt.Sprintf("%s -> {%s}", w.Message(), w.Cause().Error())
}

func (w *withMessage) Is(err error) bool {
	if e, ok := err.(*withMessage); ok {
		return w.Message() == e.Message()
	}
	return false
}

func (w *withMessage) Cause() error {
	return w.cause
}

func (w *withMessage) Unwrap() error {
	return w.cause
}

func (w *withMessage) WrapStack(err error) error {
	ws := &withStack{
		error: &withMessage{
			message: w.Message(),
			cause:   err,
		},
	}
	if !stackExists(ws) {
		ws.stack = callers()
	}
	return ws
}

func (w *withMessage) Wrapf(err error, format string, args ...any) error {
	wm := &withMessage{
		message: w.Message(),
		cause:   err,
	}
	ws := &withStack{
		error: &withMessage{
			message: fmt.Sprintf(format, args...),
			cause:   wm,
		},
	}
	if !stackExists(ws) {
		ws.stack = callers()
	}
	return ws
}

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if w.Cause() != nil {
				fmt.Fprintf(s, "%+v\n", w.Cause())
			}
			io.WriteString(s, w.Message())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type ErrorCode interface {
	Code() int
	ErrorMessage
}

type withCode struct {
	*withMessage
	code int
}

func (w *withCode) Code() int {
	return w.code
}

func (w *withCode) WrapStack(err error) error {
	ws := &withStack{
		error: &withCode{
			withMessage: &withMessage{
				message: w.Message(),
				cause:   err,
			},
			code: w.Code(),
		},
	}
	if !stackExists(ws) {
		ws.stack = callers()
	}
	return ws
}

func (w *withCode) Wrapf(err error, format string, args ...any) error {
	wm := &withMessage{
		message: w.Message(),
		cause:   err,
	}
	wc := &withCode{
		withMessage: wm,
		code:        w.Code(),
	}
	ws := &withStack{
		error: &withMessage{
			message: fmt.Sprintf(format, args...),
			cause:   wc,
		},
	}
	if !stackExists(ws) {
		ws.stack = callers()
	}
	return ws
}

func (w *withCode) Error() string {
	s := fmt.Sprintf("[%d: %s]", w.Code(), w.Message())
	if w.Cause() == nil {
		return s
	}
	return fmt.Sprintf("%s -> {%s}", s, w.Cause())
}

func (w *withCode) Is(err error) bool {
	if e, ok := err.(*withCode); ok {
		return w.Code() == e.Code()
	}
	return false
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if w.Cause() != nil {
				fmt.Fprintf(s, "%+v\n", w.Cause())
			}
			fmt.Fprintf(s, "%d: %s", w.Code(), w.Message())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error {
	return w.error
}

func (w *withStack) Unwrap() error {
	return w.error
}

func (w *withStack) Wrap(err error) error {
	if err == nil {
		return nil
	}

	e := &withStack{
		error: err,
	}
	if !stackExists(err) {
		e.stack = callers()
	}

	return e
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.error)
			if w.stack != nil {
				w.stack.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type causer interface {
	Cause() error
}

func stackExists(err error) bool {
	for err != nil {
		if e, ok := err.(*withStack); ok {
			if e.stack != nil {
				return true
			}
			err = e.Cause()
			continue
		}

		e, ok := err.(causer)
		if !ok {
			break
		}
		err = e.Cause()
	}
	return false
}

// New returns an error with no code but a message
// New no stack
func New(message string) error {
	return &withMessage{
		message: message,
	}
}

func NewWithMessage(format string, args ...any) ErrorMessage {
	return &withMessage{
		message: fmt.Sprintf(format, args...),
	}
}

// NewWithCode returns an error with code and message
// NewWithCode no stack
func NewWithCode(code int, format string, args ...any) ErrorCode {
	return &withCode{
		code: code,
		withMessage: &withMessage{
			message: fmt.Sprintf(format, args...),
		},
	}
}

// NewWithStack returns an error without code but with message and stack
func NewWithStack(format string, args ...any) error {
	err := &withMessage{
		message: fmt.Sprintf(format, args...),
	}
	return &withStack{
		error: err,
		stack: callers(),
	}
}

// WithStack annotates err with a stack trace only once
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if stackExists(err) {
		return err
	}
	if e, ok := err.(*withStack); ok {
		e.stack = callers()
		return e
	}
	return &withStack{
		error: err,
		stack: callers(),
	}
}

func WithStackForce(err error) error {
	if err == nil {
		return nil
	}

	return &withStack{
		error: err,
		stack: callers(),
	}
}

func Wrap(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}
	if stackExists(err) {
		return err
	}
	if e, ok := err.(*withStack); ok {
		e.stack = callers()
		return e
	}
	return &withStack{
		error: err,
		stack: callers(),
	}
}

func WrapForce(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		message: fmt.Sprintf(format, args...),
		cause:   err,
	}

	return &withStack{
		error: err,
		stack: callers(),
	}
}

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok || cause.Cause() == nil {
			break
		}

		err = cause.Cause()
	}
	return err
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}

// LatestCode returns the latest ErrorCode
func LatestCode(err error) ErrorCode {
	for err != nil {
		ex, ok := err.(ErrorCode)
		if ok {
			return ex
		}

		e, ok := err.(causer)
		if !ok {
			break
		}
		err = e.Cause()
	}
	return nil
}

// LatestMessage returns the latest ErrorMessage
func LatestMessage(err error) ErrorMessage {
	for err != nil {
		ex, ok := err.(ErrorMessage)
		if ok {
			return ex
		}

		e, ok := err.(causer)
		if !ok {
			break
		}
		err = e.Cause()
	}
	return nil
}
