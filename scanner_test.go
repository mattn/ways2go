package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: `select * from foo where id = /*a*/'foo bar'`,
			want:  []string{"select", "*", "from", "foo", "where", "id", "=", "/*", "a", "*/", "'foo bar'"},
		},
		{
			input: `insert into foo(id, text) values(0, 'foo'/*a*/)`,
			want:  []string{"insert", "into", "foo", "(", "id", ",", "text", ")", "values", "(", "0", ",", "'foo'", "/*", "a", "*/", ")"},
		},
	}

	for _, tt := range tests {
		var got []string
		scan := NewScanner(strings.NewReader(tt.input))
		for scan.Scan() {
			if scan.Token() != SPACE {
				got = append(got, scan.Text())
			}
		}
		if !cmp.Equal(got, tt.want) {
			t.Fatalf("want %v, but got %v", tt.want, got)
		}
	}
}
