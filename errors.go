package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	sep       = "."
	dirLevels = 1
)

type wrappedError struct {
	err error
}

func (e *wrappedError) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// New создает новую ошибку, добавляя к сообщению информацию о пакете и функции,
// из которой была вызвана ошибка.
//
// Параметры:
//   - msg: Массив строк, объединяется в одно строковое сообщение.
func New(msg ...string) error {
	return newWrap(strings.Join(msg, " "), 3)
}

// Newf создает новую форматированную ошибку, добавляя к сообщению информацию о пакете и функции,
// из которой была вызвана ошибка.
//
// Параметры:
//   - format: Формат строки, как в fmt.Errorf.
//   - args: Аргументы для форматирования строки.
func Newf(format string, args ...any) error {
	return newWrap(format, 3, args...)
}

// Wrap оборачивает существующую ошибку, добавляя к ней информацию о пакете и функции,
// из которой была вызвана ошибка.
//
// Параметры:
//   - err: Существующая ошибка, которую нужно обернуть.
//   - msg: Массив строк, объединяется в одно строковое сообщение.
func Wrap(err error, msg ...string) error {
	return wrap(err, strings.Join(msg, " "), 3)
}

// Wrapf оборачивает существующую ошибку, добавляя к ней информацию о пакете и функции,
// из которой была вызвана ошибка.
//
// Параметры:
//   - err: Существующая ошибка, которую нужно обернуть.
//   - format: Формат строки, как в fmt.Errorf.
//   - args: Аргументы для форматирования строки.
func Wrapf(err error, format string, args ...any) error {
	return wrap(err, format, 3, args...)
}

// newWrap является вспомогательной функцией для создания ошибки с информацией о пакете и функции.
//
// Параметры:
//   - format: Сообщение об ошибке.
//   - skip: Количество фреймов стека вызовов для пропуска при определении вызывающей функции.
//   - args: Аргументы форматированной строки
func newWrap(format string, skip int, args ...any) error {
	packageName, functionName := getFuncInfo(skip)

	err := fmt.Errorf(format, args...)

	wErr := &wrappedError{
		err: fmt.Errorf("[%s%s%s] %w", packageName, sep, functionName, err),
	}

	return wErr
}

// wrap является вспомогательной функцией для оборачивания ошибки с информацией о пакете и функции.
//
// Параметры:
//   - err: Существующая ошибка, которую нужно обернуть.
//   - msg: Сообщение для добавления к ошибке.
//   - skip: Количество фреймов стека вызовов для пропуска при определении вызывающей функции.
func wrap(err error, format string, skip int, args ...any) error {
	if err == nil {
		return nil
	}

	packageName, functionName := getFuncInfo(skip)

	msgErr := fmt.Errorf(format, args...)

	if wErr := new(wrappedError); !errors.As(err, &wErr) {
		err = &wrappedError{err: err}
	}

	return fmt.Errorf("[%s%s%s] %w -> %w", packageName, sep, functionName, msgErr, err)
}

// getFuncInfo получает информацию о пакете и функции, из которой была вызвана ошибка.
//
// Параметры:
//   - skip: Количество фреймов стека вызовов для пропуска при определении вызывающей функции.
//
// Возвращает:
//   - Имя пакета и функции.
func getFuncInfo(skip int) (packageName, functionName string) {
	packageName = "unknown"
	functionName = "unknown"

	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	funcDetails := runtime.FuncForPC(pc)
	if funcDetails == nil {
		return
	}

	full := funcDetails.Name() // пример: "github.com/user/project/pkg.(*Type).Method"
	parts := strings.Split(full, "/")
	if len(parts) == 0 {
		return
	}

	last := parts[len(parts)-1] // пример: "pkg.(*Type).Method"
	nameParts := strings.SplitN(last, ".", 2)

	if len(nameParts) > 0 {
		packageName = nameParts[0]
	}
	if len(nameParts) > 1 {
		functionName = nameParts[1]
	}

	return
}

func (e *wrappedError) Unwrap() error {
	return e.err
}

func (e *wrappedError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *wrappedError) Cause() error {
	cause := e.err
	for {
		unwrapped := errors.Unwrap(cause)
		if unwrapped == nil {
			break
		}
		cause = unwrapped
	}
	return cause
}

// Join объединяет несколько ошибок с контекстом вызова.
func Join(errs ...error) error {
	var nonNil []error
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}

	if len(nonNil) == 0 {
		return nil
	}

	return newWrap("%w", 3, errors.Join(nonNil...))
}
