package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/swedishborgie/treadonme"
	"github.com/urfave/cli/v2"
)

//go:embed static/*
var staticFS embed.FS

type webserver struct {
	bindAddr       string
	macAddress     string
	connectTimeout time.Duration
	tmClient       *treadonme.Treadmill
	tmMutex        sync.Mutex
	devInfo        *treadonme.MessageDeviceInfo

	wsClients []*websocket.Conn
	wsMutex   sync.Mutex
}

type ClientMessage struct {
	Command string
}

type MessageWrapper struct {
	Error   string
	Type    string
	Message treadonme.Message
}

func main() {
	app := &cli.App{
		Name:        "treadonme",
		Description: "a small web server for getting live telemetry from sole treadmills",
		Action:      run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bind-address",
				Usage:   "the socket address to bind to",
				EnvVars: []string{"TREAD_BIND_ADDRESS"},
				Value:   ":8089",
			},
			&cli.StringFlag{
				Name:     "mac-address",
				Usage:    "the mac address of the treadmill",
				EnvVars:  []string{"TREAD_MAC_ADDRESS"},
				Required: true,
			},
			&cli.DurationFlag{
				Name:    "connect-timeout",
				Usage:   "the amount of time to wait before timing out on connect",
				EnvVars: []string{"TREAD_CONNECT_TIMEOUT"},
				Value:   60 * time.Second,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cliCtx *cli.Context) error {
	ws := &webserver{
		bindAddr:       cliCtx.String("bind-address"),
		macAddress:     cliCtx.String("mac-address"),
		connectTimeout: cliCtx.Duration("connect-timeout"),
	}

	log.Printf("starting server listening on %s looking for treadill at %s", ws.bindAddr, ws.macAddress)

	if err := ws.start(); err != nil {
		return err
	}

	return nil
}

func (ws *webserver) start() error {
	subDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(subDir)))
	http.HandleFunc("/ws", ws.wsEndpoint)

	if err := http.ListenAndServe(ws.bindAddr, nil); err != nil {
		return err
	}

	return nil
}

func (ws *webserver) startTreadmill() error {
	ws.tmMutex.Lock()
	defer ws.tmMutex.Unlock()

	tm, err := treadonme.New(ws.macAddress)
	if err != nil {
		return err
	}

	if err := tm.Connect(context.Background()); err != nil {
		return err
	}

	tm.AddListener(ws.treadmillListener)

	devInfo, err := tm.GetDeviceInfo()
	if err != nil {
		return err
	}

	ws.devInfo = devInfo

	if _, err := tm.WaitForResponse(context.Background(), treadonme.MessageTypeHeartRateType); err != nil {
		return err
	}

	if err := tm.Start(); err != nil {
		return err
	}

	ws.tmClient = tm

	return nil
}

func (ws *webserver) stopTreadmill() {
	ws.tmMutex.Lock()
	defer ws.tmMutex.Unlock()

	if err := ws.tmClient.Close(); err != nil {
		log.Printf("problem closing treadmill after workout: %s", err)
	}

	ws.tmClient = nil
	ws.devInfo = nil
}

func (ws *webserver) treadmillListener(msg treadonme.Message, err error) {
	if err != nil {
		ws.notifyClients(&MessageWrapper{Error: err.Error()})

		return
	}

	ws.notifyClients(&MessageWrapper{Type: msg.MessageType().String(), Message: msg})

	if msg.MessageType() == treadonme.MessageTypeEndWorkout {
		go ws.stopTreadmill()
	}
}

func (ws *webserver) addClient(client *websocket.Conn) {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()
	ws.wsClients = append(ws.wsClients, client)
}

func (ws *webserver) removeClient(client *websocket.Conn) {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()

	scrubbed := make([]*websocket.Conn, 0, len(ws.wsClients)-1)
	for _, c := range ws.wsClients {
		if c != client {
			scrubbed = append(scrubbed, c)
		}
	}
	ws.wsClients = scrubbed
}

func (ws *webserver) notifyClients(msg *MessageWrapper) {
	ws.wsMutex.Lock()
	defer ws.wsMutex.Unlock()

	for _, c := range ws.wsClients {
		if err := c.WriteJSON(msg); err != nil {
			log.Printf("problem writing message to websocket client: %s", err)
		}
	}
}

func (ws *webserver) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	c, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
	if err != nil {
		log.Printf("unable to upgrade: %s", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("problem closing websocket client: %s", err)
		}
	}()

	ws.addClient(c)
	defer ws.removeClient(c)

	if ws.devInfo != nil {
		if err := c.WriteJSON(&MessageWrapper{Type: treadonme.MessageTypeDeviceInfo.String(), Message: ws.devInfo}); err != nil {
			log.Printf("problem writing initial dev info to client: %s", err)

			return
		}
	}

	for {
		cm := &ClientMessage{}

		if err := c.ReadJSON(cm); err != nil {
			log.Printf("unable to read message from websocket client: %s", err)

			return
		}

		switch cm.Command {
		case "start":
			if err := ws.startTreadmill(); err != nil {
				log.Printf("problem starting treadmill: %s", err)

				if writeErr := c.WriteJSON(&MessageWrapper{Error: err.Error()}); writeErr != nil {
					log.Printf("problem writing error message to client: %s", writeErr)
				}
			}
		}
	}
}
