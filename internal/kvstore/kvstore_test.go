package kvstore

import (
	"fmt"
	"reflect"
	"testing"
)

// Test is a test structure that
// we will store in our cache
type Test struct {
	ID   string `json:"ID"`
	Data string `json:"Data"`
}

// Print prints out the structure to stdio
func (t Test) Print() {
	fmt.Printf("%+v\n", t)
}

// GetID returns the ID of the data
func (t Test) GetID() string {
	return t.ID
}

func TestAddLookup(t *testing.T) {
	tc := struct {
		Elem Test
	}{
		Elem: Test{ID: "1a2a3a", Data: "test data"},
	}

	tCache := New()
	tCache.Set(tc.Elem.ID, tc.Elem)

	if want, got := tc.Elem, tCache.Get(tc.Elem.ID).(Test); want != got {
		t.Errorf("want %+v got %+v | fail\n", want, got)
	}
}

func TestAddManyDeleteOne(t *testing.T) {
	tcs := []struct {
		Elem Test
	}{
		{Elem: Test{ID: "1a2a3a", Data: "test data"}},
		{Elem: Test{ID: "1b2b3b", Data: "another test data"}},
		{Elem: Test{ID: "1c2c3c", Data: "one more test data"}},
	}

	tCache := New()
	want := make(map[string]Test)
	for _, tc := range tcs {
		tCache.Set(tc.Elem.ID, tc.Elem)
		want[tc.Elem.ID] = tc.Elem
	}
	delete(want, "1c2c3c")
	tCache.Delete("1c2c3c")
	// tCache.PRINT()
	//fmt.Printf("%+v", tCache)

	if got := tCache.Data; reflect.DeepEqual(want, got) {
		t.Errorf("want %+v got %+v | fail\n", want, got)
	}
}
