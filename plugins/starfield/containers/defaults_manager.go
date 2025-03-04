package containers

import (
	"crypto/rand"
	"math/big"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var Default proxy.RegisteredServer

func GetLobbyForPlayer(player proxy.Player) proxy.RegisteredServer {
	for _, srv := range GlobalContainers.GetAllContainers() {
		if srv.Tag == "lobby" {
			if srv.Info.Players().Len() < 20 {
				Default = srv.Info
				return srv.Info
			}
		}
	}
	return Default
}

func CreateLobby() {
	Logger.Info("create debug")
	name := "lobby-"
	prefix, _ := generateRandomString(4)
	name += prefix
	CreateContainer(name, "lobby", "anton691/simple-lobby:latest")
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}
