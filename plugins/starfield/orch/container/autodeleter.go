package container

func Check() {
	for _, c := range GetContainers() {
		if c.Tag != "lobby" && c.Info.Players().Len() == 0 && len(c.Pending) == 0 {
			DeleteContainer(c)
		}
	}
}
