package containers

import (
	"fmt"
	"starfield/plugins/starfield/containers/node"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type ContainerRecord struct {
	Name string
	Tag  string
	Node node.Node
	Info proxy.RegisteredServer
}

type Containers struct {
	Nodes      map[string]node.Node
	Containers map[string]ContainerRecord
}

var GlobalContainers = &Containers{
	Nodes:      make(map[string]node.Node),
	Containers: make(map[string]ContainerRecord),
}

func (c *Containers) AddNode(name string, n node.Node) {
	c.Nodes[name] = n
}

func (c *Containers) AddContainer(containerName, tag, nodeName string, info proxy.RegisteredServer) error {
	n, ok := c.Nodes[nodeName]
	if !ok {
		return fmt.Errorf("node %q not found", nodeName)
	}
	c.Containers[containerName] = ContainerRecord{Name: containerName, Tag: tag, Node: n, Info: info}
	return nil
}

func (c *Containers) RemoveContainer(containerName string) error {
	_, exists := c.Containers[containerName]
	if !exists {
		return fmt.Errorf("container %s not found", containerName)
	}
	delete(c.Containers, containerName)
	return nil
}

func (c *Containers) GetContainer(containerName string) (ContainerRecord, error) {
	container, exists := c.Containers[containerName]
	if !exists {
		return ContainerRecord{}, fmt.Errorf("container %s not found", containerName)
	}
	return container, nil
}

func (c *Containers) GetAllContainers() []ContainerRecord {
	list := make([]ContainerRecord, 0, len(c.Containers))
	for _, rec := range c.Containers {
		list = append(list, rec)
	}
	return list
}

func (c *Containers) GetAllContainerbyTag(tag string) []ContainerRecord {
	list := make([]ContainerRecord, 0, len(c.Containers))
	for _, rec := range c.Containers {
		if rec.Tag == tag {
			list = append(list, rec)
		}
	}
	return list
}
