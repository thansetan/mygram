package config

import "testing"

type testcase struct {
	in  string
	out bool
}

var testcases = []testcase{
	{"/", true},
	{"/api/v1/", true},
	{"/api/v1", false},
	{"/api/v1//", false},
	{"/api/v1//v2/", false},
	{"/api/v1/v2/", true},
	{"/api/v1/v2", false},
	{"", false},
	{"//", false},
	{"/api///v1///v2//", false},
	{"//api/v1/v2/", false},
	{"api", false},
	{"api/v1", false},
	{"/?/", false},
	{"/api/v1/?/", false},
	{"/api/v1/#/", false},
	{"/api/v1/#/v2/", false},
	{"/api/./?/#/", false},
	{"/api/v1/./?/#/", false},
	{"/#/", false},
}

func TestValidBasePath(t *testing.T) {
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.in, func(t *testing.T) {
			conf := App{BasePath: tc.in}
			got := conf.isValidBasePath()
			if got != tc.out {
				t.Errorf("conf.App.BasePath = %s > conf.isValidBasePath() = %t, want %t", tc.in, got, tc.out)
			}
		})
	}
}
