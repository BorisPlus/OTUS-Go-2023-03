package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	testCases := []struct {
		level          LogLevel
		debugMessage   string
		infoMessage    string
		warningMessage string
		errorMessage   string
		expectedStack  []string
	}{
		{
			DEBUG,
			"DEBUG message",
			"INFO message",
			"WARNING message",
			"ERROR message",
			[]string{"DEBUG message", "INFO message", "WARNING message", "ERROR message"},
		},
		{
			INFO,
			"DEBUG message",
			"INFO message",
			"WARNING message",
			"ERROR message",
			[]string{"INFO message", "WARNING message", "ERROR message"},
		},
		{
			WARNING,
			"DEBUG message",
			"INFO message",
			"WARNING message",
			"ERROR message",
			[]string{"WARNING message", "ERROR message"},
		},
		{
			ERROR,
			"DEBUG message",
			"INFO message",
			"WARNING message",
			"ERROR message",
			[]string{"ERROR message"},
		},
	}
	for _, testCase := range testCases {
		outputInto := &bytes.Buffer{}
		logger := NewLogger(testCase.level, outputInto)
		logger.Debug(testCase.debugMessage)
		logger.Info(testCase.infoMessage)
		logger.Warning(testCase.warningMessage)
		logger.Error(testCase.errorMessage)
		output := outputInto.String()
		for _, expected := range testCase.expectedStack {
			if !strings.Contains(output, expected) {
				t.Errorf("Error logger %s output: %q not contains %s\n", logger.level, output, expected)
			}
		}
	}
}

func TestFormatSuffix(t *testing.T) {
	expected := "I am suffix from formatting trmplate +100500\n"
	outputInto := &bytes.Buffer{}
	logger := NewLogger("INFO", outputInto)
	logger.Info("%s +%d", "I am suffix from formatting trmplate", 100500)
	output := outputInto.String()
	if !strings.HasSuffix(outputInto.String(), expected) {
		t.Errorf("Logger %s output: %q not contains %q\n", logger.level, output, expected)
	}
}
