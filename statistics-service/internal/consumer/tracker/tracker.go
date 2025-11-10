package tracker

// ProcessedEventTracker 已处理事件跟踪器
type ProcessedEventTracker interface {
	IsProcessed(eventID string) (bool, error)
	MarkProcessed(eventID string) error
}

type DefaultTracker struct {
}

func NewDefaultTracker() *DefaultTracker {
	return &DefaultTracker{}
}

func (t *DefaultTracker) IsProcessed(eventID string) (bool, error) {
	return false, nil
}

func (t *DefaultTracker) MarkProcessed(eventID string) error {
	return nil
}
