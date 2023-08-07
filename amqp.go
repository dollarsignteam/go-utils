package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
)

// AMQP utility instance
var AMQP AMQPUtil

// AMQPUtil is a utility for AMQP
type AMQPUtil struct{}

// AMQPConfig is the configuration for the AMQP client
type AMQPConfig struct {
	URL string
}

// MessageHandlerFunc is the function to handle the received message
type MessageHandlerFunc func(message *amqp.Message, err error) (bool, error)

// AMQPClient is a wrapper around the AMQP client
type AMQPClient struct {
	Connection      *amqp.Conn
	SenderSession   *amqp.Session
	ReceiverSession *amqp.Session
	SenderList      []*amqp.Sender
	ReceiverList    []*amqp.Receiver
	mutex           sync.Mutex
	isClosed        bool
}

var amqpTimeoutTTL = 60 * time.Second

// New creates a new AMQP client
func (AMQPUtil) New(config AMQPConfig) (*AMQPClient, error) {
	ctx := context.TODO()
	connection, err := amqp.Dial(ctx, config.URL, nil)
	if err != nil {
		return nil, err
	}
	senderSession, err := connection.NewSession(ctx, nil)
	if err != nil {
		defer connection.Close()
		return nil, err
	}
	receiverSession, err := connection.NewSession(ctx, nil)
	if err != nil {
		defer func() {
			_ = senderSession.Close(ctx)
			_ = connection.Close()
		}()
		return nil, err
	}
	return &AMQPClient{
		Connection:      connection,
		SenderSession:   senderSession,
		ReceiverSession: receiverSession,
	}, nil
}

// Close closes all senders, receivers, sender session,
// receiver session and the connection of the AMQPClient
func (client *AMQPClient) Close() {
	client.isClosed = true
	var wg sync.WaitGroup
	for _, sender := range client.SenderList {
		wg.Add(1)
		go func(s *amqp.Sender) {
			defer wg.Done()
			_ = s.Close(context.TODO())
		}(sender)
	}
	for _, receiver := range client.ReceiverList {
		wg.Add(1)
		go func(r *amqp.Receiver) {
			defer wg.Done()
			_ = r.Close(context.TODO())
		}(receiver)
	}
	wg.Wait()
	_ = client.SenderSession.Close(context.TODO())
	_ = client.ReceiverSession.Close(context.TODO())
	_ = client.Connection.Close()
}

// NewSender creates a new sender for the given queue
func (client *AMQPClient) NewSender(queue string) (*amqp.Sender, error) {
	sender, err := client.SenderSession.NewSender(context.TODO(), queue, nil)
	if err != nil {
		return nil, err
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.SenderList = append(client.SenderList, sender)
	return sender, nil
}

// NewReceiver creates a new receiver for the given queue
func (client *AMQPClient) NewReceiver(queue string) (*amqp.Receiver, error) {
	receiver, err := client.ReceiverSession.NewReceiver(context.TODO(), queue, nil)
	if err != nil {
		return nil, err
	}
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.ReceiverList = append(client.ReceiverList, receiver)
	return receiver, nil
}

// NewSenderList creates a list of senders for the given queue
func (client *AMQPClient) NewReceiverList(queue string, count int) ([]*amqp.Receiver, error) {
	var receiverList []*amqp.Receiver
	for i := 0; i < count; i++ {
		receiver, err := client.NewReceiver(queue)
		if err != nil {
			return nil, err
		}
		receiverList = append(receiverList, receiver)
	}
	return receiverList, nil
}

// NewPublisher creates a new publisher for the given topic
func (client *AMQPClient) NewPublisher(topic string) (*amqp.Sender, error) {
	return client.NewSender(fmt.Sprintf("topic://%s", topic))
}

// NewSubscriber creates a new subscriber for the given topic
func (client *AMQPClient) NewSubscriber(topic string) (*amqp.Receiver, error) {
	return client.NewReceiver(fmt.Sprintf("topic://%s", topic))
}

// Send sends the given message to the given sender
func (client *AMQPClient) Send(sender *amqp.Sender, message *amqp.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), amqpTimeoutTTL)
	defer cancel()
	if message.Header == nil {
		message.Header = &amqp.MessageHeader{}
	}
	message.Header.Durable = true
	return sender.Send(ctx, message, nil)
}

// Publish publishes the given message to the given publisher
func (client *AMQPClient) Publish(publisher *amqp.Sender, message *amqp.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), amqpTimeoutTTL)
	defer cancel()
	return publisher.Send(ctx, message, nil)
}

// Received receives messages from the given receiver
// and handles them with the provided message handler function.
// If the message handler function returns false, the loop stops.
// If the message handler function returns an error, the message is rejected;
// otherwise, it is accepted for further processing.
func (client *AMQPClient) Received(receiver *amqp.Receiver, messageHandlerFunc MessageHandlerFunc) {
	for !client.isClosed {
		message, err := receiver.Receive(context.TODO(), nil)
		ok, err := messageHandlerFunc(message, err)
		if err != nil {
			_ = receiver.RejectMessage(context.TODO(), message, nil)
		} else {
			_ = receiver.AcceptMessage(context.TODO(), message)
		}
		if !ok {
			_ = receiver.Close(context.TODO())
			break
		}
	}
}
