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
