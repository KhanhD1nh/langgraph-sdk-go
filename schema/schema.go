package schema

import (
	"time"
)

// Json represents a JSON-like structure, which can be nil or a map with string keys and any values
type Json map[string]any

// RunStatus represents the status of a run
type RunStatus string

const (
	RunStatusPending     RunStatus = "pending"     // The run is waiting to start
	RunStatusRunning     RunStatus = "running"     // The run is currently executing
	RunStatusError       RunStatus = "error"       // The run encountered an error and stopped
	RunStatusSuccess     RunStatus = "success"     // The run completed successfully
	RunStatusTimeout     RunStatus = "timeout"     // The run exceeded its time limit
	RunStatusInterrupted RunStatus = "interrupted" // The run was manually stopped or interrupted
)

// ThreadStatus represents the status of a thread
type ThreadStatus string

const (
	ThreadStatusIdle        ThreadStatus = "idle"        // The thread is not currently processing any task
	ThreadStatusBusy        ThreadStatus = "busy"        // The thread is actively processing a task
	ThreadStatusInterrupted ThreadStatus = "interrupted" // The thread's execution was interrupted
	ThreadStatusError       ThreadStatus = "error"       // An exception occurred during task processing
)

// StreamMode defines the mode of streaming
type StreamMode string

const (
	StreamModeValues        StreamMode = "values"         // Stream only the values
	StreamModeMessages      StreamMode = "messages"       // Stream complete messages
	StreamModeUpdates       StreamMode = "updates"        // Stream updates to the state
	StreamModeEvents        StreamMode = "events"         // Stream events occurring during execution
	StreamModeDebug         StreamMode = "debug"          // Stream detailed debug information
	StreamModeCustom        StreamMode = "custom"         // Stream custom events
	StreamModeMessagesTuple StreamMode = "messages-tuple" // Stream messages as tuples
)

// DisconnectMode specifies behavior on disconnection
type DisconnectMode string

const (
	DisconnectModeCancel   DisconnectMode = "cancel"   // Cancel the operation on disconnection
	DisconnectModeContinue DisconnectMode = "continue" // Continue the operation even if disconnected
)

// MultitaskStrategy defines how to handle multiple tasks
type MultitaskStrategy string

const (
	MultitaskStrategyReject    MultitaskStrategy = "reject"    // Reject new tasks when busy
	MultitaskStrategyInterrupt MultitaskStrategy = "interrupt" // Interrupt current task for new ones
	MultitaskStrategyRollback  MultitaskStrategy = "rollback"  // Roll back current task and start new one
	MultitaskStrategyEnqueue   MultitaskStrategy = "enqueue"   // Queue new tasks for later execution
)

// OnConflictBehavior specifies behavior on conflict
type OnConflictBehavior string

const (
	OnConflictBehaviorRaise     OnConflictBehavior = "raise"      // Raise an exception when a conflict occurs
	OnConflictBehaviorDoNothing OnConflictBehavior = "do_nothing" // Ignore conflicts and proceed
)

// OnCompletionBehavior defines action after completion
type OnCompletionBehavior string

const (
	OnCompletionBehaviorDelete OnCompletionBehavior = "delete" // Delete resources after completion
	OnCompletionBehaviorKeep   OnCompletionBehavior = "keep"   // Retain resources after completion
)

// All represents a wildcard or 'all' selector
type All string

const AllWildcard All = "*"

// IfNotExists specifies behavior if the thread doesn't exist
type IfNotExists string

const (
	IfNotExistsCreate IfNotExists = "create" // Create a new thread if it doesn't exist
	IfNotExistsReject IfNotExists = "reject" // Reject the operation if the thread doesn't exist
)

// CancelAction specifies the action to take when cancelling the run
type CancelAction string

const (
	CancelActionInterrupt CancelAction = "interrupt" // Simply cancel the run
	CancelActionRollback  CancelAction = "rollback"  // Cancel the run, then delete the run and associated checkpoints
)

// Config represents configuration options for a call
type Config struct {
	Tags           []string               `json:"tags,omitempty"`            // Tags for this call and any sub-calls
	RecursionLimit int                    `json:"recursion_limit,omitempty"` // Maximum number of times a call can recurse
	Configurable   map[string]interface{} `json:"configurable,omitempty"`    // Runtime values for attributes previously made configurable
}

// Checkpoint represents a checkpoint in the execution process
type Checkpoint struct {
	ThreadID      string                 `json:"thread_id"`                // Unique identifier for the thread associated with this checkpoint
	CheckpointNS  string                 `json:"checkpoint_ns"`            // Namespace for the checkpoint, used for organization and retrieval
	CheckpointID  *string                `json:"checkpoint_id,omitempty"`  // Optional unique identifier for the checkpoint itself
	CheckpointMap map[string]interface{} `json:"checkpoint_map,omitempty"` // Optional dictionary containing checkpoint-specific data
}

