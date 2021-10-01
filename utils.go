// @Time    : 2021/10/1 5:46 下午
// @Author  : HuYuan
// @File    : utils.go
// @Email   : huyuan@virtaitech.com

package main

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"time"
)

// 创建 kubelet device plugin client
func dialKubelet(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := grpc.DialContext(ctx, unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}
