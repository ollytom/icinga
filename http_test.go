package icinga

import "testing"

func TestFilterEncode(t *testing.T) {
	expr := `match("*.example.com"), host.name) && "test" in host.groups`
	want := "filter=match%28%22%2A.example.com%22%29%2C%20host.name%29%20%26%26%20%22test%22%20in%20host.groups"
	got := filterEncode(expr)
	if want != got {
		t.Fail()
	}
	t.Logf("want %s, got %s", want, got)
}
