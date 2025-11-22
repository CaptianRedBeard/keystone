package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgentManager_RegisterGetListUnregister(t *testing.T) {
	manager := NewManager()
	a1 := BuildTestAgent("a1", "Agent 1")
	a2 := BuildTestAgent("a2", "Agent 2")

	// Register
	require.NoError(t, manager.Register(a1))
	require.NoError(t, manager.Register(a2))

	// Duplicate register
	require.Error(t, manager.Register(a1))

	// Get
	got, err := manager.Get("a1")
	require.NoError(t, err)
	require.Equal(t, a1.Name(), got.Name())

	// List
	list := manager.List()
	require.Len(t, list, 2)

	// Unregister
	require.NoError(t, manager.Unregister("a1"))
	require.Error(t, manager.Unregister("nonexistent"))
}
