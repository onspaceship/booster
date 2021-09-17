package socket

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/onspaceship/booster/pkg/config"
	"github.com/onspaceship/booster/pkg/socket/handler"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Options = config.SocketOptions

type socket struct {
	conn       *websocket.Conn
	refCounter *refCounter

	*Options
}

func StartListener(exit chan bool) {
	socket := New()

	log.Info("Connecting to Ground Control...")

	wait.Forever(socket.Connect, 5*time.Second)

	exit <- true
}

func New() *socket {
	options, err := config.NewSocketOptions()
	if err != nil {
		log.WithError(err).Fatal("failed to configure Ground Control socket")
	}

	return &socket{
		Options:    options,
		refCounter: newRefCounter(),
	}
}

func (socket *socket) Connect() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	socket.ensureAgentId()

	backoff := wait.Backoff{Duration: 2 * time.Second, Factor: 1.25, Jitter: 0.1, Steps: math.MaxInt32}
	err := wait.ExponentialBackoff(backoff, func() (done bool, err error) {
		conn, resp, err := websocket.DefaultDialer.DialContext(ctx, socket.Host, http.Header{
			"X-Agent-ID": {socket.AgentId},
		})

		if err != nil {
			logline := log.WithError(err)

			if resp != nil {
				logline = logline.WithField("status", resp.Status)
			}

			logline.Error("Could not reach Ground Control")
			return false, nil
		}

		log.Info("Connected to Ground Control!")

		socket.conn = conn
		return true, nil
	})

	if err != nil {
		log.Fatal("Retry attempts exceeded when connecting to Ground Control")
	}

	defer socket.conn.Close()

	go socket.listen(cancel)

	socket.joinTopic("booster")
	socket.joinTopic("booster:" + socket.AgentId)

	go socket.heartbeat(cancel)

	<-ctx.Done()
}

func (socket *socket) listen(done context.CancelFunc) {
	defer done()

	var message Message

	for {
		err := socket.conn.ReadJSON(&message)
		if err != nil {
			log.WithError(err).Error("Error reading from Ground Control")
			return
		}

		if !isPhoenixEvent(message.Event) && message.Topic == "booster:"+socket.AgentId {
			log.WithField("event", message.Event).WithField("payload", message.Payload).Debug("New message from Ground Control")

			handler.Handle(string(message.Event), message.Payload, socket.Options)
		}
	}

}

func (socket *socket) heartbeat(done context.CancelFunc) {
	defer done()

	for range time.Tick(30 * time.Second) {
		socket.sendHeartbeat()
	}
}
