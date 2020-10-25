package utils

import (
	"testing"
)

func TestHash(t *testing.T) {
	tc := struct {
		ts    string
		pvkey string
		pbkey string
	}{
		ts:    "1",
		pvkey: "abcd",
		pbkey: "1234",
	}

	want := "ffd275c5130566a2916217b101f26150"

	if got := GetAPIKeyHash(tc.ts, tc.pvkey, tc.pbkey); want != got {
		t.Errorf("want [%s] got [%s] | fail\n", want, got)
	}
}
