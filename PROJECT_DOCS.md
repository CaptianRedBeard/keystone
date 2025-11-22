# Project Documentation

## cmd/ticket_lifecycle_test.go
```go
```

## cmd/agent.go
```go
// newAgentCmd creates the "agent" command and subcommands
func newAgentCmd(managerProvider func() *agent.AgentManager) *cobra.Command {

// loadOrCreateTicket ensures the ticket exists
func loadOrCreateTicket(ticketID string) (*tickets.Ticket, *tickets.Store, error) {

```

## cmd/ticket.go
```go
// newTicketCmd creates the root "ticket" command and injects the store
func newTicketCmd(store *tickets.Store) *cobra.Command {

```

## cmd/config_test.go
```go
// TestConfigCommand verifies the config command output.
func TestConfigCommand(t *testing.T) {

```

## cmd/agent_register.go
```go
```

## cmd/root_test.go
```go
```

## cmd/agent_test.go
```go
// MockAgent implements a minimal agent for testing ticket handling.
type MockAgent struct {

```

## cmd/print.go
```go
// getJSONFlag returns true if the --json flag is set on the command.
func getJSONFlag(cmd *cobra.Command) bool {

// Print outputs obj as JSON if --json is set, otherwise prints msg.
func Print(obj interface{}, msg string, cmd *cobra.Command) {

// PrintError prints an error in JSON if --json is set, otherwise human-readable.
func PrintError(context, errMsg string, cmd *cobra.Command) {

```

## cmd/workflow_test.go
```go
// containsError returns true if CLI output contains an error string
func containsError(output string) bool {

// mockManagerProvider returns a manager with a single mock agent registered
func mockManagerProvider() *agent.AgentManager {

// runWorkflowCLI executes workflow commands through the CLI and captures output
func runWorkflowCLI(t *testing.T, args ...string) string {

```

## cmd/test_helpers.go
```go
// managerProvider returns a new AgentManager for tests
func managerProvider() *agent.AgentManager {

// captureOutput runs a cobra command and captures its stdout/stderr.
func captureOutput(f func(cmd *cobra.Command)) string {

// parseJSONOutput trims and parses CLI JSON output
func parseJSONOutput(output string, t *testing.T) map[string]interface{} {

// runCLITicketStep simulates a CLI agent run that increments a ticket.
func runCLITicketStep(tmpDir, ticketID string) error {

```

## cmd/usage_test.go
```go
```

## cmd/workflow.go
```go
// newWorkflowCmd creates the "workflow" CLI group and its subcommands
func newWorkflowCmd(managerProvider func() *agent.AgentManager, store *tickets.Store) *cobra.Command {

// newWorkflowRunCmd runs a workflow from a YAML file by workflow ID
func newWorkflowRunCmd(managerProvider func() *agent.AgentManager, store *tickets.Store) *cobra.Command {

```

## cmd/usage.go
```go
```

## cmd/root.go
```go
// NewRootCmd creates the root CLI command.
// managerProvider returns an AgentManager given a directory.
// configLoader returns a Config from a path (can be mocked for tests).
// out is optional stdout/stderr writer (useful for testing).
// store is optional ticket store.
func NewRootCmd(

// Execute runs the CLI root command using default loader and manager
func Execute() {

// loadManagerWithConfig loads the agent manager from a given dir or default.
func loadManagerWithConfig(dir string) *agent.AgentManager {

```

## cmd/ticket_test.go
```go
```

## cmd/config.go
```go
// newConfigCmd returns a CLI command to view/edit configuration.
// You can inject a custom load function for testing.
func newConfigCmd(loadFn func(string) (*config.Config, error)) *cobra.Command {

```

## cmd/agent_register_test.go
```go
```

## cmd/print_test.go
```go
```

## main.go
```go
```

## internal/logger/logger.go
```go
// InitDefault initializes logger using the default dev path
func InitDefault(verbose bool) error {

// Init initializes the logger system with verbosity control.
func Init(logPath string, verbose bool) error {

// InitWithWriter allows using a custom writer (useful in tests)
func InitWithWriter(logPath string, verbose bool, out io.Writer) error {

// Close should be called on program exit (or in tests)
func Close() {

// logEntry handles the actual log formatting and optional console output.
func logEntry(symbol, message string, jsonMode bool) {

```

## internal/logger/logger_test.go
```go
```

## internal/workflow/workflow.go
```go
// Save persists a workflow
func (s *Store) Save(wf Workflow) error {

// Load retrieves a workflow
func (s *Store) Load(id string) (*Workflow, error) {

// List all workflows
func (s *Store) List() ([]Workflow, error) {

// Delete removes a workflow by ID
func (s *Store) Delete(id string) error {

```

