package lib

import (
	"github.com/Yakwilik/errors"
	"github.com/Yakwilik/errors/teststub"
)

type Interface interface {
	DoWork(needError bool) error
}

type stubStruct struct {
	stubDependency teststub.StubInterface
}

func (s *stubStruct) DoWork(needError bool) error {
	_, err := s.stubDependency.GetData(needError)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func New() Interface {
	return &stubStruct{
		stubDependency: teststub.New(false),
	}
}
