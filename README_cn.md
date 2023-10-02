# 集成NodeExporter等插件的统一采集客户端

# 特性
集成FileBeat（官方）
集成Node_Exporter（官方）
集成Process,Network,ShellScript采集器（参考https://github.com/chaolihf/OneAgent)

# 安装
go get github.com/chaolihf/node_exporter
go get github.com/Shopify/sarama@v1.37.2
go get github.com/elastic/beats/v7/libbeat/publisher/queue/diskqueue
在import中增加node_expoter_main "github.com/chaolihf/node_exporter"
在main第一行增加go node_expoter_main.Main()

# 问题
由于node_exporter，filebeat都会自己解析参数，并且参数不对的时候就会报错，这个就导致无法正确使用参数来运行。

# 运行
./filebeat