## internal/workflow/store_test.go
```go
```

## internal/workflow/engine.go
```go
// Engine coordinates workflow execution
type Engine struct {

// NewEngine creates a new workflow engine
func NewEngine(manager *agent.AgentManager, verbose bool) *Engine {

// Run executes a workflow sequentially, updating the ticket after each step
func (e *Engine) Run(ctx context.Context, wf Workflow, ticket *tickets.Ticket) ([]StepResult, error) {

// Helper: replace {{key}} in template with value
func replacePlaceholder(template, key, val string) string {

```

## internal/workflow/engine_test.go
```go
// --- Mock provider implementation ---
type MockProvider struct{}

// --- Mock agent implementation ---
type MockAgent struct {

```

## internal/agent/agent_test_helpers.go
```go
// Package agent provides helper functions and mocks for testing agents.
package agent

// MockProvider is a simple test provider returning predictable responses.
type MockProvider struct{}

// NewMockTicket returns a pre-populated ticket for testing.
func NewMockTicket() *tickets.Ticket {

// WriteYAML writes an object to a file as YAML.
func WriteYAML(t *testing.T, path string, data any) {

// BuildTestAgent constructs a simple test agent with default memory/model.
func BuildTestAgent(id, name string) Agent {

// HandleInput runs an agent's Handle method and asserts no error occurred.
func HandleInput(t *testing.T, agent Agent, input string) string {

// WriteTempAgentConfig writes a temporary agent config YAML for tests.
func WriteTempAgentConfig(t *testing.T, dir, id, provider string) string {

```

## internal/agent/config_lifecycle_edge_test.go
```go
// Package agent contains edge tests for LifecycleManager and config handling.
package agent

```

## internal/agent/config_loader.go
```go
// Package agent provides base implementations and helpers for AI agents.
package agent

// BuildAgent constructs an Agent from an AgentConfig.
func BuildAgent(cfg AgentConfig) Agent {

// LoadAgentsFromConfig scans a directory and loads all YAML agent configs into the manager.
func LoadAgentsFromConfig(manager *AgentManager, configDir string) error {

// LoadDefaultAgent ensures a fallback dummy agent exists in the manager.
func LoadDefaultAgent(manager *AgentManager) {

```

## internal/agent/manager.go
```go
// Package agent provides base implementations and helpers for AI agents.
package agent

// AgentManager manages a thread-safe collection of Agents.
type AgentManager struct {

// NewManager creates and returns a new AgentManager.
func NewManager() *AgentManager {

// Register adds a new agent to the manager.
func (m *AgentManager) Register(a Agent) error {

// Get retrieves an agent by ID.
func (m *AgentManager) Get(id string) (Agent, error) {

// List returns a slice of all registered agents.
func (m *AgentManager) List() []Agent {

// Unregister removes an agent by ID.
func (m *AgentManager) Unregister(id string) error {

```

## internal/agent/base_test.go
```go
// Package agent contains tests for AgentBase.
package agent

// TestAgentBase_GettersAndHandle verifies getters and Handle method work correctly.
func TestAgentBase_GettersAndHandle(t *testing.T) {

// TestAgentBase_Handle_NoProvider verifies that Handle returns error when provider is nil.
func TestAgentBase_Handle_NoProvider(t *testing.T) {

```

## internal/agent/config_test.go
```go
```

## internal/agent/manager_test.go
```go
```

## internal/agent/agent_test.go
```go
// Package agent contains core agent tests.
package agent

// TestAgent_Handle verifies agent handles input and returns a response.
func TestAgent_Handle(t *testing.T) {

// TestAgent_Getters verifies all getters return expected values.
func TestAgent_Getters(t *testing.T) {

// TestAgent_HandleWithNilProvider ensures agent returns error when provider is nil.
func TestAgent_HandleWithNilProvider(t *testing.T) {

```

## internal/agent/lifecycle_test.go
```go
```

## internal/agent/interface.go
```go
```

## internal/agent/lifecycle.go
```go
// Package agent provides base implementations and helpers for AI agents.
package agent

// LifecycleManager manages agent configurations and provider resolution.
type LifecycleManager struct {

// NewLifecycleManager creates a new LifecycleManager with optional config directory and provider map.
func NewLifecycleManager(configDir string, providersMap map[string]providers.Provider) *LifecycleManager {

// Manager returns the internal AgentManager.
func (lm *LifecycleManager) Manager() *AgentManager {

// RegisterProvider adds a provider under the specified name.
func (lm *LifecycleManager) RegisterProvider(name string, p providers.Provider) {

// ResolveProvider retrieves a registered provider by name.
func (lm *LifecycleManager) ResolveProvider(name string) (providers.Provider, error) {

// SaveOrMergeConfig saves an AgentConfig to YAML in the config directory.
func (lm *LifecycleManager) SaveOrMergeConfig(cfg AgentConfig) error {

// LoadAgent loads a single agent by ID and registers it.
func (lm *LifecycleManager) LoadAgent(agentID string) error {

// LoadAgentsFromDir loads all YAML agent configs from the config directory.
func (lm *LifecycleManager) LoadAgentsFromDir() error {

// loadAgentConfig reads a YAML file and unmarshals it into an AgentConfig.
func (lm *LifecycleManager) loadAgentConfig(path string) (AgentConfig, error) {

```

