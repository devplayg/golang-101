// First: go get github.com/golang/glog
// Then : go test -v -vmodule=*=3 -logtostderr netgrace_test.go

// This test package shows how to gracefully shutdown a net.Listener.
package netgrace

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"testing"
	"time"
)

var (
	// A list of connections created by the server. We keep track of them so we
	// can gracefully shut them down if they are still alive when the server
	// goes down.
	conns []net.Conn

	// The quit channel for the server. If the server detects that this channel
	// is closed, then it's a signal for it to shutdown as well.
	quit chan bool

	// The TCP server we are going to start
	ln net.Listener
)

func init() {
	quit = make(chan bool)
	conns = make([]net.Conn, 0, 10)
}

func startServer() (err error) {
	// When the server exits, we want to clean up after ourselves.
	// Most likely all the connections have already been nil'ed but just in case.
	defer func() {
		glog.Info("Shutting down server")
		for i, conn := range conns {
			if conn != nil {
				glog.Infof("Closing connection #%d", i)
				conn.Close()
			}
		}
	}()

	// Start a echo server on port 5388/tcp
	ln, err = net.Listen("tcp", ":5388")
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		glog.Info("Listening for connections")

		// Accept() blocks waiting for new connection. The problem is that when
		// the listerner is blocked, and if there's never a new connection, then
		// this goroutine never exits (unless the program exits). We want a way
		// to gracefully stop this, so that this goroutine doesn't hang around
		// forever. For example, this is important if you are running tests with
		// a client and server, where you want the server to exit after a certain
		// period of time. If you don't somehow trigger the server to exit while
		// it's blocked on Accept(), then the test will block.
		//
		// Our requirement is to Accept() connections. If an error is returned,
		// we will try to handle it and then continue Accept(). However, we don't
		// want to continue Accept() if we are told to quit.
		//
		// One way to do this is Close() the net.Listener, which forces Accept()
		// to return with the errClosing error. Then we can check to see if the
		// error is because net.Listener is closed, or it's some other error.
		// If net.Listener is closed, then we would quit. But if it's some other
		// error, maybe we can handle it and then go back to Accept().
		//
		// Unfortunately errClosing is NOT exported from the net package so there's
		// no way to check if Accept() returned because the net.Listener is closed.
		// Others have tried to compare the actual error message ("use of closed
		// network connection"). However that is not the ideal since the actual
		// error message may change in the future.
		//
		// There was a long discussion thread,
		// https://code.google.com/p/go/issues/detail?id=4373,
		// that took place in Nov 2012 regarding the possibility of exporting
		// errClosing. But at the end the Go Authors decided it is not a
		// prudent thing to do so the issue was tagged as wontfix.
		//
		// What we need is a way to check to see if we are being told to quit,
		// after Accept() returns. The trick here is a quit channel. It's quite
		// simple actually. When you want the server to quit, first you
		// close the quit channel, which tells the Accept() goroutine to quit
		// if the goroutine checks the quit channel. Then you close the
		// net.Listener which then forces Accept() to return.
		//
		// Without further ado, we start accepting connections
		conn, err := ln.Accept()

		if err != nil {
			glog.Error(err.Error())

			// When Accept() returns with a non-nill error, we check the quit
			// channel to see if we should continue or quit. If quit, then we quit.
			// Otherwise we continue
			select {
			case <-quit:
				return nil
			default:
				// thanks to martingx on reddit for noticing I am missing a default
				// case. without the default case the select will block.
			}

			continue
		}

		// Now that we have a connection, we add it to the list, and go handle it.
		conns = append(conns, conn)
		go handleConnection(conn, len(conns)-1)
	}

	return nil

}

func handleConnection(conn net.Conn, id int) error {
	// Again, we are cleaning up after ourselves here. At the exit of this function,
	// we want to close the connection, and set the appropriate slot in the list
	// to nil so we don't leak any memory.
	defer func() {
		glog.Infof("Closing connection #%d", id)
		conn.Close()
		conns[id] = nil
	}()

	glog.Infof("Starting connection #%d", id)

	for {
		// This is effectively an echo server. Reads from the connection and
		// immediately write it back. If io.Copy() returns, then it's either
		// because the socket is closed (err == nil), or there's some type of
		// real error. Either case we return.
		if _, err := io.Copy(conn, conn); err != nil {
			glog.Error(err.Error())
			return err
		}
		return nil
	}
}

func Test10Clients(t *testing.T) {
	defer func() {
		glog.Infof("Stopping server...")

		// When we exit this test, we want to make sure we clean up after
		// ourselves so we don't leave anything behind. In this case, by
		// closing the quit channel, we are telling the server to stop
		// accepting new connection.
		close(quit)

		// We then close the net.Listener, which will force Accept() to
		// return if it's blocked waiting for new connections.
		ln.Close()

		// The above order matters somewhat. If you ln.Close() first, then 
		// you run the risk of Accept() returning but the quit channel 
		// hasn't been closed. It is not the end of the world, however,
		// it just means you will likely see quite a few more errors before
		// the goroutine detects quit is closed.
	}()

	go startServer()

	time.Sleep(time.Second)

	// In this test, we start 10 clients, each making a single connection
	// to the server. Then each will write 10 messages to the server, and
	// read the same 10 messages back. After that the client quits.
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				glog.Infof("Quiting client #%d", id)
			}()

			conn, err := net.Dial("tcp", ":5388")
			if err != nil {
				glog.Error(err.Error())
				return
			}
			defer conn.Close()

			for i := 0; i < 10; i++ {
				fmt.Fprintf(conn, "client #%d, count %d\n", id, i)
				res, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					glog.Error(err.Error())
					return
				}
				glog.Infof("Received: %s", res)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// We sleep for a couple of seconds, let the clients run their jobs,
	// then we exit, which triggers the defer function that will shutdown
	// the server.
	time.Sleep(2 * time.Second)

	// So instead of just quiting, we clean up first, well, in the defer block.
	// This is expecially important if this is a long running program.
}