package channels_test

import (
	"encoding/json"
	"testing"

	"eagain.net/go/oppositus/channels"
)

type channelTest struct {
	val  channels.Channel
	json string
}

var channelTests = []channelTest{
	{channels.Stable, `"stable"`},
	{channels.Beta, `"beta"`},
	{channels.Alpha, `"alpha"`},
}

func TestChannelJSONMarshal(t *testing.T) {
	for i, test := range channelTests {
		buf, err := json.Marshal(test.val)
		if err != nil {
			t.Errorf("#%d: marshal(%v): %v", i, test.val, err)
			continue
		}
		if g, e := string(buf), test.json; g != e {
			t.Errorf("#%d: wrong json: %#q != %#q", i, g, e)
		}
	}
}

func TestChannelJSONUnmarshal(t *testing.T) {
	for i, test := range channelTests {
		var val channels.Channel
		if err := json.Unmarshal([]byte(test.json), &val); err != nil {
			t.Errorf("#%d: unmarshal(%#q): %v", i, test.json, err)
			continue
		}
		if g, e := val, test.val; g != e {
			t.Errorf("#%d: wrong value: %v != %v", i, g, e)
		}
	}
}
