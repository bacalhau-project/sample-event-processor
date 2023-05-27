package app

type Event struct {
	Source string
	Offset int64
	Data   string
}

type Checkpointer interface {
	Restore() (int64, error)
	Save(offset int64) error
}

type EventProcessor interface {
	ProcessEvent(event Event) error
}
