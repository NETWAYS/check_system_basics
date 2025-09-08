package filter

import (
	"reflect"
	"testing"
)

type TestObject struct {
	Name string
}

func (to TestObject) GetFilterableValue(val uint) string {
	switch val {
	case Name:
		return to.Name
	default:
		return to.Name
	}
}

const (
	Name = iota
)

var (
	objList = []TestObject{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
		{
			Name: "bar",
		},
		{
			Name: "baz",
		},
		{
			Name: "foo1",
		},
		{
			Name: "foo2",
		},
		{
			Name: "000",
		},
	}
)

func TestEmptyFilter(t *testing.T) {

	filterList := []string{}

	filterOpts := Options{
		EmptyFilterNoMatch: true,
	}

	testResult1, err := Filter(objList, &filterList, Name, filterOpts)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual([]TestObject{}, testResult1) {
		t.Fatalf("expected %v, got %v", []TestObject{}, testResult1)
	}
}

func TestSimpleIncludeFilter(t *testing.T) {

	filterList := []string{
		"^foo$",
	}

	filterOpts := Options{
		MatchIncludedInResult: true,
		RegexpMatching:        true,
	}

	testResult1, err := Filter(objList, &filterList, Name, filterOpts)

	testResult1Expected := []TestObject{
		{
			Name: "foo",
		},
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(testResult1Expected, testResult1) {
		t.Fatalf("expected %v, got %v", testResult1Expected, testResult1)
	}
}

func TestSimpleExcludeFilter(t *testing.T) {

	filterList := []string{
		"^foo$",
	}

	filterOpts := Options{
		MatchIncludedInResult: false,
		RegexpMatching:        true,
	}

	testResult1, err := Filter(objList, &filterList, Name, filterOpts)

	testResult1Expected := []TestObject{
		{
			Name: "bar",
		},
		{
			Name: "bar",
		},
		{
			Name: "baz",
		},
		{
			Name: "foo1",
		},
		{
			Name: "foo2",
		},
		{
			Name: "000",
		},
	}

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(testResult1Expected, testResult1) {
		t.Fatalf("expected %v, got %v", testResult1Expected, testResult1)
	}
}
