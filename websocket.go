package Ftx

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
	"time"
)

const (
	WebSocketEndPoint = "wss://ftx.com/ws/"
)

type SubChannel int
const (
	Sub_Fills = iota
	Sub_Orders
)

type subscribeFunc func(conn *websocket.Conn,market *string,isSubscribe bool) error
type SubscribeMap map[SubChannel]subscribeFunc

type FtxWebsocket struct {
	ErrChan chan error
	subscribMap SubscribeMap
	ctx context.Context
	cancel context.CancelFunc
}

func NewFtxWebsocket() *FtxWebsocket {
	ctx, cancel := context.WithCancel(context.Background())
	return &FtxWebsocket{
		ErrChan: make(chan error, 1),
		subscribMap: SubscribeMap{
			Sub_Fills: FillChannel,
			Sub_Orders: OrderChannel,
		},
		ctx: ctx,
		cancel: cancel,
	}
}

func(self *FtxWebsocket) str2SubChannel(channel string) SubChannel {
	switch channel {
	case "fills":
		return Sub_Fills
	case "orders":
		return Sub_Orders
	}
	return -1
}

func FillChannel(conn *websocket.Conn, market *string, isSubscribe bool) error {
	result := ""
	sub := "subscribe"
	if !isSubscribe {
		sub = "unsubscribe"
	}
	result, _ = sjson.Set(result, "op", sub)
	result, _ = sjson.Set(result, "channel", "fills")
	err := conn.WriteMessage(websocket.TextMessage, []byte(result))
	if err != nil {
		return err
	}
	return nil
}

func OrderChannel(conn *websocket.Conn, market *string, isSubscribe bool) error {
	result := ""
	sub := "subscribe"
	if !isSubscribe {
		sub = "unsubscribe"
	}
	result, _ = sjson.Set(result, "op", sub)
	result, _ = sjson.Set(result, "channel", "orders")
	err := conn.WriteMessage(websocket.TextMessage, []byte(result))
	if err != nil {
		return err
	}
	return nil
}

func(self *FtxWebsocket) Login(conn *websocket.Conn) error {
	msec := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	sign := strconv.FormatInt(msec, 10) + "websocket_login"
	sign = ComputeHmacSha256(sign, ApiSecret)
	result := "{}"
	result, _ = sjson.Set(result, "args.key", ApiKey)
	result, _ = sjson.Set(result, "args.sign", sign)
	result, _ = sjson.Set(result, "args.time", msec)
	result, _ = sjson.Set(result, "op", "login")
	err := conn.WriteMessage(websocket.TextMessage, []byte(result))
	if err != nil {
		return err
	}
	return nil
}

func(self *FtxWebsocket) PingJson() string {
	result := "{}"
	result, _ = sjson.Set(result, "op", "ping")
	return result
}

func(self *FtxWebsocket) KeepPing(conn *websocket.Conn) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	ALL:
	for {
		select {
		case <-self.ctx.Done():
			break ALL
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte(self.PingJson())); err != nil {
				self.ErrChan <- err
				break ALL
			}
		}
	}
	self.cancel()
}

func(self *FtxWebsocket) Connect(market *string, subCannel SubChannel, msgFunc func(msg []byte)) {
	conn, _, err := websocket.DefaultDialer.Dial(WebSocketEndPoint, nil)
	if err != nil {
		self.ErrChan <- err
		return
	}
	err = self.Login(conn)
	if err != nil {
		self.ErrChan <- err
		return
	}
	go self.KeepPing(conn)
	err = self.subscribMap[subCannel](conn, market, true)
	if err != nil {
		self.cancel()
		self.ErrChan <- err
		return
	}
	go func() {
		ALL:
		for {
			select {
			case <- self.ctx.Done():
				break ALL
			default:
				_, data, err := conn.ReadMessage()
				if err != nil {
					self.ErrChan <- err
					break ALL
				}
				_type := gjson.GetBytes(data, "type").String()
				switch _type {
				case "error", "info":
					self.ErrChan <- NewApiError(gjson.GetBytes(data, "code").Int(), gjson.GetBytes(data, "msg").String())
					break ALL
				case "subscribed":
					self.ErrChan <- nil
					//
				case "unsubscribed":
					self.ErrChan <- nil
					break ALL
				case "partial":
				case "update":
					msgFunc(data)
				}
			}
		}
		self.cancel()
	}()
}