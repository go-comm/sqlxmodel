package sqlxmodel

import "testing"

func TestIfClauseAppendWhere(t *testing.T) {

	var tests = []struct {
		Param string
		Want  bool
	}{
		{Param: "", Want: false},
		{Param: "1=1", Want: true},
		{Param: "where=1", Want: true},
		{Param: "where 1=1", Want: false},
		{Param: "order=1", Want: true},
		{Param: "order by 1", Want: false},
		{Param: "left join", Want: false},
		{Param: "join", Want: false},
		{Param: "join_group", Want: true},
		{Param: "wher", Want: true},
		{Param: "where", Want: false},
	}

	for _, test := range tests {
		got := IfClauseAppendWhere(test.Param)
		if got != test.Want {
			t.Fatalf("%v, got %v, but want %v", test.Param, got, test.Want)
		}
	}

}
