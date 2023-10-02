package main

import (
	"flag"

	node_exporter_main "github.com/chaolihf/node_exporter"
	"github.com/elastic/beats/v7/filebeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
)

/*
创建包含Node_Export相关的参数
*/
func getSettingsIncludeNodeExporter() instance.Settings {
	go node_exporter_main.Main()
	flag.String("web.config.file", "", "[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
	flag.String("web.listen-address", ":9100", "Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.")
	settings := cmd.FilebeatSettings()
	runFlags := settings.RunFlags
	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.config.file"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.listen-address"))
	return settings
}
