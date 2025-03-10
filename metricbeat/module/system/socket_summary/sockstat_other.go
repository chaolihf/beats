// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !linux

package socket_summary

import (
	"github.com/chaolihf/gopsutil/net"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/elastic-agent-system-metrics/metric/system/resolve"
)

// a stub function for non-linux systems
// get a list of platform-specific enhancements and apply them to our mapStr object.
func applyEnhancements(data mapstr.M, sys resolve.Resolver) (mapstr.M, error) {
	return data, nil
}

// connections gets connection information
func connections(kind string) ([]net.ConnectionStat, error) {
	return net.Connections(kind)
}
