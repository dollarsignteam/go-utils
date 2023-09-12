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
	queueName := "test-queue"
	messageCount := 1000
	sender, err := suite.amqpClient.NewSender(queueName)
	if err != nil {
		suite.FailNow("failed to create sender", err)
		return
	}
	suite.NotNil(sender)
	wg1 := sync.WaitGroup{}
	for i := 0; i < messageCount; i++ {
		wg1.Add(1)
		go func(index int) {
			defer wg1.Done()
			id := fmt.Sprintf("index[%d]", index)
			message := amqp.NewMessage([]byte(id))
			message.Properties = &amqp.MessageProperties{
				CorrelationID: id,
				GroupID:       utils.PointerOf(id),
			}
			err := suite.amqpClient.Send(sender, message, false)
			if err != nil {
				suite.Nil(err)
			}
		}(i)
	}
	suite.T().Log("Waiting for all messages to be sent...")
	wg1.Wait()
	count := 0
	receiverListCount := 10
	// r, e := suite.amqpClient.NewReceiver(queueName)
	// suite.Nil(e)
	rList, eList := suite.amqpClient.NewReceiverList(queueName, receiverListCount)
	suite.Nil(eList)
	suite.Len(rList, receiverListCount)
	wg2 := sync.WaitGroup{}
	wg2.Add(messageCount)
	handlerFunc := func(message *amqp.Message, err error) utils.AMQPMessageHandlerResult {
		count++
		return utils.AMQPMessageHandlerResult{
			Callback: func() {
				suite.T().Log("process")
				defer wg2.Done()
			},
		}
	}
	// go suite.amqpClient.Received(r, handlerFunc)
	go suite.amqpClient.ReceivedList(rList, handlerFunc)
	suite.T().Log("Waiting for all messages to be received...")
	wg2.Wait()
	suite.Equal(messageCount, count)
}

func (suite *AMQPTestSuite) TestPublishToTopic() {
	publisher, err := suite.amqpClient.NewPublisher("test-topic")
	if err != nil {
		suite.FailNow("failed to create publisher", err)
		return
	}
	suite.NotNil(publisher)
	wg1 := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg1.Add(1)
		go func(index int) {
			defer wg1.Done()
			id := fmt.Sprintf("index[%d]", index)
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
	wg1.Wait()
}

func TestIntegrationAMQPTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPTestSuite))
}
