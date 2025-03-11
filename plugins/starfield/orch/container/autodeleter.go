package container

import "time"

var deletionTimers = make(map[*Container]time.Time)
var deletionDelay = 2 * time.Second

func Check() {
	mutex.Lock()
	defer mutex.Unlock()
	now := time.Now()
	var toDelete []*Container
	lobbyDeleted := false
	containers := GetContainers()
	for _, c := range containers {
		if c.Tag == "lobby" {
			if c.Info.Players().Len() == 0 && len(c.Pending) == 0 && canDeleteLobbyServer(c) && !lobbyDeleted {
				if start, exists := deletionTimers[c]; !exists {
					deletionTimers[c] = now
				} else {
					if now.Sub(start) >= deletionDelay {
						toDelete = append(toDelete, c)
						delete(deletionTimers, c)
						lobbyDeleted = true
					}
				}
			} else {
				delete(deletionTimers, c)
			}
		} else {
			if c.Info.Players().Len() == 0 && len(c.Pending) == 0 {
				toDelete = append(toDelete, c)
			}
		}
	}
	for _, c := range toDelete {
		DeleteContainer(c)
	}
}

func canDeleteLobbyServer(exclude *Container) bool {
	for _, c := range GetContainers() {
		if c == exclude {
			continue
		}
		if c.Tag == "lobby" && c.Online && c.Info.Players().Len() < defaultMin {
			return true
		}
	}
	return false
}
