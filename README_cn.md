# 集成NodeExporter等插件的统一采集客户端

# 特性
集成FileBeat（官方）  
集成Node_Exporter（官方）  
集成Process,Network,ShellScript采集器（参考https://github.com/chaolihf/OneAgent)  

# 安装
go get github.com/chaolihf/node_exporter  
go get github.com/Shopify/sarama@v1.37.2  
go get github.com/elastic/beats/v7/libbeat/publisher/queue/diskqueue  

修改冲突的模块beats/filebeat/fileset/factory.go，将模块名称参数module修改为filebeat-module(应该是用于监控的)
增加oneagent.go文件，filebeat.yml,auditbeat.yml,config.json文件

# 配置
config.json文件格式
```
{
    "module":["filebeat","node_exporter"]
}
```

module中可以激活0个或3个模块

# 编译和运行
cd beats/oneagent
go build -ldflags '-w -s'  -o oneagent
./oneagent

需要给此用户以root用户运行的权限，可以在sudo中执行此应用
Cmnd_Alias ONEAGENT_COMMANDS = /home/oneops/oneagent
User_Alias ONEAGENT_USERS = oneops
ONEAGENT_USERS    ALL=(ALL:ALL) ONEAGENT_COMMANDS


# 问题
## 1、命令行参数缺失
由于node_exporter，filebeat都会自己解析参数，并且参数不对的时候就会报错，这个就导致无法正确使用参数来运行。  
现在已支持的参数包括两个都支持的参数；  
--web.config.file  
--web.listen-address  

## 2、为什么不集成auditbeat等
曾经做个尝试，但由于两个原因放弃 
首先是两个模块都会用到metricbeat，从而导致多个变量、参数冲突，需要逐个解决，改动较大
其次时auditbeat需要root权限，而filebeat,node是建议在普通用户下运行，这样会将整个应用放在root下运行，安全性存在一定的问题