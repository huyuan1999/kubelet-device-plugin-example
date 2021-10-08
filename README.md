## kubelet-device-plugin-example

kubelet 设备插件实现示例

```bash
$ go build -o ./deployment/example-plugin
$ cd ./deployment/ && docker build -t example-plugin:v1
$ kubelet apply -f plugin.yaml
# 等待 plugin running 之后再执行
$ kubelet apply -f client.yaml
```
