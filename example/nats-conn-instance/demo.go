package main

import (
	nats "github.com/nats-io/nats.go"
	"github.com/xswitch-cn/xctrl/ctrl"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var natsUrl = os.Getenv("NATS_URL")

func init() {
	if natsUrl == "" {
		natsUrl = "nats://localhost:4222"
	}
}

func main() {
	shutdown := make(chan os.Signal, 1)

	instance, err := ctrl.NewCtrlInstance(true, natsUrl)
	if err != nil {
		log.Panic(err)
	}

	conn := instance.GetNATSConn()
	if conn == nil {
		log.Panic("nats conn is nil.")
	}
	conn.SetClosedHandler(func(conn *nats.Conn) {
		log.Printf("nats connection: %s has been closed", conn.ConnectedAddr())
	})

	conn.SetDisconnectErrHandler(func(conn *nats.Conn, err error) {
		log.Printf("nats connection: %s has been closed", conn.ConnectedAddr())
	})

	conn.SetReconnectHandler(func(conn *nats.Conn) {
		log.Printf("Reconnect to nats: %s", conn.ConnectedAddr())
	})

	// MaxReconnect sets the number of reconnect attempts that will be tried before giving up. If negative, then it will never give up trying to reconnect. Default is 60
	//conn.Opts.MaxReconnect = -1
	// ReconnectWait sets the time to backoff after attempting to (and failing to) reconnect. Default is 2 * time.Second
	//conn.Opts.ReconnectWait = time.Second

	go func() {
		for {
			if conn.IsConnected() {
				log.Println("nats is connected")
			} else {
				log.Println("nats isn't connected")
			}
			time.Sleep(time.Second)
		}
	}()

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTSTP)

	<-shutdown
	log.Println("shutting down")
	os.Exit(0)

}
