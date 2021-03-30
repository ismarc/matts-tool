package app

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func (s *Suite) SetupSuite() {
}

// AfterTest comment
func (s *Suite) AfterTest(_, _ string) {
}

func TestIncludeProcessor(t *testing.T) {
	suite.Run(t, new(Suite))
}
