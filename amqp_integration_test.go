//go:build integration

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestIntegrationAMQP(t *testing.T) {
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
