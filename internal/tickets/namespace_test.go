package tickets

import "testing"

func TestNamespaced(t *testing.T) {
	cases := []struct {
		agentID string
		field   string
		want    string
	}{
		{"agentA", "foo", "agent.agentA.foo"},
		{"x1", "state", "agent.x1.state"},
		{"A_B-C", "value123", "agent.A_B-C.value123"},
		{"", "", "agent._._"}, // updated for auto-fill
	}

	for _, c := range cases {
		got := Namespaced(c.agentID, c.field)
		if got != c.want {
			t.Errorf("Namespaced(%q, %q) = %q; want %q", c.agentID, c.field, got, c.want)
		}
	}
}
