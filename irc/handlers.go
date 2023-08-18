package irc

type Event int
type Handler func(*Conn, Message) error

const (
	// Dummy value that corresponds to the the zero value of the Event type.
	NO_EVENT Event = iota
	NOTICE
	RPL_WELCOME
	PRIVMSG
)

// AddHandler registers the given handler for the given event. If a previous handler
// for this event was already registered, the previous one is discarded.
//
// An handler accept two parameters: the IRC connection and the message. Any error
// returned by an handler stops the dispatching (see [Dispatch()]).
func (i *Conn) AddHandler(event Event, handler Handler) {
	i.handlers[event] = handler
}

func (i *Conn) OnConnect(handler Handler) {
	i.AddHandler(RPL_WELCOME, handler)
}

func (i *Conn) OnNotice(handler Handler) {
	i.AddHandler(NOTICE, handler)
}

func (i *Conn) OnPrivmsg(handler Handler) {
	i.AddHandler(PRIVMSG, handler)
}
