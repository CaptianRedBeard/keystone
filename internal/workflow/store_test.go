package workflow

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowStoreLifecycle(t *testing.T) {
	// Temporary directory for testing
	tmpDir := filepath.Join(os.TempDir(), "keystone_workflows_test")
	os.MkdirAll(tmpDir, 0o755)
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	wf1 := Workflow{
		ID: "wf1",
		Steps: []Step{
			{AgentID: "a1", Input: "input1"},
		},
	}
	wf2 := Workflow{
		ID: "wf2",
		Steps: []Step{
			{AgentID: "a2", Input: "input2"},
		},
	}

	// Save workflows (pass by value)
	err := store.Save(wf1)
	assert.NoError(t, err)

	err = store.Save(wf2)
	assert.NoError(t, err)

	// Load workflow
	loaded, err := store.Load("wf1")
	assert.NoError(t, err)
	assert.Equal(t, wf1.ID, loaded.ID)
	assert.Equal(t, wf1.Steps[0].AgentID, loaded.Steps[0].AgentID)

	// List workflows
	list, err := store.List()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// Delete workflow
	err = store.Delete("wf1")
	assert.NoError(t, err)
	list, _ = store.List()
	assert.Len(t, list, 1)
	assert.Equal(t, "wf2", list[0].ID)
}
