package cluster

import (
	"bytes"
	"fmt"
	"io/ioutil"

	nodeletconfig "github.com/platform9/nodelet2go/pkg/config"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type BootstrapConfig struct {
	SSHUser                string        `json:"sshUser,omitempty"`
	SSHPrivateKeyFile      string        `json:"sshPrivateKeyFile,omitempty"`
	ClusterId              string        `json:"clusterName,omitempty"`
	AllowWorkloadsOnMaster bool          `json:"allowWorksloadsOnMaster,omitempty"`
	MasterIp               string        `json:"masterIp,omitempty"`
	MasterVipEnabled       string        `json:"masterVipEnabled,omitempty"`
	MasterVipInterface     string        `json:"masterVipInterface,omitempty"`
	CalicoV4Interface      string        `json:"calicoV4Interface,omitempty"`
	CalicoV6Interface      string        `json:"calicoV6Interface,omitempty"`
	MTU                    string        `json:"mtu,omitempty"`
	Privileged             string        `json:"privileged,omitempty"`
	MasterNodes            []NodeletHost `json:"masterNodes"`
	WorkerNodes            []NodeletHost `json:"workerNodes"`
}

type NodeletHost struct {
	NodeName            string `json:"nodeName"`
	V4InterfaceOverride string `json:calicoV4Interface,omitempty"`
	V6InterfaceOverride string `json:calicoV6Interface,omitempty"`
}

func CreateCluster(cfgPath string) error {
	clusterCfg, err := ParseBootstrapConfig(cfgPath)
	if err != nil {
		fmt.Printf("Failed to Parse Cluster Config: %s", err)
		return fmt.Errorf("Failed to Parse Cluster Config: %s", err)
	}

	sshKey, err := ioutil.ReadFile(clusterCfg.SSHPrivateKeyFile)
	if err != nil {
		return fmt.Errorf("Failed to read private key: %s", clusterCfg.SSHPrivateKeyFile)
	}
	fmt.Printf("Got sshKey: %s", sshKey)

	return nil
}

func ParseBootstrapConfig(cfgPath string) (*BootstrapConfig, error) {
	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening bootstrap config file: %s", cfgFile)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(cfgFile), 4096)
	bootstrapConfig := &BootstrapConfig{
		AllowWorkloadsOnMaster: false,
		CalicoV4Interface:      "first-found",
		ClusterId:              "airctl-mgmt",
		SSHUser:                "root",
		SSHPrivateKeyFile:      "/root/.ssh/id_rsa",
		Privileged:             "true",
		MTU:                    "1440",
	}

	err = decoder.Decode(bootstrapConfig)
	if err != nil {
		return nil, fmt.Errorf("Error decoding bootstrap config\n")
	}

	return bootstrapConfig, nil
}

func GenClusterState(cfg *BootstrapConfig) error {
	for _, host := range cfg.MasterNodes {
		var nodeletCfg *nodeletconfig.NodeletConfig
		nodeletCfg = new(nodeletconfig.NodeletConfig)
		nodeletCfg.AllowWorkloadsOnMaster = cfg.AllowWorkloadsOnMaster
		nodeletCfg.CalicoV4Interface = cfg.CalicoV4Interface
		nodeletCfg.CalicoV6Interface = cfg.CalicoV6Interface
		nodeletCfg.ClusterId = cfg.ClusterId
		nodeletCfg.HostId = host.NodeName
		nodeletCfg.Mtu = cfg.MTU
		nodeletCfg.Privileged = cfg.Privileged
		nodeletCfg.NodeletRole = "master"

		err := nodeletconfig.GenNodeletConfigLocal(nodeletCfg)
		if err != nil {
			fmt.Printf("Failed to generate config: %s", err)
			return fmt.Errorf("Failed to generate config: %s", err)
		}
	}

	return nil
}
