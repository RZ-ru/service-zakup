package request // Описывает возможные статусы и переходы.

type Status string

const (
	StatusDraft      Status = "draft"
	StatusSubmitted  Status = "submitted"
	StatusApproved   Status = "approved"
	StatusRejected   Status = "rejected"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

var allowedTransitions = map[Status][]Status{
	StatusDraft: {
		StatusSubmitted,
		StatusCanceled,
	},
	StatusSubmitted: {
		StatusApproved,
		StatusRejected,
		StatusCanceled,
	},
	StatusApproved: {
		StatusInProgress,
		StatusCanceled,
	},
	StatusInProgress: {
		StatusCompleted,
		StatusCanceled,
	},
	StatusRejected:  {},
	StatusCompleted: {},
	StatusCanceled:  {},
}

func (s Status) Valid() bool {
	switch s {
	case StatusDraft,
		StatusSubmitted,
		StatusApproved,
		StatusRejected,
		StatusInProgress,
		StatusCompleted,
		StatusCanceled:
		return true
	default:
		return false
	}
}
