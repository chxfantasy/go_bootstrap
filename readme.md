# 编译
go build --tags=jsoniter main.go
# 启动命令
./go_bootstrap -env dev

// 使用配置文件 conf/dev.conf.yaml
// 线上环境，切换成 conf/prod.conf.yaml，./go_bootstrap -env prod