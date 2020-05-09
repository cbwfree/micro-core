package srv

import (
	"fmt"
	"math/rand"
	"strings"
)

// 服务节点
type ServiceNode struct {
	UUID     string            `json:"uuid"`
	Id       string            `json:"id"`
	Version  string            `json:"version"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

// CheckServiceNode 检查服务节点是否存在
func CheckServiceNode(name string, nodeId string) bool {
	nameId := fmt.Sprintf("%s-%s", name, nodeId)

	for _, s := range APP().GetServices(name) {
		for _, n := range s.Nodes {
			if n.Id == nameId {
				return true
			}
		}
	}

	return false
}

// GetServiceNode 获取指定服务指定节点
func GetServiceNode(srvName string, nodeId string) *ServiceNode {
	prefix := fmt.Sprintf("%s-", srvName)
	nameId := fmt.Sprintf("%s-%s", srvName, nodeId)

	for _, s := range APP().GetServices(srvName) {
		for _, n := range s.Nodes {
			if n.Id == nameId {
				return &ServiceNode{
					UUID:     strings.Replace(n.Id, prefix, "", 1),
					Id:       n.Id,
					Version:  s.Version,
					Address:  n.Address,
					Metadata: n.Metadata,
				}
			}
		}
	}

	return nil
}

// GetServiceNodes 获取指定服务所有节点
func GetServiceNodes(srvName string) []*ServiceNode {
	var nodes []*ServiceNode

	prefix := fmt.Sprintf("%s-", srvName)

	for _, s := range APP().GetServices(srvName) {
		for _, n := range s.Nodes {
			nodes = append(nodes, &ServiceNode{
				UUID:     strings.Replace(n.Id, prefix, "", 1),
				Id:       n.Id,
				Version:  s.Version,
				Address:  n.Address,
				Metadata: n.Metadata,
			})
		}
	}

	return nodes
}

// GetServiceNodesVersion 获取指定版本服务节点
func GetServiceNodesVersion(srvName string, version string) []*ServiceNode {
	prefix := fmt.Sprintf("%s-", srvName)

	var nodes []*ServiceNode
	for _, s := range APP().GetServices(srvName) {
		if s.Version != version {
			continue
		}

		for _, n := range s.Nodes {
			nodes = append(nodes, &ServiceNode{
				UUID:     strings.Replace(n.Id, prefix, "", 1),
				Id:       n.Id,
				Version:  s.Version,
				Address:  n.Address,
				Metadata: n.Metadata,
			})
		}
	}

	return nodes
}

// GetServiceNodeIds 获取指定服务所有节点的ID列表
func GetServiceNodeIds(srvName string) []string {
	var nodes []string

	prefix := fmt.Sprintf("%s-", srvName)

	for _, s := range APP().GetServices(srvName) {
		for _, n := range s.Nodes {
			nodes = append(nodes, strings.ReplaceAll(n.Id, prefix, ""))
		}
	}
	return nodes
}

// 获取随机节点
func GetRandomServiceNode(srvName string) (*ServiceNode, error) {
	nodes := GetServiceNodes(srvName)
	length := len(nodes)
	if length == 0 {
		return nil, fmt.Errorf("not found %s service node", srvName)
	}
	if length == 1 {
		return nodes[0], nil
	}
	return nodes[rand.Intn(length)], nil
}

// 获取随机节点ID
func GetRandomNodeId(srvName string) (string, error) {
	if Name() == srvName {
		return Id(), nil
	}

	node, err := GetRandomServiceNode(srvName)
	if err != nil {
		return "", err
	}

	return node.UUID, nil
}
