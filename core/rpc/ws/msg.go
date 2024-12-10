package ws

import "encoding/json"

type MsgType int

const (
	MsgTypeReq MsgType = iota
	MsgTypeResp
	MsgTypeNtfn
)

func (mt MsgType) String() string {
	switch mt {
	case MsgTypeReq:
		return "request"
	case MsgTypeResp:
		return "response"
	case MsgTypeNtfn:
		return "notification"
	}
	return "unknown"
}

// Message is the universal websocket message, which may encapsulate different
// types of domain-specific messages. The message and the payload are intended
// to both be JSON encoded.
type Message struct {
	ID      uint64          `json:"id"`
	Type    MsgType         `json:"type"`
	Payload json.RawMessage `json:"payload"` // e.g. marshalled jsonrpc.{Request,Response}
}
