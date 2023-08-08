package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestAMQPNew(t *testing.T) {
	var tests = []struct {
		name           string
		config         utils.AMQPConfig
		expectedClient *utils.AMQPClient
		expectedError  *string
	}{
		{
			name: "invalid URL",
			config: utils.AMQPConfig{
				URL: "invalid-url",
			},
			expectedClient: nil,
			expectedError:  utils.PointerOf("dial tcp :5672: connect: connection refused"),
		},
		// {
		// 	name: "Test Case 1",
		// 	config: utils.AMQPConfig{
		// 		URL: "amqp://guest:guest@localhost:5672/",
		// 	},
		// 	expectedClient: &utils.AMQPClient{
		// 		// Connection:      connection,
		// 		// SenderSession:   senderSession,
		// 		// ReceiverSession: receiverSession,
		// 	},
		// 	expectedError: nil,
		// },

		// {
		// 	name: "Test Case 3 - Connection Error",
		// 	config: utils.AMQPConfig{
		// 		URL: "amqp://guest:guest@localhost:5672/",
		// 	},
		// 	expectedClient: nil,
		// 	expectedError:  nil,
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := utils.AMQP.New(test.config)
			assert.Equal(t, test.expectedClient, client)
			if test.expectedError == nil {
				assert.EqualError(t, err, *test.expectedError)
			}
		})
	}
}
