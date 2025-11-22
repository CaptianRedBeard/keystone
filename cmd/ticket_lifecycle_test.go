package cmd

import (
	"testing"

	"keystone/internal/tickets"
)

func TestCLITicketLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	tickets.TicketDir = tmpDir

	store := tickets.NewStore(tmpDir)
	namespace := "default"
	ticketID := "t1"

	// --- Step 1: Create and save initial ticket ---
	tkt := tickets.NewTicket(ticketID, namespace, nil)
	if err := store.Save(tkt); err != nil {
		t.Fatalf("failed to save initial ticket: %v", err)
	}

	// --- Step 2: Simulate first CLI run ---
	tkt.Step++
	tkt.Hops++
	if err := store.Save(tkt); err != nil {
		t.Fatalf("failed to save after first run: %v", err)
	}

	// --- Step 3: Reload ticket using correct userID/ticketID order ---
	reloaded, err := store.Load(namespace, ticketID) // ✅ fixed
	if err != nil {
		t.Fatalf("failed to reload ticket after first run: %v", err)
	}
	if reloaded.Step != 1 {
		t.Errorf("expected step=1 after first run, got %d", reloaded.Step)
	}
	if reloaded.Hops != 1 {
		t.Errorf("expected hops=1 after first run, got %d", reloaded.Hops)
	}

	// --- Step 4: Simulate second CLI run ---
	reloaded.Step++
	reloaded.Hops++
	if err := store.Save(reloaded); err != nil {
		t.Fatalf("failed to save after second run: %v", err)
	}

	// Reload again
	reloaded, err = store.Load(namespace, ticketID) // ✅ fixed
	if err != nil {
		t.Fatalf("failed to reload ticket after second run: %v", err)
	}
	if reloaded.Step != 2 {
		t.Errorf("expected step=2 after second run, got %d", reloaded.Step)
	}
	if reloaded.Hops != 2 {
		t.Errorf("expected hops=2 after second run, got %d", reloaded.Hops)
	}
}
