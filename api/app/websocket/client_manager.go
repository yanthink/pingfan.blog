package websocket

import (
	"blog/app"
	"go.uber.org/zap"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients           map[*Client]bool      // 全部客户端连接
	ClientsLock       sync.RWMutex          // 全部客户端连接读写锁
	Users             map[uint64][]*Client  // 已登录的用户客户端连接
	UsersLock         sync.RWMutex          // 已登录的用户客户端连接读写锁
	TempUsers         map[string]*Client    // 临时用户客户端链接
	TempUsersLock     sync.RWMutex          // 临时用户客户端读写锁
	Register          chan *Client          // 连接处理
	Login             chan *Client          // 用户登录处理
	Unregister        chan *Client          // 断开连接处理程序
	UserBroadcast     chan *UserMessage     // 给用户广播消息
	TempUserBroadcast chan *TempUserMessage // 给临时用户广播消息
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:           map[*Client]bool{},
		Users:             map[uint64][]*Client{},
		TempUsers:         map[string]*Client{},
		Register:          make(chan *Client, 1000),
		Login:             make(chan *Client, 1000),
		Unregister:        make(chan *Client, 1000),
		UserBroadcast:     make(chan *UserMessage, 1000),
		TempUserBroadcast: make(chan *TempUserMessage, 1000),
	}
	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	clientManager.ClientsLock.RLock()
	defer clientManager.ClientsLock.RUnlock()

	now := uint64(time.Now().Unix())
	for client := range clientManager.Clients {
		if client.HeartbeatTime+300 < now {
			app.Logger.Debug(
				"心跳时间超时关闭连接",
				zap.String("Addr", client.Addr),
				zap.Uint64("UserID", client.UserID),
				zap.String("TempUserID", client.TempUserID),
				zap.Time("LoginTime", time.Unix(int64(client.LoginTime), 0)),
				zap.Time("HeartbeatTime", time.Unix(int64(client.HeartbeatTime), 0)),
			)
			client.Socket.Close()
		}
	}
}

// SendToUser 给用户发送消息
func SendToUser(message *UserMessage) {
	clientManager.UserBroadcast <- message
}

// SendToTempUser 给临时用户广播消息
func SendToTempUser(message *TempUserMessage) {
	clientManager.TempUserBroadcast <- message
}

func IsOnline(id uint64) (ok bool) {
	_, ok = clientManager.Users[id]

	return
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
}

func (manager *ClientManager) DelClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

func (manager *ClientManager) AddUser(client *Client) {
	if client.UserID <= 0 {
		return
	}

	manager.UsersLock.Lock()
	defer manager.UsersLock.Unlock()

	for _, c := range manager.Users[client.UserID] {
		if c == client {
			return
		}
	}

	manager.Users[client.UserID] = append(manager.Users[client.UserID], client)
}

func (manager *ClientManager) DelUser(client *Client) {
	manager.UsersLock.Lock()
	defer manager.UsersLock.Unlock()

	if _, ok := manager.Users[client.UserID]; !ok {
		return
	}

	for i, c := range manager.Users[client.UserID] {
		if c == client {
			manager.Users[client.UserID] = append(manager.Users[client.UserID][:i], manager.Users[client.UserID][i+1:]...)
			if len(manager.Users[client.UserID]) == 0 {
				delete(manager.Users, client.UserID)
			}
		}
	}
}

func (manager *ClientManager) AddTempUser(client *Client) {
	if client.TempUserID == "" {
		return
	}

	manager.TempUsersLock.Lock()
	defer manager.TempUsersLock.Unlock()

	manager.TempUsers[client.TempUserID] = client
}

func (manager *ClientManager) DelTempUser(client *Client) {
	if client.TempUserID == "" {
		return
	}

	manager.TempUsersLock.Lock()
	defer manager.TempUsersLock.Unlock()

	delete(manager.TempUsers, client.TempUserID)

	client.TempUserID = ""
}

func (manager *ClientManager) RegisterHandle(client *Client) {
	manager.AddClient(client)

	app.Logger.Debug("用户连接")
}

func (manager *ClientManager) LoginHandle(client *Client) {
	if client.UserID > 0 {
		manager.AddUser(client)
		manager.DelTempUser(client)
	}

	if client.TempUserID != "" {
		manager.AddTempUser(client)
	}

	app.Logger.Debug("用户登录", zap.Uint64("UserID", client.UserID), zap.String("TempUserID", client.TempUserID))
}

func (manager *ClientManager) UnregisterHandle(client *Client) {
	manager.DelClient(client)
	manager.DelUser(client)
	manager.DelTempUser(client)

	app.Logger.Debug("用户断开连接", zap.Uint64("UserID", client.UserID), zap.String("TempUserID", client.TempUserID))
}

func (manager *ClientManager) UserBroadcastHandle(message *UserMessage) {
	manager.UsersLock.RLock()
	defer manager.UsersLock.RUnlock()

	for _, client := range manager.Users[message.UserID] {
		client.SendMsg(message.Response)
	}
}

func (manager *ClientManager) TempUserBroadcastHandle(message *TempUserMessage) {
	manager.TempUsersLock.RLock()
	defer manager.TempUsersLock.RUnlock()

	if client, ok := manager.TempUsers[message.TempUserID]; ok {
		client.SendMsg(message.Response)
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.RegisterHandle(conn)
		case conn := <-manager.Login:
			manager.LoginHandle(conn)
		case conn := <-manager.Unregister:
			manager.UnregisterHandle(conn)
		case message := <-manager.UserBroadcast:
			manager.UserBroadcastHandle(message)
		case message := <-manager.TempUserBroadcast:
			manager.TempUserBroadcastHandle(message)
		}
	}
}
