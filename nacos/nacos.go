package nacos

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
)

type NacosConf struct {
	ServerAddr     string
	ServerPort     uint64
	DataId         string
	NamespaceId    string
	Group          string
	ApiServiceName string
	RpcServiceName string
	ClusterName    string
	Username       string
	Password       string
	LogDir         string
	CacheDir       string
	LogLevel       string
}

func (cfg NacosConf) CreateNacosClient() (naming_client.INamingClient, config_client.IConfigClient) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(cfg.ServerAddr, cfg.ServerPort),
	}
	cc := &constant.ClientConfig{
		NamespaceId:         cfg.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              cfg.LogDir,
		CacheDir:            cfg.CacheDir,
		LogLevel:            cfg.LogLevel,
		Username:            cfg.Username,
		Password:            cfg.Password,
	}
	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc})
	if err != nil {
		log.Fatalf("Failed to create Nacos naming client: %v", err)
		return nil, nil
	}
	configClient, err := clients.NewConfigClient(vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc})
	if err != nil {
		log.Fatalf("Failed to create Nacos configuration client: %v", err)
		return nil, nil
	}
	return namingClient, configClient
}

func (cfg NacosConf) RegisterServiceInstance(client naming_client.INamingClient, listenOn string, serviceName string) {
	host, portStr, err := net.SplitHostPort(listenOn)
	if err != nil {
		logx.Errorf("Failed to split listen on %s: %v", listenOn, err)
		return
	}
	port, err := strconv.ParseUint(portStr, 10, 64)
	if err != nil {
		logx.Errorf("Failed to parse listen on %s: %v", listenOn, err)
		return
	}
	if host == "0.0.0.0" {
		host = GetLocalIP()
	}
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          host,
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: cfg.ClusterName,
		GroupName:   cfg.Group,
		Metadata:    map[string]string{"version": "1.0"},
	})
	proc.AddShutdownListener(func() {
		_, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          host,
			Port:        port,
			ServiceName: serviceName,
			Cluster:     cfg.ClusterName,
			GroupName:   cfg.Group,
			Ephemeral:   true,
		})
		if err != nil {
			logx.Errorf("Failed to deregister instance: %v", err)
		} else {
			logx.Infof("Deregistered instance: %v", host)
		}
	})
}

func (cfg NacosConf) SelectOneHealthyService(client naming_client.INamingClient) string {
	param := vo.SelectOneHealthInstanceParam{
		ServiceName: cfg.RpcServiceName,
		GroupName:   cfg.Group,
		Clusters:    []string{cfg.ClusterName},
	}
	instance, err := client.SelectOneHealthyInstance(param)
	if err != nil {
		return ""
	}
	return instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
}

func (cfg NacosConf) GetConfig(client config_client.IConfigClient) (string, error) {
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: cfg.DataId,
		Group:  cfg.Group,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	return content, nil
}

func (cfg NacosConf) ListenConfig(client config_client.IConfigClient, onChange func(data string)) error {
	err := client.ListenConfig(vo.ConfigParam{
		DataId: cfg.DataId,
		Group:  cfg.Group,
		OnChange: func(namespace, group, dataId, data string) {
			onChange(data)
		},
	})
	if err != nil {
		return fmt.Errorf("failed to listen config: %w", err)
	}
	proc.AddShutdownListener(func() {
		err = client.CancelListenConfig(vo.ConfigParam{
			DataId: cfg.DataId,
			Group:  cfg.Group,
		})
		if err != nil {
			logx.Info("failed to cancel listen config")
		} else {
			logx.Info("cancelled listen config")
		}
	})
	return nil
}

func (cfg NacosConf) CancelListenConfig(client config_client.IConfigClient) error {
	err := client.CancelListenConfig(vo.ConfigParam{
		DataId: cfg.DataId,
		Group:  cfg.Group,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel listen config: %w", err)
	}
	return nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
