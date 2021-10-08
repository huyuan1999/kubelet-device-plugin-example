// @Time    : 2021/9/28 2:57 下午
// @Author  : HuYuan
// @File    : main.go
// @Email   : huyuan@virtaitech.com

package main

import (
	"github.com/fsnotify/fsnotify"
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
		log.Fatalf("创建 unix socket 文件错误: %s", err.Error())
	}

	defer func() { _ = os.Remove(pluginSock) }()

	// 向 kubelet 注册插件 grpc server
	pluginapi.RegisterDevicePluginServer(server, &examplePlugin{})

	if err := server.Serve(sock); err != nil {
		log.Fatalf("启动 grpc server 错误: %s", err.Error())
	}
}

func checkRPCServer(socket string, timeout time.Duration)  {
	conn, err := dialKubelet(socket, timeout)
	if err != nil {
		log.Fatalf("dialKubelet %s %s", pluginSock, err.Error())
	}
	_ = conn.Close()
}

func watcherKubeletSocket()  {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("初始化 fsnotify watcher 错误: ", err.Error())
	}

	err = watcher.Add(pluginapi.DevicePluginPath)
	if err != nil {
		log.Fatalf("add watcher file %s error: %s\n", pluginapi.DevicePluginPath, err.Error())
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				// 退出插件程序, 由 kubelet 重新拉起
				log.Fatalf("inotify: %s created, restarting.", pluginapi.KubeletSocket)
			}
		case err := <-watcher.Errors:
			log.Printf("inotify: %s", err)
		}
	}
}

func main() {
	go watcherKubeletSocket()
	go serve()
	checkRPCServer(pluginSock, time.Second * 30)
	plugin := examplePlugin{}
	err := plugin.Register()
	if err != nil {
		log.Fatal("Register error", err)
	}
	select {

	}
}

