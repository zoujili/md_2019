package hpid

import (
	"testing"
	"github.com/stretchr/testify/suite"
	testsuite "IdentityService/pkg/test/suite"
)

func TestHPIDProvider(t *testing.T) {
	suite.Run(t, &TestNotificationSuite{})
}


type TestNotificationSuite struct {
	testsuite.UnitTestSuite
	provider                 *Provider
}

func (s *TestNotificationSuite) SetupSuite() {
	s.UnitTestSuite.SetupSuite()
	s.provider = New(s.Config)
}

func (s *TestNotificationSuite) TesAuthOK() {
	s.T().Run("should auth", func(t *testing.T) {
		//given

	})
}
