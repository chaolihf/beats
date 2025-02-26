package main

import (
	"flag"
	"fmt"
	"os"

	node_exporter_main "github.com/chaolihf/node_exporter"
	jjson "github.com/chaolihf/udpgo/json"
	"github.com/chaolihf/udpgo/lang"
	"github.com/containerd/cgroups/v3/cgroup1"
	fileBeatCmd "github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var fileLogger *zap.Logger

const (
	Module_File          = "filebeat"
	Module_Node          = "node_exporter"
	Module_LimitResource = "limit"
)

/*
命令行参数信息
*/
type CommandInfo struct {
	name        string
	valueType   string
	value       string
	description string
}

/*
限制资源信息
*/
type LimitResourceInfo struct {
	shares uint64
	period uint64
	quota  int64
	memory int64
	swap   int64
}

/*
根据配置文件激活不同的模块，获取命令行信息
*/
func activeInfos() (map[string]bool, []CommandInfo, LimitResourceInfo) {
	var oneAgentModules = make(map[string]bool)
	var commandInfos []CommandInfo
	var resourceInfo LimitResourceInfo
	oneAgentModules[Module_File] = false
	oneAgentModules[Module_Node] = false
	filePath := "config.json"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件%s出错:%s\n", filePath, err.Error())
	} else {
		jsonConfigInfos, err := jjson.NewJsonObject([]byte(content))
		if err != nil {
			fmt.Printf("文件%sJson格式出错:%s\n", filePath, err.Error())
		} else {
			jsonModuleInfos := jsonConfigInfos.GetJsonArray("module")
			for _, jsonModuleInfo := range jsonModuleInfos {
				oneAgentModules[jsonModuleInfo.GetStringValue()] = true
			}
			jsonCommandInfos := jsonConfigInfos.GetJsonArray("command")
			for _, jsonCommandInfo := range jsonCommandInfos {
				commandInfos = append(commandInfos, CommandInfo{
					name:        jsonCommandInfo.GetString("name"),
					valueType:   jsonCommandInfo.GetString("type"),
					value:       jsonCommandInfo.GetString("value"),
					description: jsonCommandInfo.GetString("description"),
				})
			}
			jsonLimitResourceInfo := jsonConfigInfos.GetJsonObject("limit")
			resourceInfo = LimitResourceInfo{
				shares: uint64(jsonLimitResourceInfo.GetInt("shares")),
				period: uint64(jsonLimitResourceInfo.GetInt("period")),
				quota:  int64(jsonLimitResourceInfo.GetInt("quota")),
				memory: int64(jsonLimitResourceInfo.GetInt("memory")),
				swap:   int64(jsonLimitResourceInfo.GetInt("swap")),
			}
		}
	}
	for moduleName := range oneAgentModules {
		fmt.Printf("激活模块%s\n", moduleName)
	}
	return oneAgentModules, commandInfos, resourceInfo
}

/*
初始化日志设置，包括日志文件路径、日志级别、日志保留天数、日志文件大小限制等信息
*/
func init() {
	fileLogger = lang.InitProductLogger("logs/oneagent.log", 300, 3, 10)
}

/*
开始应用
*/
func main() {
	defer func() {
		fileLogger.Sync()
		if r := recover(); r != nil {
			fileLogger.Info(fmt.Sprintf("程序退出原因:", r))
			fmt.Println("程序退出原因:", r)
		} else {
			fileLogger.Info("程序正常退出")
			fmt.Println("程序正常退出")
		}
	}()
	fileDone := make(chan string)
	oneAgentModules, oneAgentCommandLines, oneAgentResource := activeInfos()
	if oneAgentModules[Module_LimitResource] {
		limitResource(oneAgentResource)
	}
	if oneAgentModules[Module_File] {
		go runFileBeat(oneAgentModules, fileDone, oneAgentCommandLines)
		fileDoneInfo := <-fileDone
		fmt.Println(fileDoneInfo)
		if oneAgentModules[Module_Node] {
			fmt.Println("start node exporter")
			fileLogger.Info("启动node exporter")
			node_exporter_main.Main(fileLogger)
		}
	} else {
		if oneAgentModules[Module_Node] {
			node_exporter_main.Main(fileLogger)
			fileLogger.Info("启动node exporter")
		}
	}
}

/*
use cgroup to limit process resource,such as cpu,memory
*/
func limitResource(oneAgentResource LimitResourceInfo) {
	control, err := cgroup1.New(cgroup1.StaticPath("/oneagent"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &oneAgentResource.shares,
			Quota:  &oneAgentResource.quota,
			Period: &oneAgentResource.period,
		},
		Memory: &specs.LinuxMemory{
			Limit: &oneAgentResource.memory,
			Swap:  &oneAgentResource.swap,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer control.Delete()
		if err := control.Add(cgroup1.Process{Pid: os.Getpid()}, cgroup1.Cpu, cgroup1.Memory); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("enable limit resource ")
		}
	}
}

/*
创建包含Node_Export相关的参数
*/
func addNodeExporterFlag(runFlags *pflag.FlagSet, commandInfos []CommandInfo) {
	for _, commandInfo := range commandInfos {
		findFlag := flag.CommandLine.Lookup(commandInfo.name)
		if findFlag == nil {
			switch commandInfo.valueType {
			case "string":
				{
					flag.String(commandInfo.name, commandInfo.value, commandInfo.description)
					break
				}
			case "boolean":
				{
					flag.Bool(commandInfo.name, commandInfo.value == "true", commandInfo.description)
					break
				}
			}

		}
		runFlags.AddGoFlag(flag.CommandLine.Lookup(commandInfo.name))
	}
}

// func addNodeExporterConfigFileFlag(runFlags *pflag.FlagSet) {
// 	findFlag := flag.CommandLine.Lookup("web.config.file")
// 	if findFlag == nil {
// 		flag.String("web.config.file", "", "[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
// 	}
// 	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.config.file"))
// }

// func addNodeExporterListenAddressFlag(runFlags *pflag.FlagSet) {
// 	findFlag := flag.CommandLine.Lookup("web.listen-address")
// 	if findFlag == nil {
// 		flag.String("web.listen-address", ":9100", "Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.")
// 	}
// 	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.listen-address"))
// }

/*
运行FileBeat
*/
func runFileBeat(activeModules map[string]bool, doneChain chan<- string, commandInfos []CommandInfo) {
	defer func() {
		fileLogger.Sync()
		if r := recover(); r != nil {
			fileLogger.Info(fmt.Sprintf("filebeat退出原因:", r))
			fmt.Println("filebeat退出原因:", r)
		} else {
			fileLogger.Info("filebeat正常退出")
			fmt.Println("filebeat正常退出")
		}
	}()
	settings := fileBeatCmd.FilebeatSettings()
	if activeModules[Module_Node] {
		addNodeExporterFlag(settings.RunFlags, commandInfos)
	}
	doneChain <- "start filebeat"
	if err := fileBeatCmd.Filebeat(inputs.Init, settings).Execute(); err != nil {
		doneChain <- " file error " + err.Error()
	}
}
