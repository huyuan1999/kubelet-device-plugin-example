// @Time    : 2021/9/28 5:16 下午
// @Author  : HuYuan
// @File    : plugin.go
// @Email   : huyuan@virtaitech.com

package main

import (
	"context"
	"github.com/google/uuid"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"time"
)

type examplePlugin struct {
}

// GetDevicePluginOptions
// 插件并非必须为 GetPreferredAllocation() 或 PreStartContainer() 提供有用的实现逻辑，调用 GetDevicePluginOptions() 时所返回的 DevicePluginOptions 消息中应该设置这些调用是否可用。
// kubelet 在真正调用这些函数之前，总会调用 GetDevicePluginOptions() 来查看是否存在这些可选的函数。
func (e *examplePlugin) GetDevicePluginOptions(ctx context.Context, empty *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// ListAndWatch 返回 Device 列表构成的数据流。
// 当 Device 状态发生变化或者 Device 消失时，ListAndWatch 会返回新的列表
func (e *examplePlugin) ListAndWatch(empty *pluginapi.Empty, watch pluginapi.DevicePlugin_ListAndWatchServer) error {
	var devices []*pluginapi.Device
	// 切片里面有几个元素则 kubelet 认为这个节点里面有多少此内设备
	devices = append(devices, &pluginapi.Device{
		ID:                   uuid.NewString(),
		Health:               pluginapi.Healthy,
	})
	if err := watch.Send(&pluginapi.ListAndWatchResponse{Devices: devices}); err != nil {
		log.Println("watch.Send error ", err.Error())
	}
	// ListAndWatch 应该保持阻塞以确保 kubelet 可以实时的获取插件信息
	// 如果 ListAndWatch 函数发生错误异常退出，则应该使整个插件退出，否则 kubelet 将会抛出 EOF 和 sync 错误
	select {
	}
	//return nil
}

// GetPreferredAllocation 从一组可用的设备中返回一些优选的设备用来分配，
// 所返回的优选分配结果不一定会是设备管理器的最终分配方案。
// 此接口的设计仅是为了让设备管理器能够在可能的情况下做出更有意义的决定。
func (e *examplePlugin) GetPreferredAllocation(ctx context.Context, preferred *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse,error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Allocate 在容器创建期间调用，这样设备插件可以运行一些特定于设备的操作，
// 并告诉 kubelet 如何令 Device 可在容器中访问的所需执行的具体步骤
func (e *examplePlugin) Allocate(ctx context.Context, allocate *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse,error) {
	allocations := make([]*pluginapi.ContainerAllocateResponse, len(allocate.ContainerRequests))
	for i := range allocations {
		allocations[i] = &pluginapi.ContainerAllocateResponse{}
	}
	return &pluginapi.AllocateResponse{
		ContainerResponses: allocations,
	}, nil
}

// PreStartContainer 在设备插件注册阶段根据需要被调用，调用发生在容器启动之前。
// 在将设备提供给容器使用之前，设备插件可以运行一些诸如重置设备之类的特定于具体设备的操作
func (e *examplePlugin) PreStartContainer(ctx context.Context, before *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse,error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

// Register 不是 kubelet interface 中定义需要实现的方法
// 但作为一个插件应该向 kubelet 汇报资源, 否则这个插件将没有任何意义
// 向 kubelet 注册插件信息
func (e *examplePlugin) Register() error {
	// 获取 kubelet device plugin client
	conn, err := dialKubelet(pluginapi.KubeletSocket, 5 * time.Second)
	if err != nil {
		return err
	}
	// 向 kubelet 注册插件信息
	client := pluginapi.NewRegistrationClient(conn)
	register := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     "example.sock",
		ResourceName: "test.com/huyuan",
	}
	_, err = client.Register(context.Background(), register)
	if err != nil {
		return err
	}
	return nil
}

