package events

type TaskEvent struct {
	TaskID    string `json:"task_id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Status    string `json:"status,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
	EventType string `json:"event_type,omitempty"`
}
