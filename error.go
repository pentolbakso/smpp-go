package smpp

var SessionClosedBeforeReceiving error = sessionClosedBeforeReceiving{}

type sessionClosedBeforeReceiving struct{}

func (s sessionClosedBeforeReceiving) Error() string {
	return "smpp: session closed before receiving response"
}
