package connect

import "sync"
import "time"
import clientConnect "github.com/LilyPad/GoLilyPad/client/connect"
import packetConnect "github.com/LilyPad/GoLilyPad/packet/connect"

type FluxConnect struct {
	Client       clientConnect.Connect
	Servers      map[string]*Server
	ServersMutex *sync.RWMutex
	Players      uint16
	MaxPlayers   uint16
}

func NewFluxConnect(addr *string, user *string, pass *string, done chan bool) (connect *FluxConnect) {
	client := clientConnect.NewConnect()
	connect = &FluxConnect{
		Client:       client,
		Servers:      make(map[string]*Server),
		ServersMutex: &sync.RWMutex{},
	}

	// Initialize the connect object on before we connect
	client.RegisterEvent("preconnect", func(event clientConnect.Event) {
		if len(connect.Servers) > 0 {
			connect.ServersMutex.Lock()
			connect.Servers = make(map[string]*Server)
			connect.ServersMutex.Unlock()
		}
		connect.Players = 0
		connect.MaxPlayers = 0
	})

	clientConnect.AutoAuthenticate(client, user, pass)
	go clientConnect.AutoConnect(client, addr, done)

	// Start collecing the player count of the cluster ever second
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				if !client.Connected() {
					break
				}
				connect.QueryRemotePlayers()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return connect
}

func (this *FluxConnect) QueryRemotePlayers() {
	this.Client.RequestLater(&packetConnect.RequestGetPlayers{false}, func(statusCode uint8, result packetConnect.Result) {

		if result == nil {
			return
		}

		getPlayersResult := result.(*packetConnect.ResultGetPlayers)
		this.Players = getPlayersResult.CurrentPlayers
		this.MaxPlayers = getPlayersResult.MaximumPlayers
	})
}
