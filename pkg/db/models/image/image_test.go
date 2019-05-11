package image

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImageTestSuite struct {
	suite.Suite
}

func TestImageTestSuite(t *testing.T) {
	tests := new(ImageTestSuite)
	suite.Run(t, tests)
}
