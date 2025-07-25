package teststub

import (
	"github.com/Yakwilik/errors"
)

type StubInterface interface {
	GetData(fail bool) (any, error)
}

type mockDependency struct {
	returnError bool
}

func (m *mockDependency) GetData(fail bool) (any, error) {
	if m.returnError || fail {
		return nil, errors.New("mock error")
	}
	return "data", nil
}

func New(returnError bool) StubInterface {
	return &mockDependency{returnError: returnError}
}
