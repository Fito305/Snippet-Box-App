package main

import (
	"testing"
	"time"

	"snippetbox.felipeacosta.net/internal/assert"
)


func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expect output
	// (the want field).
	tests := []struct {
		name string
		tm time.Time
		want string
	}{
	{
		name: "UTC",
		tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
		want: "17 Mar 2022 at 10:15",
	},
	{
		name: "Empty",
		tm :	time.Time{},
		want: "",
	},
	{
		name: "CET",
		tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
		want: "17 Mar 2022 at 09:15",
	},
  }

  // Loop over the test cases.
  for _, tt := range tests {
	  // Use the t.Run() function to run a sub-test for each test case. The first parameter to this is the name of the test 
	  // (which is used to identify the sub-test in any log output) and the second parameter is and anonymous function containg the actual test for each case.
	  t.Run(tt.name, func(t *testing.T) {
		  hd := humanDate(tt.tm)

		  // Use the new assert.Equal() helper to compare the expected and actual values. 
		  assert.Equal(t, hd, tt.want)
	  })
  }
}


// func TestHumanDate(t *testing.T) {
// 	// Initialize a new time.Time object and pass it to the humanDate function.
// 	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
// 	hd := humanDate(tm)
//
// 	// Check that the output from the humanDate function is in the format we
// 	// expect. If it isn't what we expect, use the t.Errorf() function to 
// 	// indicate that the test has failed and log the expected and actual values.
// 	if hd != "17 Mar 2022 at 10:15" {
// 		t.Errorf("got %q; want %q", hd, "17 Mar 2022 at 10:15")
// 	}
// }

// This pattern is the basic one that you'll use for nearly all thest that you write in Go. The
// important things to take away are:
// - The test is just regular Go code, which calls the `humanDate()` function and checks that the result matches what we expect. 
// - Your unit test are contained in a normal Go function with the signature func(*testing.T).
// - To be valid unit test the name of this function must begin with the word Test. Typically this is then followed by the name of the function, method 
// 	 type that you're testing to help make it obvious at a glance what is being tested.
// - You can use the t.Errorf() function to mark a test as failed and log a descriptive message about the failure. It's important to note that calling
//   t.Errorf() doesn't stop execution of your test -- after you call it Go will continue executing any remaining test code as normal.
