package config

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/platform9/nodelet2go/pkg/consts"
)

type NodeletConfig struct {
	AllowWorkloadsOnMaster bool
	CalicoV4Interface      string
	CalicoV6Interface      string
	ClusterId              string
	HostId                 string
	MasterIp               string
	MasterVipEnabled       string
	MasterVipInterface     string
	Mtu                    string
	Privileged             string
	NodeletRole            string
}

func GenNodeletConfigLocal(host *NodeletConfig) error {
	nodeStateDir := filepath.Join(consts.ClusterStateDir, host.ClusterId, host.HostId)
	if _, err := os.Stat(nodeStateDir); os.IsNotExist(err) {
		os.MkdirAll(nodeStateDir, 0766)
	}

	nodeletCfgFile := filepath.Join(nodeStateDir, consts.NodeletConfigFile)

	t := template.Must(template.New(host.HostId).Parse(nodeletConfigTmpl))

	fd, err := os.Create(nodeletCfgFile)
	defer fd.Close()
	if err != nil {
		return fmt.Errorf("Failed to Create nodelet config File: %s err: %s", nodeletCfgFile, err)
	}

	err = t.Execute(fd, host)
	if err != nil {
		return fmt.Errorf("template.Execute failed for file: %s err: %s\n", nodeletCfgFile, err)
	}

	return nil
}
