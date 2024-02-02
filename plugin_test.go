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

func TestCalculateS3Key(t *testing.T) {
	tests := []struct {
		path         string
		reportSource string
		newFolder    string
		expected     string
	}{
		{
			path:         "/home/user/documents/file.txt",
			reportSource: "/home/user",
			newFolder:    "build-123",
			expected:     "build-123/documents/file.txt",
		},
		{
			path:         "/var/data/image.jpg",
			reportSource: "/var/data",
			newFolder:    "output",
			expected:     "output/image.jpg",
		},
		{
			path:         "logs/error.log",
			reportSource: "logs",
			newFolder:    "build-789",
			expected:     "build-789/error.log",
		},
		{
			path:         "data/results/report.csv",
			reportSource: "data/results",
			newFolder:    "backup",
			expected:     "backup/report.csv",
		},
		{
			path:         "folder/file.txt",
			reportSource: "folder",
			newFolder:    "",
			expected:     "file.txt",
		},
	}

	for _, tc := range tests {
		result := calculateS3Key(tc.path, tc.reportSource, tc.newFolder)
		if result != tc.expected {
			t.Errorf("Expected: %s, Got: %s", tc.expected, result)
		}
	}
}

func TestAppendToFileSlice(t *testing.T) {
	// Test case 1
	files1 := []File{{Name: "file1", URL: "http://example.com/file1"}}
	name1 := "file2"
	url1 := "http://example.com/file2"
	expected1 := []File{{Name: "file1", URL: "http://example.com/file1"}, {Name: "file2", URL: "http://example.com/file2"}}
	result1 := appendToFileSlice(files1, name1, url1)
	checkFileSlice(t, result1, expected1)

	// Test case 2
	files2 := []File{}
	name2 := "document"
	url2 := "https://example.org/document.pdf"
	expected2 := []File{{Name: "document", URL: "https://example.org/document.pdf"}}
	result2 := appendToFileSlice(files2, name2, url2)
	checkFileSlice(t, result2, expected2)

	// Test case 3
	files3 := []File{{Name: "image", URL: "https://example.com/image.jpg"}}
	name3 := "new_image"
	url3 := "https://example.com/new_image.png"
	expected3 := []File{{Name: "image", URL: "https://example.com/image.jpg"}, {Name: "new_image", URL: "https://example.com/new_image.png"}}
	result3 := appendToFileSlice(files3, name3, url3)
	checkFileSlice(t, result3, expected3)

	// Test case 4
	files4 := []File{{Name: "log", URL: "ftp://example.com/log.txt"}}
	name4 := "new_log"
	url4 := "ftp://example.com/new_log.txt"
	expected4 := []File{{Name: "log", URL: "ftp://example.com/log.txt"}, {Name: "new_log", URL: "ftp://example.com/new_log.txt"}}
	result4 := appendToFileSlice(files4, name4, url4)
	checkFileSlice(t, result4, expected4)

	// Test case 5
	files5 := []File{{Name: "file", URL: "file:///home/user/file.txt"}}
	name5 := "new_file"
	url5 := "file:///home/user/new_file.txt"
	expected5 := []File{{Name: "file", URL: "file:///home/user/file.txt"}, {Name: "new_file", URL: "file:///home/user/new_file.txt"}}
	result5 := appendToFileSlice(files5, name5, url5)
	checkFileSlice(t, result5, expected5)
}

// Helper function to check if two slices of File structs are equal
func checkFileSlice(t *testing.T, result, expected []File) {
	t.Helper()
	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, Got length: %d", len(expected), len(result))
		return
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Mismatch at index %d. Expected: %+v, Got: %+v", i, expected[i], result[i])
		}
	}
}
