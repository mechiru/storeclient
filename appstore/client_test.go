package appstore

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestLookupURL(t *testing.T) {
	for idx, c := range []struct {
		in   LookupKey
		want string
		err  error
	}{
		{
			StoreID(340368403),
			"https://itunes.apple.com/lookup?id=340368403",
			nil,
		},
		{
			BundleID("com.cookpad"),
			"https://itunes.apple.com/lookup?bundleId=com.cookpad",
			nil,
		},
		{
			LookupKey{},
			"",
			errEmptyLookupKey,
		},
	} {
		got, err := NewClient().lookupURL(c.in)
		if err != c.err {
			t.Errorf("idx=%d: got=%+v, want=%+v", idx, err, c.err)
		}
		if err != nil {
			continue
		}
		if v := got.String(); v != c.want {
			t.Errorf("idx=%d: got=%+v, want=%+v", idx, v, c.want)
		}
	}
}

// https://itunes.apple.com/lookup?id=340368403&country=JP&lang=ja_jp
func TestParseLookupResponse(t *testing.T) {
	buf, err := ioutil.ReadFile("./testdata/340368403.json")
	if err != nil {
		t.Fatal(err)
	}
	var resp LookupResponse
	if err = json.Unmarshal(buf, &resp); err != nil {
		t.Error(err)
	}
}

func TsetLookup(t *testing.T) {
	c := NewClient(Lang("ja_jp"), Country("JP"))
	if _, err := c.Lookup(context.Background(), StoreID(340368403)); err != nil {
		t.Error(err)
	}
}