## internal/agent/base.go
```go
// Package agent provides base implementations and helpers for AI agents.
package agent

// AgentBase is a simple base implementation of an Agent.
type AgentBase struct {

// AgentOption is a functional option to configure AgentBase.
type AgentOption func(*AgentBase)

// WithPromptTemplate sets the agent's prompt template.
func WithPromptTemplate(tpl string) AgentOption {

// WithParameters sets the agent's parameters.
func WithParameters(params map[string]string) AgentOption {

// WithLogging enables or disables agent logging.
func WithLogging(enabled bool) AgentOption {

// NewAgent creates a new AgentBase with defaults applied.
func NewAgent(id, name, description string, provider providers.Provider, model, memory string, opts ...AgentOption) Agent {

// ID returns the agent's ID.
func (a *AgentBase) ID() string { return a.id }

// Name returns the agent's display name.
func (a *AgentBase) Name() string { return a.name }

// Description returns the agent's description.
func (a *AgentBase) Description() string { return a.description }

// Memory returns the agent's memory identifier.
func (a *AgentBase) Memory() string { return a.memory }

// DefaultModel returns the agent's default model.
func (a *AgentBase) DefaultModel() string { return a.model }

// Provider returns the agent's provider instance.
func (a *AgentBase) Provider() providers.Provider { return a.provider }

// PromptTemplate returns the agent's prompt template.
func (a *AgentBase) PromptTemplate() string { return a.promptTemplate }

// Parameters returns the agent's parameters map.
func (a *AgentBase) Parameters() map[string]string { return a.parameters }

// LoggingEnabled returns true if logging is enabled.
func (a *AgentBase) LoggingEnabled() bool { return a.logging }

// Handle processes input using the agent's provider.
func (a *AgentBase) Handle(ctx context.Context, input string, t *tickets.Ticket) (string, error) {

```

## internal/agent/agent_edge_test.go
```go
// Package agent contains tests for agent behaviors, including edge cases.
package agent

// ErrProviderFail is used to simulate a failing provider.
var ErrProviderFail = errors.New("provider failure")

// FailingProvider simulates a provider that always fails.
type FailingProvider struct{}

// TestAgent_HandleWithProviderError ensures the agent propagates provider errors.
func TestAgent_HandleWithProviderError(t *testing.T) {

// TestAgent_HandleNilProvider ensures handling an agent with nil provider returns error.
func TestAgent_HandleNilProvider(t *testing.T) {

// TestAgent_EmptyIDNameDefaults ensures default values are set for empty ID/Name.
func TestAgent_EmptyIDNameDefaults(t *testing.T) {

// TestAgent_ParameterMerging ensures options correctly override parameters.
func TestAgent_ParameterMerging(t *testing.T) {

// TestAgent_PromptTemplateOption ensures prompt template option is applied.
func TestAgent_PromptTemplateOption(t *testing.T) {

```

## internal/agent/config.go
```go
// Package agent provides base implementations and helpers for AI agents.
package agent

// AgentConfig defines the structure of an agent YAML configuration.
type AgentConfig struct {

// Merge merges another AgentConfig (src) into this one, prioritizing non-empty fields from src.
func (dst *AgentConfig) Merge(src AgentConfig) {

// Validate ensures required fields are set correctly and applies defaults.
func (cfg *AgentConfig) Validate() error {

```

## internal/config/config_test.go
```go
```

## internal/config/config.go
```go
// ErrNoConfig is returned when the config file does not exist.
var ErrNoConfig = errors.New("config file not found")

// Config holds Keystone configuration values.
type Config struct {

// New returns a config populated with defaults, optionally overridden by environment variables.
func New() *Config {

// Load reads and unmarshals a YAML config from disk.
func Load(path string) (*Config, error) {

// Save marshals and writes the config to the given path.
func Save(path string, cfg *Config) error {

```

## internal/providers/venice/venice_test.go
```go
```

## internal/providers/venice/venice.go
```go
// GenerateResponse simulates an API call to the Venice service.
// For Phase 1, this is mocked to return a test response.
func (v *VeniceProvider) GenerateResponse(ctx context.Context, prompt string, model string) (string, error) {

// UsageInfo returns mock usage data.
func (v *VeniceProvider) UsageInfo() (providers.Usage, error) {

```

