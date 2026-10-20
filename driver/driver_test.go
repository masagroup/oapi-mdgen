package driver

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/pb33f/testify/require"
	"github.com/pmezard/go-difflib/difflib"
)

// update rewrites the golden markdown files instead of comparing against them:
//
//	go test ./driver -update
var update = flag.Bool("update", false, "update golden markdown files")

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
		wantError bool
	}{
		{"empty", true},
		{"invalid", false},
		{"notags", false},
		{"tags", false},
		{"params.body", false},
		{"params.body.ref", false},
		{"params.body.inline", false},
		{"params.body.content", false},
		{"params.body.allof.inline", false},
		{"params.body.oneof.inline", false},
		{"responses.simple", false},
		{"responses.array", false},
		{"responses.content.ref", false},
		{"responses.content.inline", false},
		{"responses.oneof", false},
		{"schemas.oneof.ref", false},
		{"schemas.oneof.discriminator", false},
		{"schemas.anyof", false},
		{"schemas.allof.ref", false},
		{"schemas.recursive", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			input, err := os.Open(filepath.Join("testdata", test.name+".yaml"))
			require.NoError(t, err)
			defer func() {
				require.NoError(t, input.Close())
			}()
			var output bytes.Buffer
			err = GenerateMarkdown(input, &output)
			if (err != nil) != test.wantError {
				t.Fatalf("GenerateMarkdown %s error = %v, wantErr %v", test.name, err, test.wantError)
			}
			goldenPath := filepath.Join("testdata", test.name+".md")
			if *update {
				require.NoError(t, os.WriteFile(goldenPath, output.Bytes(), 0o644))
				return
			}
			expected, err := os.ReadFile(goldenPath)
			require.NoError(t, err)
			requireEqualWithDiff(t, expected, output.Bytes())
		})
	}
}
