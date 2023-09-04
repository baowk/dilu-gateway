chmod +x ./go-gateway
echo "kill go-gateway service"
killall go-gateway # kill go-gateway service
nohup ./go-gateway -c config.test.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run go-gateway success"
ps -aux | grep go-gateway