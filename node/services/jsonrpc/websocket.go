package rpcserver

// TODO: move to node/services/ws

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"github.com/kwilteam/kwil-db/core/log"
	"github.com/kwilteam/kwil-db/core/rpc/ws"
)

const (
	wsReadTimeout      = time.Second
	wsWriteTimeout     = 5 * time.Second
	wsDefaultReadLimit = 100_000
)

type wsSession struct {
	// id     uint64
	conn   *websocket.Conn
	cancel context.CancelFunc
	logger log.Logger

	// reqMtx sync.RWMutex
	// reqs   map[uint64]struct{}

	outChan chan *ws.Message
}

// wsSession is generalized, hands off the WsMessage to handler func.

func (ses *wsSession) pinger(ctx context.Context) {
	defer ses.cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ping := []byte{}
out:
	for {
		select {
		case <-ticker.C:
			err := ses.conn.WriteControl(websocket.PingMessage, ping, time.Now().Add(wsWriteTimeout))
			if err != nil {
				ses.logger.Debugf("WriteControl(ping) failed: %v", err)
				break out
			}
		case <-ctx.Done():
			break out
		}
	}
}

func (ses *wsSession) reader(ctx context.Context,
	handleMessage func(ctx context.Context, msg *ws.Message, resp chan<- *ws.Message)) {
	defer ses.cancel()

	for {
		if ctx.Err() != nil { // gorilla...
			return
		}
		messageType, rd, err := ses.conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				ses.logger.Infof("conn NextReader failed %v", err)
			}
			return
		}
		if messageType != websocket.TextMessage {
			ses.logger.Warnf("got unexpected message type %v", messageType)
			continue
		}

		var msg ws.Message
		err = json.NewDecoder(rd).Decode(&msg)
		if err != nil {
			ses.logger.Warnf("failed to read incoming message: %v", err)
			continue
		}

		go handleMessage(ctx, &msg, ses.outChan)
	}
}

func (ses *wsSession) writer(ctx context.Context) {
	defer ses.cancel()

	for {
		var msg *ws.Message
		select {
		case msg = <-ses.outChan:
		case <-ctx.Done():
			return
		}

		ses.conn.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
		err := ses.conn.WriteJSON(msg)
		if err != nil {
			ses.logger.Infof("write failed: %v", err)
			return
		}
	}
}
