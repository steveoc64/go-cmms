package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"sync"
	"time"

	"itrak-cmms/shared"

	"golang.org/x/net/websocket"
)

// Wrapper for socket, which has a controlling mutex to be shared with the RPC server
// and session data once the user logs in on this connection
type Connection struct {
	ID       int
	Active   bool
	Socket   *websocket.Conn
	Mutex    *sync.Mutex
	Username string
	UserID   int
	UserRole string
	Time     time.Time
	ticker   *time.Ticker
	enc      *gob.Encoder
	r        rpc.Response
	Route    string
	Routes   []string
}

// Safely send unsolicited RPC response to a connection
func (c *Connection) Send(name string, payload interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.r.ServiceMethod = name
	c.r.Seq = 0

	if err := c.enc.Encode(&c.r); err != nil {
		log.Println("Header", name, err.Error())
		return err
	}
	if err := c.enc.Encode(payload); err != nil {
		log.Println("Payload", payload, err.Error())
		return err
	}
	// log.Println("got here with", payload)
	return nil
}

// Upgrade the session data for this connection
func (c *Connection) Login(username string, id int, role string) {
	c.Username = username
	c.UserID = id
	c.UserRole = role
	c.Route = ""
	c.Time = time.Now()
}

// Constantly Ping the Backend
func (c *Connection) KeepAlive(sec time.Duration) {
	// log.Println("sending ping to ", c.ID)

	data := shared.AsyncMessage{
		Action: "Ping",
		ID:     c.ID,
	}
	c.Send("Ping", data)
	c.ticker = time.NewTicker(time.Second * sec)
	for range c.ticker.C {
		// log.Println("sending ping to client", c.ID)
		err := c.Send("Ping", data)
		if err != nil {
			log.Println("Send error on", c.ID, err.Error())
		}
	}
}

// Send an async message to everyone but this connection
func (c *Connection) Broadcast(name string, action string, id int) {

	data := shared.AsyncMessage{
		Action: action,
		ID:     id,
	}

	for _, v := range Connections.cmap {
		if v != c && v.UserID != 0 {
			log.Println("broadcast", name, action, id, "»", v.ID)
			go v.Send(name, data)
		}
	}
}

// Send an async message to everyone but this connection, if they are admin
func (c *Connection) BroadcastAdmin(name string, action string, id int) {

	data := shared.AsyncMessage{
		Action: action,
		ID:     id,
	}

	for _, v := range Connections.cmap {
		if v != c && v.UserID != 0 && v.UserRole == "Admin" {
			log.Println("broadcastAdmin", name, action, id, "»", v.ID)
			go v.Send(name, data)
		}
	}
}

// A collection of Connections
type ConnectionsList struct {
	// conns  []*Connection
	cmap map[int]*Connection
	keys []int

	nextID int
}

func (c *ConnectionsList) Map() map[int]*Connection {
	return c.cmap
}

func (c *ConnectionsList) Keys() []int {
	return c.keys
}

// Send an async message to everyone that is connected
func (c *ConnectionsList) BroadcastAll(name string, action string, id int) {

	data := shared.AsyncMessage{
		Action: action,
		ID:     id,
	}

	for _, v := range c.cmap {
		if v.UserID != 0 {
			log.Println("BroadcastAll", name, action, id, "»", v.ID)
			go v.Send(name, data)
		}
	}
}

// Send an async message to all admints that are connected
func (c *ConnectionsList) BroadcastAllAdmin(name string, action string, id int) {

	data := shared.AsyncMessage{
		Action: action,
		ID:     id,
	}

	for _, v := range c.cmap {
		if v.UserID != 0 && v.UserRole == "Admin" {
			log.Println("BroadcastAllAdmin", name, action, id, "»", v.ID)
			go v.Send(name, data)
		}
	}
}

var Connections *ConnectionsList

// Find the connection that owns the socket, return nil if not found
func (c *ConnectionsList) Find(ws *websocket.Conn) *Connection {
	for _, conn := range c.cmap {
		if conn.Socket == ws {
			return conn
		}
	}
	return nil
}

