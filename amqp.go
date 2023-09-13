package utils

import (
	"context"
	"errors"
	"fmt"
	"sync"

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

// AMQPMessageHandlerFunc is the function to handle the received message
type AMQPMessageHandlerFunc func(message *amqp.Message, err error) *AMQPMessageHandler

// AMQPMessageHandler is the handler for the received message
type AMQPMessageHandler struct {
	IsClosed bool
	Rejected bool
}

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

// Send sends the given message to the given sender with the specified persistence flag.
// If the persistence flag is set to true, the message will be saved to disk for durability.
// If the message's header is nil, it will be initialized with a new amqp.MessageHeader.
// The Durable field of the message header will be updated to reflect the persistence flag.
// Note: Setting the persistence flag to true might result in slower performance.
func (client *AMQPClient) Send(sender *amqp.Sender, message *amqp.Message, persistent bool) error {
	if message.Header == nil {
		message.Header = &amqp.MessageHeader{}
	}
	message.Header.Durable = message.Header.Durable || persistent
	return sender.Send(context.TODO(), message, nil)
}

// Publish publishes the given message to the given publisher
func (client *AMQPClient) Publish(publisher *amqp.Sender, message *amqp.Message) error {
	return publisher.Send(context.TODO(), message, nil)
}

// Received receives messages from the given receiver
// and handles them with the provided message handler function.
// If the message handler function returns false, the loop stops.
// If the message handler function returns an error, the message is rejected;
// otherwise, it is accepted for further processing.
func (client *AMQPClient) Received(receiver *amqp.Receiver, messageHandlerFunc AMQPMessageHandlerFunc) {
	for !client.isClosed {
		message, err := receiver.Receive(context.TODO(), nil)
		err, isClosed := client.IsErrorClosed(err)
		if _, closed := client.IsErrorClosed(err); closed {
			return
		}
		h := messageHandlerFunc(message, err)
		if h == nil {
			h = &AMQPMessageHandler{}
		}
		if err == nil {
			if h.Rejected {
				_ = receiver.RejectMessage(context.TODO(), message, nil)
			} else {
				_ = receiver.AcceptMessage(context.TODO(), message)
			}
		}
		if h.IsClosed || isClosed {
			break
		}
	}
	_ = receiver.Close(context.TODO())
}

// ReceivedList receives messages from the given list of receivers
// and handles them with the provided message handler function.
func (client *AMQPClient) ReceivedList(receiverList []*amqp.Receiver, messageHandlerFunc AMQPMessageHandlerFunc) {
	var wg sync.WaitGroup
	for _, receiver := range receiverList {
		wg.Add(1)
		go func(r *amqp.Receiver) {
			defer wg.Done()
			client.Received(r, messageHandlerFunc)
		}(receiver)
	}
	wg.Wait()
}

// IsErrorClosed checks if the given error is a closed error
func (client *AMQPClient) IsErrorClosed(err error) (error, bool) {
	switch e := err.(type) {
	case *amqp.ConnError, *amqp.SessionError, *amqp.LinkError:
		if e.Error() == "EOF" {
			return errors.New("amqp: connection, session, link error"), true
		}
		return e, true
	default:
		return err, false
	}
}
