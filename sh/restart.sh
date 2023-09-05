chmod +x ./dilu-gateway
echo "kill dilu-gateway service"
killall dilu-gateway # kill dilu-gateway service
nohup ./dilu-gateway -c config.test.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run dilu-gateway success"
ps -aux | grep dilu-gateway