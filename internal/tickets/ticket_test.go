package tickets

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicketBasics(t *testing.T) {
	userID := "user1"

	t.Run("Create_and_validate_ticket", func(t *testing.T) {
		ticket := NewTicket("ticket1", userID, map[string]string{"foo": "bar"})
		assert.NotEmpty(t, ticket.ID)
		assert.Equal(t, 0, ticket.Step)
		assert.Equal(t, 0, ticket.Hops)
	})

	t.Run("Namespaced_context", func(t *testing.T) {
		ticket := NewTicket("ticket2", userID, nil)

		// Standard agentID/key
		ticket.SetNamespaced("agent1", "key", "val")
		v, ok := ticket.GetNamespaced("agent1", "key")
		assert.True(t, ok)
		assert.Equal(t, "val", v)

		// Empty agentID
		ticket.SetNamespaced("", "emptyAgent", "val")
		v2, ok2 := ticket.GetNamespaced("_", "emptyAgent")
		assert.True(t, ok2)
		assert.Equal(t, "val", v2)

		// Empty key
		ticket.SetNamespaced("agentX", "", "val2")
		v3, ok3 := ticket.GetNamespaced("agentX", "_")
		assert.True(t, ok3)
		assert.Equal(t, "val2", v3)

		// Empty agentID and key
		ticket.SetNamespaced("", "", "val3")
		v4, ok4 := ticket.GetNamespaced("_", "_")
		assert.True(t, ok4)
		assert.Equal(t, "val3", v4)
	})

	t.Run("IncrementStep_increments_step_and_hops", func(t *testing.T) {
		ticket := NewTicket("ticket3", userID, nil)
		ticket.IncrementStep(false)
		assert.Equal(t, 1, ticket.Step)
		assert.Equal(t, 1, ticket.Hops)
	})

	t.Run("Serialize_and_SerializeContext", func(t *testing.T) {
		ticket := NewTicket("ticket4", userID, map[string]string{"foo": "bar"})
		data := ticket.Serialize()
		assert.Equal(t, ticket.ID, data["id"])
		assert.Equal(t, ticket.UserID, data["user_id"])

		ctx := ticket.SerializeContext()
		assert.Equal(t, "bar", ctx[Namespaced("", "foo")])
	})

	t.Run("GetAllNamespaced_returns_all_keys_for_agent", func(t *testing.T) {
		ticket := NewTicket("ticket5", userID, nil)
		ticket.SetNamespaced("agentA", "1", "1")
		ticket.SetNamespaced("agentA", "2", "2")
		ticket.SetNamespaced("agentB", "1", "B1")

		all := ticket.GetAllNamespaced("agentA")
		assert.Len(t, all, 2)
		assert.Equal(t, "1", all["1"])
		assert.Equal(t, "2", all["2"])
	})

	t.Run("SetNamespacedWithOverwrite_blocks_collision", func(t *testing.T) {
		ticket := NewTicket("ticket6", userID, nil)
		err := ticket.SetNamespacedWithOverwrite("agentX", "key1", "val1", false)
		assert.NoError(t, err)

		err = ticket.SetNamespacedWithOverwrite("agentX", "key1", "val2", false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		err = ticket.SetNamespacedWithOverwrite("agentX", "key1", "val2", true)
		assert.NoError(t, err)
		v, ok := ticket.GetNamespaced("agentX", "key1")
		assert.True(t, ok)
		assert.Equal(t, "val2", v)
	})

	t.Run("Validate_ticket_expired_or_overhopped", func(t *testing.T) {
		ticket := NewTicket("ticket7", userID, nil)
		err := ticket.Validate()
		assert.NoError(t, err)

		ticket.Hops = ticket.MaxHops
		err = ticket.Validate()
		assert.Error(t, err)

		ticket = NewTicket("ticket8", userID, nil)
		ticket.ExpiresAt = ticket.CreatedAt.Add(-time.Hour)
		err = ticket.Validate()
		assert.Error(t, err)
	})

	t.Run("NewID_contains_parts", func(t *testing.T) {
		id := NewID("part1", "part2")
		assert.True(t, strings.Contains(id, "part1"))
		assert.True(t, strings.Contains(id, "part2"))
	})
}
