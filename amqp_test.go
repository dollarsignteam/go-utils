//go:build integration

package utils_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/suite"

	"github.com/dollarsignteam/go-utils"
)

type AMQPTestSuite struct {
	suite.Suite
	amqpClient *utils.AMQPClient
	wg         sync.WaitGroup
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

func (suite *AMQPTestSuite) TestSendToQueue() {
	sender, err := suite.amqpClient.NewSender("test-queue")
	if err != nil {
		suite.FailNow("failed to create sender", err)
		return
	}
	suite.NotNil(sender)
	for i := 0; i < 1000; i++ {
		suite.wg.Add(1)
		go func(index int) {
			defer suite.wg.Done()
			id := fmt.Sprintf("id = %d", index+1)
			message := amqp.NewMessage([]byte(id))
			message.Properties = &amqp.MessageProperties{
				CorrelationID: id,
			}
			err := suite.amqpClient.Send(sender, message, false)
			if err != nil {
				suite.Nil(err)
			}
		}(i)
	}
	suite.T().Log("Waiting for all messages to be sent...")
	suite.wg.Wait()
}

func (suite *AMQPTestSuite) TestPublishToTopic() {
	publisher, err := suite.amqpClient.NewPublisher("test-topic")
	if err != nil {
		suite.FailNow("failed to create publisher", err)
		return
	}
	suite.NotNil(publisher)
	for i := 0; i < 1000; i++ {
		suite.wg.Add(1)
		go func(index int) {
			defer suite.wg.Done()
			id := fmt.Sprintf("id = %d", index+1)
			message := amqp.NewMessage([]byte(id))
			message.Properties = &amqp.MessageProperties{
				CorrelationID: id,
			}
			err := suite.amqpClient.Publish(publisher, message)
			if err != nil {
				suite.Nil(err)
			}
		}(i)
	}
	suite.T().Log("Waiting for all messages to be sent...")
	suite.wg.Wait()
}

func TestIntegrationAMQPTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPTestSuite))
}
