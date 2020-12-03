package handler_test

import (
	"net/url"
	"testing"

	"github.com/lucbarr/leaderboard-manager/handler"
	"github.com/stretchr/testify/assert"
)

func TestParseFromURL(t *testing.T) {

	type testData struct {
		query url.Values
		value interface{}

		expectedValue interface{}
		expectedError error
	}

	tt := map[string]testData{
		"success": {
			query: url.Values{
				"potato": []string{"tomato"},
				"season": []string{"10"},
			},
			value: &struct {
				Potato string
				Season int
			}{},

			expectedValue: &struct {
				Potato string
				Season int
			}{Potato: "tomato", Season: 10},

			expectedError: nil,
		},
	}

	for testName, testData := range tt {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			retError := handler.ParseFromQuery(testData.query, testData.value)
			assert.Equal(t, testData.expectedError, retError)
			assert.Equal(t, testData.expectedValue, testData.value)
		})
	}
}
