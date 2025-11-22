package cmd

import (
	"fmt"

	"keystone/internal/tickets"

	"github.com/spf13/cobra"
)

// newTicketCmd creates the root "ticket" command and injects the store
func newTicketCmd(store *tickets.Store) *cobra.Command {
	ticketCmd := &cobra.Command{
		Use:   "ticket",
		Short: "Manage tickets",
	}

	ticketCmd.AddCommand(
		newTicketNewCmd(store),
		newTicketListCmd(store),
		newTicketDeleteCmd(store),
		newTicketPurgeCmd(store),
		newTicketCleanupCmd(store),
		newTicketInspectCmd(store),
		newTicketMonitorCmd(store),
	)

	return ticketCmd
}

// ------------------------ Helpers ------------------------

func serializeTicket(t *tickets.Ticket) map[string]interface{} {
	data := t.Serialize()
	data["context"] = t.SerializeContext()
	return data
}

func printTickets(ticketsList []*tickets.Ticket, cmd *cobra.Command) {
	if getJSONFlag(cmd) {
		out := make([]map[string]interface{}, len(ticketsList))
		for i, t := range ticketsList {
			out[i] = serializeTicket(t)
		}
		Print(out, "", cmd)
	} else {
		for _, t := range ticketsList {
			Print(nil, fmt.Sprintf("- %s (Hops: %d, Expires: %s)", t.ID, t.Hops, t.ExpiresAt), cmd)
		}
	}
}

func wrapTicketError(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}

// ------------------------ Commands ------------------------

func newTicketNewCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create a new ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			t := tickets.NewTicket(tickets.NewID("cli", "ticket", "default"), "default", nil)
			if err := store.Save(t); err != nil {
				return wrapTicketError("saving ticket", err)
			}
			Print(serializeTicket(t), fmt.Sprintf("✅ Created ticket %s for user %s", t.ID, t.UserID), cmd)
			return nil
		},
	}
}

func newTicketListCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketsList, err := store.List("default")
			if err != nil {
				return wrapTicketError("listing tickets", err)
			}
			printTickets(ticketsList, cmd)
			return nil
		},
	}
}

func newTicketDeleteCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [ticketID]",
		Short: "Delete a ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketID := args[0]
			if err := store.Delete("default", ticketID); err != nil {
				return wrapTicketError(fmt.Sprintf("deleting ticket %s", ticketID), err)
			}
			Print(map[string]string{"status": "deleted"}, fmt.Sprintf("✅ Deleted ticket %s", ticketID), cmd)
			return nil
		},
	}
}

func newTicketPurgeCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "purge",
		Short: "Purge all tickets for the default user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Purge("default"); err != nil {
				return wrapTicketError("purging tickets", err)
			}
			Print(map[string]string{"status": "purged"}, "✅ Purged all tickets", cmd)
			return nil
		},
	}
}

func newTicketCleanupCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup expired or stale tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			removed, err := store.Cleanup("default")
			if err != nil {
				return wrapTicketError("cleanup failed", err)
			}
			Print(map[string]int{"removed": removed}, fmt.Sprintf("✅ Removed %d tickets", removed), cmd)
			return nil
		},
	}
}

func newTicketInspectCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [ticketID]",
		Short: "Inspect a ticket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tkt, err := store.Load("default", args[0])
			if err != nil {
				return wrapTicketError("loading ticket", err)
			}
			Print(serializeTicket(tkt), fmt.Sprintf("Ticket %s (Hops: %d, Expires: %s)", tkt.ID, tkt.Hops, tkt.ExpiresAt), cmd)
			return nil
		},
	}
}

func newTicketMonitorCmd(store *tickets.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "monitor",
		Short: "Monitor tickets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ticketsList, err := store.List("default")
			if err != nil {
				return wrapTicketError("monitoring tickets", err)
			}
			total := len(ticketsList)
			stale := 0
			for _, t := range ticketsList {
				if t.Hops >= t.MaxHops || t.IsExpired() {
					stale++
				}
			}
			Print(map[string]int{"total_tickets": total, "stale_tickets": stale}, fmt.Sprintf("Total: %d, Stale: %d", total, stale), cmd)
			return nil
		},
	}
}
