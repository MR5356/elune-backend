package client

import (
	"crypto/md5"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/utils/fileutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type Client struct {
	client *kubernetes.Clientset
	config *rest.Config
}

func New(kubeconfig ...string) (client *Client, err error) {
	var config *rest.Config
	if len(kubeconfig) > 0 {
		filename := fmt.Sprintf("/tmp/kubeconfig-%x", md5.Sum([]byte(kubeconfig[0])))
		if err := fileutil.WriteToFile(filename, []byte(kubeconfig[0])); err != nil {
			return nil, err
		}
		defer os.Remove(filename)
		config, err = clientcmd.BuildConfigFromFlags("", filename)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	c, err := kubernetes.NewForConfig(config)

	return &Client{
		client: c,
		config: config,
	}, err
}

func (c *Client) GetVersion() (string, error) {
	info, err := c.client.DiscoveryClient.ServerVersion()
	if err != nil {
		return "", err
	}
	return info.String(), err
}