// GraphSchema defines the structure and properties of a graph
type GraphSchema struct {
	GraphID      string `json:"graph_id"`                // The ID of the graph
	InputSchema  *Json  `json:"input_schema,omitempty"`  // The schema for the graph input
	OutputSchema *Json  `json:"output_schema,omitempty"` // The schema for the graph output
	StateSchema  *Json  `json:"state_schema,omitempty"`  // The schema for the graph state
	ConfigSchema *Json  `json:"config_schema,omitempty"` // The schema for the graph config
}

// Graph represents a graph with additional properties
type Graph struct {
	Nodes []Node `json:"nodes"` // The nodes in the graph
	Edges []Edge `json:"edges"` // The edges in the graph
}

// Node represents a node in a graph
type Node struct {
	ID   string                 `json:"id"`   // The ID of the node
	Type string                 `json:"type"` // The type of the node
	Data map[string]interface{} `json:"data"` // The data associated with the node
}

// Edge represents an edge in a graph
type Edge struct {
	Source string `json:"source"` // The source node ID
	Target string `json:"target"` // The target node ID
}

// Subgraphs is a map of graph names to their schemas
type Subgraphs map[string]GraphSchema

// AssistantBase is the base model for an assistant
type AssistantBase struct {
	AssistantID string    `json:"assistant_id"`          // The ID of the assistant
	GraphID     string    `json:"graph_id"`              // The ID of the graph
	Config      Config    `json:"config"`                // The assistant config
	CreatedAt   time.Time `json:"created_at"`            // The time the assistant was created
	Metadata    Json      `json:"metadata"`              // The assistant metadata
	Version     int       `json:"version"`               // The version of the assistant
	Name        string    `json:"name"`                  // The name of the assistant
	Description *string   `json:"description,omitempty"` // The description of the assistant
}

// AssistantVersion represents a specific version of an assistant
type AssistantVersion struct {
	AssistantBase
}

// Assistant represents an assistant with additional properties
type Assistant struct {
	AssistantBase
	UpdatedAt time.Time `json:"updated_at"` // The last time the assistant was updated
}

// InterruptWhen defines when an interrupt occurred
type InterruptWhen string

const InterruptWhenDuring InterruptWhen = "during"

// Interrupt represents an interruption in the execution flow
type Interrupt struct {
	Value     interface{}   `json:"value,omitempty"`     // The value associated with the interrupt
	When      InterruptWhen `json:"when,omitempty"`      // When the interrupt occurred
	Resumable bool          `json:"resumable,omitempty"` // Whether the interrupt can be resumed
	NS        []string      `json:"ns,omitempty"`        // Optional namespace for the interrupt
}

// Thread represents a conversation thread
type Thread struct {
	ThreadID   string                 `json:"thread_id"`  // The ID of the thread
	CreatedAt  time.Time              `json:"created_at"` // The time the thread was created
	UpdatedAt  time.Time              `json:"updated_at"` // The last time the thread was updated
	Metadata   Json                   `json:"metadata"`   // The thread metadata
	Status     ThreadStatus           `json:"status"`     // The status of the thread
	Values     Json                   `json:"values"`     // The current state of the thread
	Interrupts map[string][]Interrupt `json:"interrupts"` // Interrupts which were thrown in this thread
}

// ThreadTask represents a task within a thread
type ThreadTask struct {
	ID         string       `json:"id"`                   // Task ID
	Name       string       `json:"name"`                 // Task name
	Error      *string      `json:"error,omitempty"`      // Error message, if any
	Interrupts []Interrupt  `json:"interrupts"`           // List of interrupts that occurred during task execution
	Checkpoint *Checkpoint  `json:"checkpoint,omitempty"` // Associated checkpoint, if any
	State      *ThreadState `json:"state,omitempty"`      // Current thread state
	Result     Json         `json:"result,omitempty"`     // Task result
}

// ThreadState represents the state of a thread
type ThreadState struct {
	Values           interface{}  `json:"values"`                      // The state values
	Next             []string     `json:"next"`                        // The next nodes to execute
	Checkpoint       Checkpoint   `json:"checkpoint"`                  // The ID of the checkpoint
	Metadata         Json         `json:"metadata"`                    // Metadata for this state
	CreatedAt        *string      `json:"created_at,omitempty"`        // Timestamp of state creation
	ParentCheckpoint *Checkpoint  `json:"parent_checkpoint,omitempty"` // The ID of the parent checkpoint
	Tasks            []ThreadTask `json:"tasks"`                       // Tasks to execute in this step
}

// ThreadUpdateStateResponse represents the response from updating a thread's state
type ThreadUpdateStateResponse struct {
	Checkpoint Checkpoint `json:"checkpoint"` // Checkpoint of the latest state
}

// Run represents a single execution run
type Run struct {
	RunID             string            `json:"run_id"`             // The ID of the run
	ThreadID          string            `json:"thread_id"`          // The ID of the thread
	AssistantID       string            `json:"assistant_id"`       // The assistant that was used for this run
	CreatedAt         time.Time         `json:"created_at"`         // The time the run was created
	UpdatedAt         time.Time         `json:"updated_at"`         // The last time the run was updated
	Status            RunStatus         `json:"status"`             // The status of the run
	Metadata          Json              `json:"metadata"`           // The run metadata
	MultitaskStrategy MultitaskStrategy `json:"multitask_strategy"` // Strategy to handle concurrent runs on the same thread
}