## internal/providers/provider.go
```go
```

## internal/tickets/ticket.go
```go
// Default ticket storage path, can be overridden with KEYSTONE_TICKET_DIR
var TicketDir = filepath.Join(os.Getenv("HOME"), ".keystone", "tickets")

// Ticket represents a unit of work tracked across agents.
type Ticket struct {

// OnHandoffHook is an optional function to receive handoff events (Phase 5+ GUI/monitoring).
var OnHandoffHook func(t *Ticket, nextAgentID string)

// NewTicket constructs a new Ticket with optional context.
func NewTicket(id, userID string, ctx interface{}) *Ticket {

// NewID generates a unique ticket ID using optional string parts and current timestamp.
func NewID(parts ...string) string {

// Validate checks if the ticket is expired or exceeded max hops.
func (t *Ticket) Validate() error {

// IsExpired returns true if the ticket has expired.
func (t *Ticket) IsExpired() bool {

// GetNamespaced retrieves a namespaced value for a specific agent.
func (t *Ticket) GetNamespaced(agentID, key string) (string, bool) {

// SetNamespaced sets a namespaced value for a specific agent.
func (t *Ticket) SetNamespaced(agentID, key, value string) {

// GetAllNamespaced returns all key-value pairs for a specific agent.
func (t *Ticket) GetAllNamespaced(agentID string) map[string]string {

// IncrementStep increments the ticket's step and hops, optionally logging.
func (t *Ticket) IncrementStep(verbose bool) {

// Handoff transfers the ticket to the next agent.
// It validates TTL/MaxHops, increments step and hops, and ensures
// the next agent has an initialized context namespace.
// Returns an error if the ticket cannot be handed off.
func (t *Ticket) Handoff(nextAgentID string) error {

// Serialize returns a snapshot of ticket fields.
func (t *Ticket) Serialize() map[string]interface{} {

// SerializeContext returns a copy of the ticket's context.
func (t *Ticket) SerializeContext() map[string]interface{} {

// Store handles file-based persistence of tickets.
type Store struct {

// NewStore creates a new ticket store at the given directory.
func NewStore(dir string) *Store {

// Save persists a ticket to disk.
func (s *Store) Save(t *Ticket) error {

// Load retrieves a ticket by user ID and ticket ID.
func (s *Store) Load(userID, ticketID string) (*Ticket, error) {

// List returns all tickets for a user or all users if "all" is passed.
func (s *Store) List(userID string) ([]*Ticket, error) {

// Delete removes a ticket from storage.
func (s *Store) Delete(userID, ticketID string) error {

// Purge deletes all tickets for a given user.
func (s *Store) Purge(userID string) error {

// Cleanup removes expired or over-hopped tickets for a user, returning count of removed tickets.
func (s *Store) Cleanup(userID string) (int, error) {

```

## internal/tickets/ticket_handoff_test.go
```go
// TestTicketHandoffBasic tests that Handoff increments Step and Hops
func TestTicketHandoffBasic(t *testing.T) {

// TestTicketHandoffMaxHops ensures handoff fails when MaxHops exceeded
func TestTicketHandoffMaxHops(t *testing.T) {

// TestTicketHandoffExpired ensures handoff fails for expired tickets
func TestTicketHandoffExpired(t *testing.T) {

// TestTicketHandoffHook tests OnHandoffHook is called
func TestTicketHandoffHook(t *testing.T) {

```

## internal/tickets/ticket_test.go
```go
// setupStore creates a temporary ticket store for tests and returns a cleanup func.
func setupStore(t *testing.T) (*Store, func()) {

// TestTicketLifecycle contains core ticket persistence and utility tests.
func TestTicketLifecycle(t *testing.T) {

// TestTicketHandoff contains Phase 4 tests for handoff behavior (step/hops, TTL, hooks).
func TestTicketHandoff(t *testing.T) {

```

## internal/usage/usage_test.go
```go
```

## internal/usage/usage.go
```go
// Entry represents a single usage event.
type Entry struct {

// Tracker keeps track of usage events.
type Tracker struct {

// Summary holds aggregated usage info.
type Summary struct {

// NewTracker creates a usage tracker instance.
func NewTracker() *Tracker {

// Record adds a new usage entry to the tracker.
func (t *Tracker) Record(agentID, provider string, tokens int) Entry {

// Summary aggregates usage data.
func (t *Tracker) Summary() Summary {

// List returns a copy of all usage entries.
func (t *Tracker) List() []Entry {

// generateID produces a simple timestamp-based unique ID.
func generateID() string {

```

