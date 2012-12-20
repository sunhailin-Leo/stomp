package stomp

import (
	"net"
	"time"
)

// The STOMP server has the concept of queues and topics. A message
// sent to a queue destination will be transmitted to the next available
// client that has subscribed. A message sent to a topic will be 
// transmitted to all subscribers that are currently subscribed to the
// topic.
//
// Destinations that start with this prefix are considered to be queues.
// Destinations that do not start with this prefix are considered to be topics.
const QueuePrefix = "/queue"

// Default server parameters. 
const (
	// Default address for listening for connections.
	DefaultAddr = ":61613"

	// Default maximum bytes for the header part of STOMP frames.
	// Override by setting Server.MaxHeaderBytes.
	DefaultMaxHeaderBytes = 1 << 12 // 4KB

	// Default maximum bytes for the body part of STOMP frames.
	// Override by setting Server.MaxBodyBytes
	DefaultMaxBodyBytes = 1 << 20 // 1MB

	// Default read timeout for heart-beat.
	// Override by setting Server.HeartBeat.
	DefaultHeartBeat = time.Minute
)

// Interface for authenticating STOMP clients.
type Authenticator interface {
	// Authenticate based on the given login and passcode, either of which might be nil.
	// Returns true if authentication is successful, false otherwise.
	Authenticate(login, passcode string) bool
}

// A Server defines parameters for running a STOMP server.
type Server struct {
	Addr           string        // TCP address to listen on, DefaultAddr if empty
	Authenticator  Authenticator // Authenticates login/passcodes. If nil no authentication is performed
	QueueStorage   QueueStorage  // Implementation of queue storage. If nil, in-memory queues are used.
	HeartBeat      time.Duration // Preferred value for heart-beat read/write timeout, if zero, then DefaultHeartBeat.
	MaxHeaderBytes int           // Maximum size of STOMP headers in bytes, if zero then DefaultMaxHeaderBytes.
	MaxBodyBytes   int           // Maximum size of STOMP body in bytes, if zero then DefaultMaxBodyBytes.
}

func ListenAndServe(addr string) error {
	s := &Server{Addr: addr}
	return s.ListenAndServe()
}

func Serve(l net.Listener) error {
	s := &Server{}
	return s.Serve(l)
}

// Listens on the TCP network address s.Addr and then calls
// Serve to handle requests on the incoming connections. If
// s.Addr is blank, then DefaultAddr is used.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = DefaultAddr
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

// Accepts incoming connections on the Listener l, creating a new
// service thread for each connection. The service threads read
// requests and then process each request.
func (s *Server) Serve(l net.Listener) error {
	proc := newRequestProcessor(s)
	return proc.Serve(l)
}