// Cron represents a scheduled task
type Cron struct {
	CronID    string     `json:"cron_id"`             // The ID of the cron
	ThreadID  *string    `json:"thread_id,omitempty"` // The ID of the thread
	EndTime   *time.Time `json:"end_time,omitempty"`  // The end date to stop running the cron
	Schedule  string     `json:"schedule"`            // The schedule to run, cron format
	CreatedAt time.Time  `json:"created_at"`          // The time the cron was created
	UpdatedAt time.Time  `json:"updated_at"`          // The last time the cron was updated
	Payload   Json       `json:"payload"`             // The run payload to use for creating new run
}

// RunCreate defines the parameters for initiating a background run
type RunCreate struct {
	ThreadID          *string            `json:"thread_id,omitempty"`          // The identifier of the thread to run
	AssistantID       string             `json:"assistant_id"`                 // The identifier of the assistant to use for this run
	Input             Json               `json:"input,omitempty"`              // Initial input data for the run
	Metadata          Json               `json:"metadata,omitempty"`           // Additional metadata to associate with the run
	Config            *Config            `json:"config,omitempty"`             // Configuration options for the run
	CheckpointID      *string            `json:"checkpoint_id,omitempty"`      // The identifier of a checkpoint to resume from
	InterruptBefore   []string           `json:"interrupt_before,omitempty"`   // List of node names to interrupt execution before
	InterruptAfter    []string           `json:"interrupt_after,omitempty"`    // List of node names to interrupt execution after
	Webhook           *string            `json:"webhook,omitempty"`            // URL to send webhook notifications about the run's progress
	MultitaskStrategy *MultitaskStrategy `json:"multitask_strategy,omitempty"` // Strategy for handling concurrent runs on the same thread
}

// Item represents a single document or data entry in the graph's Store
type Item struct {
	Namespace []string               `json:"namespace"`  // The namespace of the item
	Key       string                 `json:"key"`        // The unique identifier of the item within its namespace
	Value     map[string]interface{} `json:"value"`      // The value stored in the item
	CreatedAt time.Time              `json:"created_at"` // The timestamp when the item was created
	UpdatedAt time.Time              `json:"updated_at"` // The timestamp when the item was last updated
}

// ListNamespaceResponse is the response structure for listing namespaces
type ListNamespaceResponse struct {
	Namespaces [][]string `json:"namespaces"` // A list of namespace paths, where each path is a list of strings
}

// SearchItem is an Item with an optional relevance score from search operations
type SearchItem struct {
	Item
	Score *float64 `json:"score,omitempty"` // Relevance/similarity score
}

// SearchItemsResponse is the response structure for searching items
type SearchItemsResponse struct {
	Items []SearchItem `json:"items"` // A list of items matching the search criteria
}

// StreamPart represents a part of a stream response
type StreamPart struct {
	Event    string `json:"event"`    // The type of event for this stream part
	Data     string `json:"data"`     // The data payload associated with the event
	MetaData string `json:"metadata"` // Additional metadata associated with the event
}

// Send is a structure for directing input to a specific node
type Send struct {
	Node  string `json:"node"`            // The node to send input to
	Input Json   `json:"input,omitempty"` // The input to send to the node
}

// Command represents a command to execute in the graph
type Command struct {
	Goto   any            `json:"goto,omitempty"`   // Where to go next in the graph
	Update map[string]any `json:"update,omitempty"` // Updates to apply to the state
	Resume any            `json:"resume,omitempty"` // Value to resume with
}

// The field to sort by.
type AssistantSortBy string

const (
	AssistantSortByAssistantID AssistantSortBy = "assistant_id" // The assistant ID
	AssistantSortByGraphID     AssistantSortBy = "graph_id"     // The graph ID
	AssistantSortByName        AssistantSortBy = "name"         // The name of the assistant
	AssistantSortByCreatedAt   AssistantSortBy = "created_at"   // The time the assistant was created
	AssistantSortByUpdatedAt   AssistantSortBy = "updated_at"   // The last time the assistant was updated
)

// The order to sort by.
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"  // Ascending order
	SortOrderDesc SortOrder = "desc" // Descending order
)

// The field to sort by.
type ThreadSortBy string

const (
	ThreadSortByThreadID  ThreadSortBy = "thread_id"  // The thread ID
	ThreadSortByStatus    ThreadSortBy = "status"     // The status of the thread
	ThreadSortByCreatedAt ThreadSortBy = "created_at" // The time the thread was created
	ThreadSortByUpdatedAt ThreadSortBy = "updated_at" // The last time the thread was updated
)

// RunCreateMetadata represents metadata for a run creation request
type RunCreateMetadata struct {
	RunID    string  `json:"run_id"`              // The ID of the run
	ThreadID *string `json:"thread_id,omitempty"` // The ID of the thread
}
