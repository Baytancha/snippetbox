package main

import (
	"testing"
	"time"

	"github.com/Baytancha/snip56/internal/assert"
)

// func TestHumanDate(t *testing.T) {
// 	// Initialize a new time.Time object and pass it to the humanDate function.
// 	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
// 	hd := humanDate(tm)

//		// Check that the output from the humanDate function is in the format we
//		// expect. If it isn't what we expect, use the t.Errorf() function to
//		// indicate that the test has failed and log the expected and actual
//		// values.
//		if hd != "17 Mar 2022 at 10:15" {
//			t.Errorf("got %q; want %q", hd, "17 Mar 2022 at 10:15")
//		}
//	}

func TestExample(t *testing.T) {
	t.Run("Example sub-test 1", func(t *testing.T) {
		// Do a test.
	})

	t.Run("Example sub-test 2", func(t *testing.T) {
		// Do another test.
	})

	t.Run("Example sub-test 3", func(t *testing.T) {
		// And another...
	})
}

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	// (the want field).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// and anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)

			// if hd != tt.want {
			// 	t.Errorf("got %q; want %q", hd, tt.want)
			// }
		})
	}
}
