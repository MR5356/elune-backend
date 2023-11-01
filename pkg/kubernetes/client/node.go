package client

import (
	"context"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type NodeInfo struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	Age              string `json:"age"`
	Version          string `json:"version"`
	InternalIP       string `json:"internalIP"`
	ExternalIP       string `json:"externalIP"`
	OsImage          string `json:"osImage"`
	KernelVersion    string `json:"kernelVersion"`
	ContainerRuntime string `json:"containerRuntime"`
	OperatingSystem  string `json:"operatingSystem"`
	Architecture     string `json:"architecture"`
}

func (c *Client) GetNodes() (nodes []NodeInfo, err error) {
	list, err := c.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, node := range list.Items {
		nd := NodeInfo{
			Name:             node.Name,
			Status:           string(node.Status.Conditions[len(node.Status.Conditions)-1].Type),
			Age:              humanize.CustomRelTime(node.CreationTimestamp.Time, time.Now(), "", "", magnitudes),
			Version:          node.Status.NodeInfo.KernelVersion,
			OsImage:          node.Status.NodeInfo.OSImage,
			KernelVersion:    node.Status.NodeInfo.KernelVersion,
			ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
			OperatingSystem:  node.Status.NodeInfo.OperatingSystem,
			Architecture:     node.Status.NodeInfo.Architecture,
		}
		internalIP, externalIP := "", ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				if len(internalIP) > 0 {
					internalIP += ","
				}
				internalIP += addr.Address
				break
			} else if addr.Type == "ExternalIP" {
				if len(externalIP) > 0 {
					externalIP += ","
				}
				externalIP += addr.Address
				break
			}
		}

		nd.InternalIP = internalIP
		nd.ExternalIP = externalIP

		nodes = append(nodes, nd)
		logrus.Infof("nodes: %+v", nodes)

	}
	return
}
