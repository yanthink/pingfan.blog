package bootstrap

import "blog/app/websocket"

func SetupWebsocket() {
	websocket.Start()
}
