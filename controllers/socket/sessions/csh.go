package sessions

import (
	"github.com/andreweggleston/DeathByDagger/helpers/wsevent"
	"sync"
	"time"
)

var (
	socketsMu        = new(sync.RWMutex)
	IDSockets   = make(map[string][]*wsevent.Client) //id -> client array, since players can have multiple tabs open
	connectedMu      = new(sync.Mutex)
	connectedTimer   = make(map[string]*time.Timer)
)

func AddSocket(cshusername string, so *wsevent.Client) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	IDSockets[cshusername] = append(IDSockets[cshusername], so)
	if len(IDSockets[cshusername]) == 1 {
		connectedMu.Lock()
		timer, ok := connectedTimer[cshusername]
		if ok {
			timer.Stop()
			delete(connectedTimer, cshusername)
		}
		connectedMu.Unlock()
	}
}

func RemoveSocket(sessionID, cshusername string) {
	socketsMu.Lock()
	defer socketsMu.Unlock()

	clients := IDSockets[cshusername]
	for i, socket := range clients {
		if socket.ID == sessionID {
			clients[i] = clients[len(clients)-1]
			clients[len(clients)-1] = nil
			clients = clients[:len(clients)-1]
			break
		}
	}

	IDSockets[cshusername] = clients

	if len(clients) == 0 {
		delete(IDSockets, cshusername)
	}
}

func GetSockets(cshusername string) (sockets []*wsevent.Client, success bool) {
	socketsMu.RLock()
	defer socketsMu.RUnlock()

	sockets, success = IDSockets[cshusername]
	return
}

func IsConnected(cshusername string) bool {
	_, ok := GetSockets(cshusername)
	return ok
}

func ConnectedSockets(cshusername string) int {
	socketsMu.RLock()
	l := len(IDSockets[cshusername])
	socketsMu.RUnlock()

	return l
}

func AfterDisconnectedFunc(cshusername string, d time.Duration, f func()) {
	connectedMu.Lock()
	connectedTimer[cshusername] = time.AfterFunc(d, func() {
		if !IsConnected(cshusername) {
			f()
		}

		connectedMu.Lock()
		delete(connectedTimer, cshusername)
		connectedMu.Unlock()
	})
	connectedMu.Unlock()
}