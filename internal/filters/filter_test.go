package filters_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"eagain.net/go/oppositus/internal/filters"
)

type filterJSONTest struct {
	val  filters.Filters
	json string
}

var filterJSONTests = []filterJSONTest{
	{
		[]filters.Filter{
			filters.Include{Glob: filters.Glob("*")},
		},
		`["+ *"]`,
	},

	{
		[]filters.Filter{
			filters.Include{Glob: filters.Glob("foo")},
			filters.Exclude{Glob: filters.Glob("*")},
		},
		`["+ foo","- *"]`,
	},
}

func TestFilterJSONMarshal(t *testing.T) {
	for i, test := range filterJSONTests {
		buf, err := json.Marshal(test.val)
		if err != nil {
			t.Errorf("#%d: marshal(%#v): %v", i, test.val, err)
			continue
		}
		if g, e := string(buf), test.json; g != e {
			t.Errorf("#%d: wrong json: %#q != %#q", i, g, e)
		}
	}
}

func TestFilterJSONUnmarshal(t *testing.T) {
	for i, test := range filterJSONTests {
		var val filters.Filters
		if err := json.Unmarshal([]byte(test.json), &val); err != nil {
			t.Errorf("#%d: unmarshal(%#q): %v", i, test.json, err)
			continue
		}
		if g, e := val, test.val; !reflect.DeepEqual(g, e) {
			t.Errorf("#%d: wrong value: %#v != %#v", i, g, e)
		}
	}
}

func TestFilterMatch(t *testing.T) {
	tests := []struct {
		filters filters.Filters
		path    string
		want    bool
	}{

		{[]filters.Filter{
			filters.Include{Glob: filters.Glob("*")},
		}, "foo", true},

		{[]filters.Filter{
			filters.Exclude{Glob: filters.Glob("*")},
		}, "foo", false},

		{[]filters.Filter{
			filters.Include{Glob: filters.Glob("does-not-match")},
		}, "foo", true},

		{[]filters.Filter{
			filters.Exclude{Glob: filters.Glob("does-not-match")},
			filters.Include{Glob: filters.Glob("foo*")},
			filters.Exclude{Glob: filters.Glob("*")},
		}, "foobar", true},
	}

	for i, test := range tests {
		got := test.filters.Match(test.path)
		if g, e := got, test.want; g != e {
			log.Printf("#%d: mismatch: %q: %v != %v", i, test.path, g, e)
		}
	}
}
