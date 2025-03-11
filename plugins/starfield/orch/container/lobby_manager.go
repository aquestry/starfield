package container

import (
	"crypto/rand"
	"github.com/aquestry/starfield/plugins/starfield/logger"
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
	for _, c := range GetContainers() {
		if c.Tag == "lobby" && c.Info.Players().Len() < defaultMax {
			servers = append(servers, c.Info)
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
	for _, s := range servers {
		count := s.Players().Len()
		if count < lowestCount {
			lowest = s
			lowestCount = count
		}
	}
	return lowest
}

func isOtherServerUnderMin(other proxy.RegisteredServer) bool {
	for _, s := range getAvailableServers() {
		if s.Players().Len() < defaultMin && s != other {
			return true
		}
	}
	return false
}

func getServerBetweenMinAndMaxPlayers() proxy.RegisteredServer {
	for _, s := range getAvailableServers() {
		count := s.Players().Len()
		if count >= defaultMin && count < defaultMax {
			return s
		}
	}
	return nil
}

func CreateLobby() {
	name := "Lobby-" + generateRandomString(4)
	_, err := CreateContainer(name, "lobby", defaultServerTemplate)
	if err != nil {
		logger.L.Info("create", "error", err.Error())
	}
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