// Get the connection by ID
func (c *ConnectionsList) Get(id int) *Connection {
	return c.cmap[id]
}

// Add a websocket to the list, creates a matching Mutex, and returns the meta-Connection
func (c *ConnectionsList) Add(ws *websocket.Conn) *Connection {
	conn := &Connection{
		ID:     c.nextID + 1,
		Socket: ws,
		Mutex:  new(sync.Mutex),
		enc:    gob.NewEncoder(ws),
	}
	// c.conns = append(c.conns, conn)
	c.nextID++
	if c.cmap == nil {
		c.cmap = make(map[int]*Connection)
	}
	c.cmap[c.nextID] = conn
	c.keys = append(c.keys, c.nextID)

	// Now create a keepalive pinger for this connection
	go conn.KeepAlive(55)

	return conn
}

// Remove the websocket from the list by ID
func (c *ConnectionsList) Drop(conn *Connection) *ConnectionsList {
	fmt.Println("Remove connection ", conn.ID)

	// theConn := c.cmap[conn.ID]
	// println("theConn = ", theConn)
	delete(c.cmap, conn.ID)

	// must remove the offending key from the keys array now
	for i := 0; i < len(c.keys); i++ {
		if c.keys[i] == conn.ID {
			c.keys = append(c.keys[:i], c.keys[i+1:]...) // NOTE - variadic tail, append takes varargs of type, not an array
		}
	}

	c.BroadcastAll("login", "delete", conn.ID)
	return c
}

// Show all the active websocket connections
func (c *ConnectionsList) Show(header string) *ConnectionsList {
	fmt.Println("==================================")
	fmt.Println(header)
	for _, key := range c.keys {
		conn := c.cmap[key]
		req := conn.Socket.Request()
		theIP := ""
		if theIP = req.Header.Get("X-Real-Ip"); theIP == "" {
			theIP = req.RemoteAddr
		}
		fmt.Printf("  %d:%s\t\t%s\n", conn.ID, theIP, req.Header["User-Agent"])

		if conn.UserID != 0 {
			fmt.Println("\t\t\t",
				"User:", conn.Username, conn.UserID,
				"Route:", conn.Route,
				"Time:", time.Since(conn.Time))
		}
	}
	fmt.Println("==================================")
	return c
}

func webSocket(ws *websocket.Conn) {

	// ws := c.Socket()
	ws.PayloadType = websocket.BinaryFrame

	conn := Connections.Add(ws)
	Connections.Show("Connections Grows To:")

	// Create a custom RPC server for this socket
	buf := bufio.NewWriter(ws)
	srv := &myServerCodec{
		rwc:    ws,
		conn:   conn,
		dec:    gob.NewDecoder(ws),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	rpc.ServeCodec(srv)
}

// gobbing RPC Codec with a Mutex to allow sharing of the line with other senders
type myServerCodec struct {
	rwc    io.ReadWriteCloser
	conn   *Connection
	mutex  *sync.Mutex
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool
}

// On receiving a new header, lock the connection until the whole RPC call has finished
func (c *myServerCodec) ReadRequestHeader(r *rpc.Request) error {
	err := c.dec.Decode(r)
	if err != nil {
		log.Println("Dropped Connection:", err.Error(), ", connection:", c.conn.ID)
		Connections.Drop(c.conn)
	}
	c.conn.Mutex.Lock()
	return err
}

func (c *myServerCodec) ReadRequestBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *myServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	// as soon as we are done, unlock the connection Mutex
	defer c.conn.Mutex.Unlock()

	if err = c.enc.Encode(r); err != nil {
		if c.encBuf.Flush() == nil {
			// Gob couldn't encode the header. Should not happen, so if it does,
			// shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = c.enc.Encode(body); err != nil {
		if c.encBuf.Flush() == nil {
			// Was a gob problem encoding the body but the header has been written.
			// Shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *myServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}
