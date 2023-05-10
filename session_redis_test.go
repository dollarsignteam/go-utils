package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession_Scan(t *testing.T) {
	tests := []struct {
		name        string
		session     Session
		dest        any
		expectedErr error
	}{
		{
			name: "valid data",
			session: Session{
				Meta: SessionMeta{},
				Data: map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			dest: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			expectedErr: nil,
		},
		{
			name: "invalid data",
			session: Session{
				Meta: SessionMeta{},
				Data: "invalid json",
			},
			dest:        new(int),
			expectedErr: errors.New("json: cannot unmarshal string into Go value of type int"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.session.Scan(test.dest)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSessionMeta_param(t *testing.T) {
	tests := []struct {
		name     string
		meta     SessionMeta
		expected sessionKeyParam
	}{
		{
			name: "no group ID",
			meta: SessionMeta{
				ID:     "1234",
				UserID: 5678,
			},
			expected: sessionKeyParam{
				ID:      "1234",
				userID:  "5678",
				groupID: "",
			},
		},
		{
			name: "with group ID",
			meta: SessionMeta{
				ID:      "abcd",
				UserID:  9876,
				GroupID: "xyz",
			},
			expected: sessionKeyParam{
				ID:      "abcd",
				userID:  "9876",
				groupID: "xyz",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.meta.param()
			assert.Equal(t, test.expected, result)
		})
	}
}

func BenchmarkSession_Scan(b *testing.B) {
	session := Session{
		Meta: SessionMeta{},
		Data: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	dest := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	for i := 0; i < b.N; i++ {
		err := session.Scan(dest)
		assert.Nil(b, err)
	}
}
