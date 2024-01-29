package main

import (
	"testing"
)

func TestGetNewFolder(t *testing.T) {
	tests := []struct {
		pipelineSeqID string
		reportTarget  string
		expected      string
	}{
		{
			pipelineSeqID: "123",
			reportTarget:  "",
			expected:      "build-123",
		},
		{
			pipelineSeqID: "456",
			reportTarget:  "reports",
			expected:      "reports/build-456",
		},
		{
			pipelineSeqID: "789",
			reportTarget:  "/output",
			expected:      "/output/build-789",
		},
		{
			pipelineSeqID: "abc",
			reportTarget:  "logs",
			expected:      "logs/build-abc",
		},
		{
			pipelineSeqID: "xyz",
			reportTarget:  "output/folder",
			expected:      "output/folder/build-xyz",
		},
		{
			pipelineSeqID: "111",
			reportTarget:  "",
			expected:      "build-111",
		},
	}

	for _, tc := range tests {
		result := getNewFolder(tc.pipelineSeqID, tc.reportTarget)
		if result != tc.expected {
			t.Errorf("Expected: %s, Got: %s", tc.expected, result)
		}
	}
}
