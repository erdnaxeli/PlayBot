package irc

type (
	Event   int
	Handler func(*Conn, Message) error
)

const (
	// NoEvent is a dummy value that corresponds to the the zero value of the Event type.
	NoEvent Event = iota
	Mode
	Notice
	RPLWelcome
	PrivMsg
)

// AddHandler registers the given handler for the given event. If a previous handler
// for this event was already registered, the previous one is discarded.
//
// An handler accept two parameters: the IRC connection and the message. Any error
// returned by an handler stops the dispatching (see [Dispatch()]).
func (i *Conn) AddHandler(event Event, handler Handler) {
	i.handlers[event] = handler
}

// OnConnect sets an handler to run on connection.
func (i *Conn) OnConnect(handler Handler) {
	i.AddHandler(RPLWelcome, handler)
}

// OnMode sets an handler to run on MODE event.
func (i *Conn) OnMode(handler Handler) {
	i.AddHandler(Mode, handler)
}

// OnNotice sets an handler to run on NOTICE event.
func (i *Conn) OnNotice(handler Handler) {
	i.AddHandler(Notice, handler)
}

// OnPrivMsg sets an handler to run on each message.
func (i *Conn) OnPrivMsg(handler Handler) {
	i.AddHandler(PrivMsg, handler)
}
