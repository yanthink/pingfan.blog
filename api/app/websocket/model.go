package websocket

type Request struct {
	Event string         `json:"event"` // 事件名称
	Data  map[string]any `json:"data"`  // 数据
}

type Response struct {
	Event string `json:"event"` // 事件名称
	Data  any    `json:"data"`  // 数据
}

type UserMessage struct {
	*Response
	UserID uint64
}

type TempUserMessage struct {
	*Response
	TempUserID string
}
