package go2way

import (
	"testing"
)

func TestGo2way(t *testing.T) {
	tests := []struct {
		query string
		input map[string]interface{}
		sign  NamedSign
		want  string
	}{
		{
			query: `select * from foo where id = /*id*/0.5 /* IF true */and bar = /*id*/ /*END*/`,
			input: map[string]interface{}{"true": true},
			sign:  Question,
			want:  `select * from foo where id = ? and bar = ?`,
		},
		{
			query: `insert into foo(id, bar) values(1,/*bar*/'bar')`,
			input: map[string]interface{}{"true": true},
			sign:  Question,
			want:  `insert into foo(id, bar) values(1,?)`,
		},
		{
			query: `insert into foo(id, bar) values(1,/*bar*/'bar')`,
			input: map[string]interface{}{"true": true},
			sign:  Dollar,
			want:  `insert into foo(id, bar) values(1,$bar)`,
		},
		{
			query: `insert into foo(id, bar) values(1,/*bar*/'bar')`,
			input: map[string]interface{}{"true": true},
			sign:  Colon,
			want:  `insert into foo(id, bar) values(1,:bar)`,
		},
	}

	for _, tt := range tests {
		got, err := Eval(tt.query, tt.input, tt.sign)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Fatalf("want %v, but %v", tt.want, got)
		}
	}
}
