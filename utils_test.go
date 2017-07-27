package main

////////////////////////////////////////////////////////////////////////////////

import (
	"testing"
)

////////////////////////////////////////////////////////////////////////////////

func TestHalfs(t *testing.T) {
	for _, tc := range []struct {
		value               string
		expFirst, expSecond string
	}{
		{
			value:     "ODDLENGTHSTRING",
			expFirst:  "ODDLENG",
			expSecond: "THSTRING",
		},
	} {
		first := firstHalf(tc.value)
		second := secondHalf(tc.value)

		if tc.expFirst != first {
			t.Errorf("First mismatch: expected: %#v, actual: %#v\n", tc.expFirst, first)
		}
		if tc.expSecond != second {
			t.Errorf("Second mismatch: expected: %#v, actual: %#v\n", tc.expSecond, second)
		}
		if tc.value != first+second {
			t.Errorf("Data loss from combining: expected: %#v, actual: %#v\n", tc.value, first+second)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
