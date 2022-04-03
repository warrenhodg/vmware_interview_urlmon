package urlchecker

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	testCases := []struct {
		url string
		up  bool
	}{
		{
			"http://httpstat.us/503",
			false,
		},
		{
			"http://httpstat.us/200",
			true,
		},
		{
			"not-a-url",
			false,
		},
		{
			"",
			false,
		},
	}

	c := NewURLChecker(nil)

	for _, testCase := range testCases {
		t.Run(testCase.url, func(t *testing.T) {
			up, _ := c.Check(context.Background(), testCase.url)
			if up != testCase.up {
				t.Fatal()
			}
		})
	}
}
