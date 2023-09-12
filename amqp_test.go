//go:build integration

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dollarsignteam/go-utils"
)

type AMQPTestSuite struct {
	suite.Suite
	amqpClient *utils.AMQPClient
}

func (suite *AMQPTestSuite) SetupTest() {
	client, err := utils.AMQP.New(utils.AMQPConfig{
		URL: "amqp://admin:admin@localhost:5672"})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.amqpClient = client
}

func (suite *AMQPTestSuite) TearDownTest() {
	if suite.amqpClient != nil {
		suite.amqpClient.Close()
	}
}

func (suite *AMQPTestSuite) TestAddOne() {
	suite.Equal(1, 1)
}

func TestIntegrationAMQPTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPTestSuite))
}
