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
    "module":["filebeat","node_exporter","limit"],
    "command":[
        {
            "name":"web.config.file",
            "type":"string",
            "value":"",
            "description":"[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md"
        }
    ],
    "limit":{
        "shares":20,//设置CPU资源的共享权重。这里将共享权重设置为20，表示该控制组在CPU资源竞争时会获得较低的优先级。
	    "period":1000000,//设置CPU时间周期的长度。这里将周期设置为1秒
	    "quota":100000,//设置CPU使用时间配额。这里将配额设置为0.1秒
	    "memory":100000000,//设置内存限制。这里将内存限制设置为100MB。
	    "swap":1000000000//设置交换空间限制。这里将交换空间限制设置为1GB。
    },
}
```

module中可以激活0个或3个模块
type可以为string,boolean值

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

## 3、为什么开启filebeat模块后无法执行shell，出现operation not permitted错误
是因为filebeat默认启动seccomp模块，可以在配置文件中seccomp.enabled false禁用