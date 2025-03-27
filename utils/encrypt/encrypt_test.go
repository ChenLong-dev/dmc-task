package encrypt

import "testing"

func TestMD5Hash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "with space",
			input:    "hello world",
			expected: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
		{
			name:     "case sensitivity",
			input:    "Hello",
			expected: "8b1a9953c4611296a827abf8c47804d7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5Hash(tt.input); got != tt.expected {
				t.Errorf("MD5Hash() = %v, want %v", got, tt.expected)
			}
		})
	}
}
