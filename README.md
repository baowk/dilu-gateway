# dilu-gateway
golang原生实现的一个代理，通过配置来转发请求。
实现了一个基于dilu的jwt和权限实现。

## 目的
我们以往的项目做了微服务的拆分，但是我们没有使用第三方的微服务架构，为了统一权限等，基于以往经验，模仿了spring cloud gateway的网管。
结合[dilu-rd](https://github.com/baowk/dilu-rd)的服务注册发现模块，已经实现了微服务。

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