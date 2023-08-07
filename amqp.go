package utils

import (
	"context"
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

// AMQPClient is a wrapper around the AMQP client
type AMQPClient struct {
	Connection      *amqp.Conn
	SenderSession   *amqp.Session
	ReceiverSession *amqp.Session
	SenderList      []*amqp.Sender
	ReceiverList    []*amqp.Receiver
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
			senderSession.Close(ctx)
			connection.Close()
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
// receiver session and the connection of the AMQPClient.
func (client *AMQPClient) Close() {
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
