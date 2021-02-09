package main

import "testing"

func TestParseSearchQuery(t *testing.T) {
	type test struct {
		input string
		want  SearchQuery
	}

	tc := []test{
		{input: `hoge
---===--- `, want: SearchQuery{Pattern: "hoge"}},
		{input: `fuga---===---`, want: SearchQuery{Pattern: "fuga---===---"}},
		{input: `piyo
---===---
fuga`, want: SearchQuery{Pattern: "piyo"}},
		{input: `piyo
---===---
▿ includes
app/,*.jpg`, want: SearchQuery{Pattern: "piyo", Includes: "app/,*.jpg"}},
	}

	for _, tc := range tc {
		tc := tc

		q, err := ParseSearchQuery(tc.input)
		if err != nil {
			t.Error(err)
		}

		if *q != tc.want {
			t.Errorf("bad: want %+v got %+v", tc.want, *q)
		}
	}
}
