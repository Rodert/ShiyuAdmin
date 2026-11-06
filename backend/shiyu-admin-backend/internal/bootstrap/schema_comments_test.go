package bootstrap

import "testing"

func TestMySQLCommentLiterals(t *testing.T) {
	for _, tc := range []struct {
		name, value, want string
		noBackslash       bool
	}{
		{"unicode", "系统用户表", "'系统用户表'", false},
		{"quote", "owner's files", "'owner''s files'", false},
		{"default escape mode", `path\'quoted`, `'path\\''quoted'`, false},
		{"no backslash escapes", `path\'quoted`, `'path\''quoted'`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := quoteMySQLString(tc.value, tc.noBackslash); got != tc.want {
				t.Fatalf("literal = %q, want %q", got, tc.want)
			}
		})
	}
}
