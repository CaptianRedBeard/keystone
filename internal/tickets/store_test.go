package tickets

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupStore creates a temporary ticket store for tests and returns a cleanup func.
func setupStore(t *testing.T) (*Store, func()) {
	t.Helper()
	tmpDir := filepath.Join(os.TempDir(), "keystone_store_test")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	store := NewStore(tmpDir)
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return store, cleanup
}

func TestStoreSaveLoadDelete(t *testing.T) {
	store, cleanup := setupStore(t)
	defer cleanup()

	ticket := NewTicket("ticket1", "user1", map[string]string{"foo": "bar"})

	// Save the ticket
	if err := store.Save(ticket); err != nil {
		t.Fatalf("failed to save ticket: %v", err)
	}

	// Load the ticket and verify
	loaded, err := store.Load("user1", "ticket1")
	if err != nil {
		t.Fatalf("failed to load ticket: %v", err)
	}
	if loaded.ID != ticket.ID {
		t.Errorf("expected ID %s, got %s", ticket.ID, loaded.ID)
	}

	// Delete the ticket and confirm removal
	if err := store.Delete("user1", "ticket1"); err != nil {
		t.Fatal(err)
	}
	list, _ := store.List("user1")
	for _, tck := range list {
		if tck.ID == "ticket1" {
			t.Error("ticket not deleted")
		}
	}
}

func TestStorePurge(t *testing.T) {
	store, cleanup := setupStore(t)
	defer cleanup()

	// Save multiple tickets for different users
	t1 := NewTicket("ticket2", "userA", nil)
	t2 := NewTicket("ticket3", "userB", nil)
	_ = store.Save(t1)
	_ = store.Save(t2)

	// Purge tickets for userA only
	if err := store.Purge("userA"); err != nil {
		t.Fatal(err)
	}

	listA, _ := store.List("userA")
	if len(listA) != 0 {
		t.Error("expected 0 tickets for userA after purge")
	}
	listB, _ := store.List("userB")
	if len(listB) != 1 {
		t.Error("expected 1 ticket for userB to remain")
	}
}

func TestStoreCleanupExpiredOrOverHopped(t *testing.T) {
	store, cleanup := setupStore(t)
	defer cleanup()

	// Create expired ticket and ticket exceeding max hops
	expired := NewTicket("ticket4", "user1", nil)
	expired.ExpiresAt = expired.CreatedAt.Add(-time.Hour)
	overHops := NewTicket("ticket5", "user1", nil)
	overHops.Hops = overHops.MaxHops

	_ = store.Save(expired)
	_ = store.Save(overHops)

	// Cleanup should remove both tickets
	removed, err := store.Cleanup("user1")
	if err != nil {
		t.Fatal(err)
	}
	if removed != 2 {
		t.Errorf("expected 2 tickets removed, got %d", removed)
	}
}
