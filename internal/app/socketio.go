package app

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"

	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

// TODO: make something useful to the template
// This function adds Socket.IO routes to the Echo server.
func addSocketIoRoutes(
	e *echo.Echo,
	injector *do.Injector,
) {
	log.DEBUG = true // Enable debug logging for the engine.io library.
	c := socket.DefaultServerOptions()
	// Serves the client-side /socket.io/socket.io.js library over HTTP — useful during development.
	c.SetServeClient(true)

	// Overrides the default heartbeat values:
	// •	Server sends a ping every 300ms.
	// •	If no pong is received within 200ms, the connection is considered lost.
	// Normally defaults are much higher (e.g., 25 s and 20 s)
	c.SetPingInterval(300 * time.Millisecond)
	c.SetPingTimeout(200 * time.Millisecond)

	// Limits HTTP upload size to ~1 MB and times out slow handshakes after 1 second.

	c.SetMaxHttpBufferSize(1000000)
	c.SetConnectTimeout(1000 * time.Millisecond)

	// Enables CORS so any origin is allowed, and credentials (cookies, auth headers) can be sent.

	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})

	//Creates a new Socket.IO server instance using the default HTTP handler options.

	socketio := socket.NewServer(nil, nil)

	// Listens for new socket connections on the default namespace ("/").
	// The client object represents the connected socket.
	socketio.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)

		// When client emits message, the server echoes it back on message-back.
		client.On("message", func(args ...interface{}) {
			client.Emit("message-back", args...)
		})

		//Sends an auth event with authentication data from the handshake.
		client.Emit("auth", client.Handshake().Auth)

		//Sends an auth event with authentication data from the handshake.
		client.Emit("message", "Hello from server! please send me a `message` event, and I will echo it back to you")

		// Handles a message with an acknowledgment callback: server sends arguments back via ack(...).
		client.On("message-with-ack", func(args ...interface{}) {
			ack, ok := args[len(args)-1].(socket.Ack)

			if !ok {
				client.Emit("trouble", "last argument is not an ack function")
				return
			}

			ack(args[:len(args)-1], nil)
		})

	})

	// Creates a custom namespace (/custom) with its own connection handler; sends same auth event.
	socketio.Of("/custom", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.Emit("auth", client.Handshake().Auth)
	})

	// Add the Socket.IO server as a handler for the Echo framework.
	// This allows the Echo server to handle WebSocket and HTTP requests for Socket.IO.
	e.Any("/socket.io/*", echo.WrapHandler(socketio.ServeHandler(c)))
}
