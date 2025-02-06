package bagman

// Object describes literally all of the objects that we will use in our proxy
type Object interface {
	Import(any) error
	Export() (any, error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	Kind() string
}

// Contented describes objects that have content
type Contented interface {
	Content() (any, error)
}

// Partly describes objects that have parts
type PartIndexer interface {
	PartsOf(string) []Part
	AllParts() []Part
}

type PartAdder interface {
	AddPart(Part) error
}

// Updateable describes objects that can be updated
type Updateable interface {
	Update(any) error
}

// Part describes objects that are parts of other objects
type Part interface {
	Object
	Contented
	Updateable
}

type Message interface {
	Object
	Contented
	PartIndexer
}

type Request interface {
	Object
	MessagesFor(string) []Message
}

type Logger interface {
	Debug(string, ...any)
	Error(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
}
