package ldap

import "testing"

func TestExtractCN(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "CN=mygroup,OU=groups,DC=example,DC=com",
			want:  "mygroup",
		},
		{
			input: "cn=lowercase,ou=groups,dc=example,dc=com",
			want:  "lowercase",
		},
		{
			input: "CN=group with spaces,OU=groups,DC=example,DC=com",
			want:  "group with spaces",
		},
		{
			input: "CN=only",
			want:  "only",
		},
		{
			input: "OU=no-cn,DC=example,DC=com",
			want:  "OU=no-cn,DC=example,DC=com", // no CN prefix — return as-is
		},
		{
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		got := ExtractCN(tt.input)
		if got != tt.want {
			t.Errorf("ExtractCN(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
