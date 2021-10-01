// @Time    : 2021/9/28 2:57 下午
// @Author  : HuYuan
// @File    : main.go
// @Email   : huyuan@virtaitech.com

package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"net"
	"os"
	"path"
	"time"
)

var pluginSock = path.Join(pluginapi.DevicePluginPath, "example.sock")

// 启动插件的 grpc 服务器，并向 kubelet 注册
func serve()  {
	_ = os.Remove(pluginSock)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	sock, err := net.Listen("unix", pluginSock)
	if err != nil {
		logrus.Fatalf("创建 unix socket 文件错误: %s", err.Error())
	}

	defer func() { _ = os.Remove(pluginSock) }()

	// 向 kubelet 注册插件 grpc server
	pluginapi.RegisterDevicePluginServer(server, &examplePlugin{})

	if err := server.Serve(sock); err != nil {
		logrus.Fatalf("启动 grpc server 错误: %s", err.Error())
	}
	conn, err := dialKubelet(pluginSock, 5*time.Second)
	if err != nil {
		log.Fatalf("dialKubelet %s %s", pluginSock, err.Error())
	}
	_ = conn.Close()
}

func main() {
	plugin := examplePlugin{}
	err := plugin.Register()
	if err != nil {
		log.Fatal("Register error", err)
	}
	serve()
}





