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

func (suite *AMQPTestSuite) TearDownSuite() {
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
	wg := sync.WaitGroup{}
	count := 0
	receiverListCount := 9
	r, e := suite.amqpClient.NewReceiver(queueName)
	suite.Nil(e)
	rList, eList := suite.amqpClient.NewReceiverList(queueName, receiverListCount)
	suite.Nil(eList)
	suite.Len(rList, receiverListCount)
	handlerFunc := func(message *amqp.Message, err error) *utils.AMQPMessageHandler {
		suite.Nil(err)
		if err != nil {
			return &utils.AMQPMessageHandler{
				Rejected: true,
				IsClosed: false,
			}
		}
		count++
		defer wg.Done()
		return nil
	}
	go suite.amqpClient.Received(r, handlerFunc)
	go suite.amqpClient.ReceivedList(rList, handlerFunc)
	suite.T().Log("Waiting for all messages to be received...")
	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		go func(index int) {
			id := fmt.Sprintf("index[%d]", index)
			message := amqp.NewMessage([]byte(id))
			message.Properties = &amqp.MessageProperties{
				GroupID: utils.PointerOf(id),
			}
			err := suite.amqpClient.Send(sender, message, false)
			suite.Nil(err)
		}(i)
	}
	suite.T().Log("Waiting for all messages to be sent...")
	wg.Wait()
	suite.Equal(messageCount, count)
}

func (suite *AMQPTestSuite) TestPublishToTopic() {
	topicName := "test-topic"
	messageCount := 1000
	publisher, err := suite.amqpClient.NewPublisher(topicName)
	if err != nil {
		suite.FailNow("failed to create publisher", err)
		return
	}
	suite.NotNil(publisher)
	wg := sync.WaitGroup{}
	count := 0
	r, e := suite.amqpClient.NewSubscriber(topicName)
	suite.Nil(e)
	handlerFunc := func(message *amqp.Message, err error) *utils.AMQPMessageHandler {
		suite.Nil(err)
		if err != nil {
			return &utils.AMQPMessageHandler{
				Rejected: true,
				IsClosed: false,
			}
		}
		count++
		defer wg.Done()
		return &utils.AMQPMessageHandler{
			Rejected: true,
		}
	}
	go suite.amqpClient.Received(r, handlerFunc)
	suite.T().Log("Waiting for all messages to be received...")
	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		go func(index int) {
			id := fmt.Sprintf("index[%d]", index)
			message := amqp.NewMessage([]byte(id))
			err := suite.amqpClient.Publish(publisher, message)
			suite.Nil(err)
		}(i)
	}
	suite.T().Log("Waiting for all messages to be publish...")
	wg.Wait()
	suite.Equal(messageCount, count)
}

func TestIntegrationAMQPTestSuite(t *testing.T) {
	suite.Run(t, new(AMQPTestSuite))
}
