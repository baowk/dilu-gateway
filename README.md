# dilu-gateway
golang原生实现的一个代理，通过配置来转发请求。
实现了一个基于dilu的jwt和权限实现。

## 安装
```bash
go mod tidy
````

## 启动
```bash
go run main.go c --config resources/config.dev.yaml 
````

## 配置
详见配置文件，resources\config.dev.yaml