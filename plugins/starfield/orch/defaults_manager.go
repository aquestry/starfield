package orch

import (
	"crypto/rand"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"math/big"
	"sync"
)

var (
	defaultMin            = 1
	defaultMax            = 2
	defaultServerTemplate = "anton691/simple-lobby:latest"
	mutex                 sync.Mutex
	defaultLobby          proxy.RegisteredServer
)

func GetTargetLobby() proxy.RegisteredServer {
	mutex.Lock()
	defer mutex.Unlock()

	target := getServerWithLowestPlayerCount()
	if target != nil {
		count := target.Players().Len()
		if count+1 == defaultMin && !isOtherServerUnderMin(target) {
			CreateLobby()
		}
	}
	target = getServerBetweenMinAndMaxPlayers()
	if target != nil {
		return target
	}
	return getServerWithLowestPlayerCount()
}

func getAvailableServers() []proxy.RegisteredServer {
	var servers []proxy.RegisteredServer
	for _, srv := range GetContainers() {
		if srv.Tag == "lobby" && srv.Info.Players().Len() < defaultMax {
			servers = append(servers, srv.Info)
		}
	}
	return servers
}

func getServerWithLowestPlayerCount() proxy.RegisteredServer {
	servers := getAvailableServers()
	if len(servers) == 0 {
		return defaultLobby
	}
	lowest := servers[0]
	lowestCount := lowest.Players().Len()
	for _, srv := range servers {
		count := srv.Players().Len()
		if count < lowestCount {
			lowest = srv
			lowestCount = count
		}
	}
	return lowest
}

func isOtherServerUnderMin(other proxy.RegisteredServer) bool {
	for _, srv := range getAvailableServers() {
		if srv.Players().Len() < defaultMin && srv != other {
			return true
		}
	}
	return false
}

func getServerBetweenMinAndMaxPlayers() proxy.RegisteredServer {
	for _, srv := range getAvailableServers() {
		count := srv.Players().Len()
		if count >= defaultMin && count < defaultMax {
			return srv
		}
	}
	return nil
}

func CreateLobby() {
	name := "Lobby-" + generateRandomString(4)
	CreateContainer(name, "lobby", defaultServerTemplate)
}

func generateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}
