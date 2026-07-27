package driver

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/pb33f/testify/require"
	"github.com/pmezard/go-difflib/difflib"
)

type StreamType uint8

const (
	StreamFile   StreamType = 1
	StreamBuffer StreamType = 2
)

// Output Type constant
// const OutputType = StreamBuffer
const OutputType = StreamFile

// Stream defines a unified interface for reading, writing, seeking,
// closing, and retrieving raw byte contents.
type Stream interface {
	io.ReadWriteSeeker
	io.Closer
	Bytes() ([]byte, error)
}

// ============================================================================
// Memory / String Buffer Wrapper
// ============================================================================

// MemBuffer wraps an in-memory byte slice to satisfy the Stream interface.
type MemBuffer struct {
	buf []byte
	off int64
}

// NewMemBuffer creates a new memory-backed Stream initialized with string data.
func NewMemBuffer(initialData string) *MemBuffer {
	return &MemBuffer{
		buf: []byte(initialData),
		off: 0,
	}
}

func (m *MemBuffer) Read(p []byte) (n int, err error) {
	if m.off >= int64(len(m.buf)) {
		return 0, io.EOF
	}
	n = copy(p, m.buf[m.off:])
	m.off += int64(n)
	return n, nil
}

func (m *MemBuffer) Write(p []byte) (n int, err error) {
	minLen := m.off + int64(len(p))
	if minLen > int64(len(m.buf)) {
		// Expand the slice to accommodate data written past the current length
		newBuf := make([]byte, minLen)
		copy(newBuf, m.buf)
		m.buf = newBuf
	}
	n = copy(m.buf[m.off:], p)
	m.off += int64(n)
	return n, nil
}

func (m *MemBuffer) Seek(offset int64, whence int) (int64, error) {
	var newOff int64
	switch whence {
	case io.SeekStart:
		newOff = offset
	case io.SeekCurrent:
		newOff = m.off + offset
	case io.SeekEnd:
		newOff = int64(len(m.buf)) + offset
	default:
		return 0, errors.New("stream: invalid whence")
	}

	if newOff < 0 {
		return 0, errors.New("stream: negative seek position")
	}
	m.off = newOff
	return m.off, nil
}

func (m *MemBuffer) Close() error {
	m.buf = nil // Allow garbage collection
	return nil
}

func (m *MemBuffer) Bytes() ([]byte, error) {
	// Return a defensive copy to prevent external modification
	out := make([]byte, len(m.buf))
	copy(out, m.buf)
	return out, nil
}

// ============================================================================
// OS File Wrapper
// ============================================================================

// FileWrapper wraps an *os.File to satisfy the Stream interface.
type FileWrapper struct {
	file *os.File
}

// NewFileWrapper wraps an existing *os.File.
func NewFileWrapper(f *os.File) *FileWrapper {
	return &FileWrapper{file: f}
}

func (f *FileWrapper) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f *FileWrapper) Write(p []byte) (int, error) {
	return f.file.Write(p)
}

func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *FileWrapper) Close() error {
	return f.file.Close()
}

func (f *FileWrapper) Bytes() (res []byte, err error) {
	// Preserve the current seek offset so Bytes() is side-effect free
	currPos, err := f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	defer func() {
		if _, err2 := f.file.Seek(currPos, io.SeekStart); err2 != nil && err == nil {
			err = err2
		}
	}()

	// Rewind to start and read everything
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.ReadAll(f.file)
}

func createOutput(t require.TestingT, path string) Stream {
	switch OutputType {
	case StreamBuffer:
		return NewMemBuffer("")
	case StreamFile:
		output, err := os.Create(path)
		require.NoError(t, err)
		require.NotNil(t, output)
		return NewFileWrapper(output)
	default:
		return nil
	}
}

func createInput(t require.TestingT, path string) io.ReadCloser {
	input, err := os.Open(path)
	require.NoError(t, err)
	require.NotNil(t, input)
	return input
}

func requireEqualWithDiff(t *testing.T, expected, actual []byte) {
	t.Helper()
	if bytes.Equal(expected, actual) {
		return
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(expected)),
		B:        difflib.SplitLines(string(actual)),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	t.Errorf("Mismatch (-want +got):\n%s", text)
}

func TestGenerateMarkdown(t *testing.T) {
	for _, test := range []struct {
		name      string
		file      string
		wantError bool
	}{
		{"notags", "notags", false},
		{"tags", "tags", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := createInput(t, fmt.Sprintf("testdata/%s.yaml", test.file))
			defer input.Close()
			output := createOutput(t, fmt.Sprintf("testdata/%s.md", test.file))
			defer output.Close()
			err := GenerateMarkdown(input, output)
			// Check error expectations
			if (err != nil) != test.wantError {
				t.Fatalf("GenerateMarkdown %s error = %v, wantErr %v", test.name, err, test.wantError)
			}
			// check result
			expected, err := os.ReadFile(fmt.Sprintf("testdata/%s.md", test.file))
			require.NoError(t, err)
			actual, err := output.Bytes()
			require.NoError(t, err)
			requireEqualWithDiff(t, expected, actual)
		})
	}
}
