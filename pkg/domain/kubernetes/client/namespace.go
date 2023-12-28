package client

import (
	"context"
	"fmt"
	"github.com/MR5356/elune-backend/pkg/utils/structutil"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNamespace() {
	list, err := c.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("list: %+v", structutil.Struct2String(list.Items))
	node := list.Items[0]
	logrus.Infof("%s\t%s\t%s\t%s\t%s\t%s\t,\n",
		node.Name,
		node.Status.Addresses[0].Address,
		node.Status.NodeInfo.OSImage,
		node.Status.NodeInfo.KubeletVersion,
		node.Status.NodeInfo.OperatingSystem,
		node.Status.NodeInfo.Architecture)

	namespaceList, err := c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, namespace := range namespaceList.Items {
		fmt.Println("k8s namespace")
		fmt.Printf("%s\t%+v\t%+v\t%s\t\n",
			namespace.Name,
			namespace.Status.Phase,
			namespace.CreationTimestamp.Day(),
			namespace.Status.Phase,
		)
	}
}
