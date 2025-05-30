package beanstalk_test

import (
	"testing"

	"github.com/artiifact/go-beanstalk"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel_String(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		actual   beanstalk.LogLevel
	}{
		{"debug", "debug", beanstalk.DebugLogLevel},
		{"info", "info", beanstalk.InfoLogLevel},
		{"warning", "warning", beanstalk.WarningLogLevel},
		{"error", "error", beanstalk.ErrorLogLevel},
		{"panic", "panic", beanstalk.PanicLogLevel},
		{"fatal", "fatal", beanstalk.FatalLogLevel},
		{"undefined", "Level(99)", beanstalk.LogLevel(99)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.actual.String())
		})
	}
}
