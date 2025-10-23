package converter

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConserterSuite struct {
	suite.Suite
}

func (s *ConserterSuite) SetupTest() {
}

func (s *ConserterSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ConserterSuite))
}
