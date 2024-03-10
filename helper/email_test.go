package helper_test

import (
	"final-project/helper"
	"testing"
)

type testcase struct {
	email string
	out   bool
}

var testcases = []testcase{
	{"budi.wow.keren@gmail.com", true},
	{"budi..wow@gmail.com", false}, // local part of address must not contain consecutive periods (RFC 5322)
	{"budi@wow.keren.sekali.com", true},
	{"budi@wow.keren..sekali.com", false}, // label in the domain part can't be empty (RFC 1035)
	{"budi.wow+awwww@wow.keren.sekali.com", true},
	{"budi.wow%aww@wow.keren.sekali.com", true},
	{"budi@wow.1aa", false}, // TLD must start with a letter (RFC 1123)
	{"budi@wow.aa1", true},
}

func TestAsaskevichGoValidatorPattern(t *testing.T) {
	emailValidator := helper.IsValidEmailRegex(helper.AsaskevichGoValidatorEmailPattern)
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := emailValidator(tt.email); got != tt.out {
				t.Errorf("[Asaskevich] IsValidEmailRegex(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func TestGoPlaygroundValidatorPattern(t *testing.T) {
	emailValidator := helper.IsValidEmailRegex(helper.GoPlaygroundValidatorEmailPattern)
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := emailValidator(tt.email); got != tt.out {
				t.Errorf("[GoPlayground] IsValidEmailRegex(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func TestStackOverflowPattern(t *testing.T) {
	emailValidator := helper.IsValidEmailRegex(helper.StackOverflowEmailPattern)
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := emailValidator(tt.email); got != tt.out {
				t.Errorf("[StackOverflow] IsValidEmailRegex(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func TestRegexLibPattern(t *testing.T) {
	emailValidator := helper.IsValidEmailRegex(helper.RegexLibEmailPattern)
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := emailValidator(tt.email); got != tt.out {
				t.Errorf("[RegexLib] IsValidEmailRegex(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func TestGoMailPKG(t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := helper.IsValidEmail(tt.email, true); got != tt.out {
				t.Errorf("[GoMailPKG] IsValidEmail(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func TestGoMailPKGNoLocal(t *testing.T) {
	for _, tt := range testcases {
		t.Run(tt.email, func(t *testing.T) {
			if got := helper.IsValidEmail(tt.email, false); got != tt.out {
				t.Errorf("[GoMailPKGNoLocal] IsValidEmail(%s) = %v, want %v", tt.email, got, tt.out)
			}
		})
	}
}

func BenchmarkStackOverflowPattern(b *testing.B) {
	emailValidator := helper.IsValidEmailRegex(helper.StackOverflowEmailPattern)
	for i := 0; i < b.N; i++ {
		emailValidator("budi.keren.sekali.anjay+tiktok@wowkeren.com.uk.qq.zz.wow.awww.com")
		emailValidator("budi...wow@gmail.com")
		emailValidator("budi@wow..gmail.com")
		emailValidator("budi@wow.1aa")
	}
}

func BenchmarkGoMailPKGNoLocal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		helper.IsValidEmail("budi.keren.sekali.anjay+tiktok@wowkeren.com.uk.qq.zz.wow.awww.com", false)
		helper.IsValidEmail("budi...wow@gmail.com", false)
		helper.IsValidEmail("budi@wow..gmail.com", false)
		helper.IsValidEmail("budi@wow.1aa", false)
	}
}
