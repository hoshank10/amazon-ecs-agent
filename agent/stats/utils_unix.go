//go:build !windows
// +build !windows

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package stats

import (
	"fmt"

	"github.com/aws/amazon-ecs-agent/agent/config"
	"github.com/cihub/seelog"
	"github.com/docker/docker/api/types"
)

// dockerStatsToContainerStats returns a new object of the ContainerStats object from docker stats.
func dockerStatsToContainerStats(dockerStats *types.StatsJSON) (*ContainerStats, error) {
	cpuUsage := dockerStats.CPUStats.CPUUsage.TotalUsage / numCores
	memoryUsage := getMemUsage(dockerStats.MemoryStats)
	storageReadBytes, storageWriteBytes := getStorageStats(dockerStats)
	networkStats := getNetworkStats(dockerStats)
	return &ContainerStats{
		cpuUsage:          cpuUsage,
		memoryUsage:       memoryUsage,
		storageReadBytes:  storageReadBytes,
		storageWriteBytes: storageWriteBytes,
		networkStats:      networkStats,
		timestamp:         dockerStats.Read,
	}, nil
}

func getMemUsage(mem types.MemoryStats) uint64 {
	if config.CgroupV2 {
		// for cgroupv2 systems, mem usage calculation uses the same method that the docker cli uses
		// https://github.com/docker/cli/blob/e198123693b1aaa724041fff602c7d75c8fe4b57/cli/command/container/stats_helpers.go#L227-L249
		// see https://github.com/aws/amazon-ecs-agent/issues/3323
		if v, ok := mem.Stats["inactive_file"]; ok && v < mem.Usage {
			return mem.Usage - v
		}
	}
	if v, ok := mem.Stats["cache"]; ok && v < mem.Usage {
		return mem.Usage - v
	}
	return mem.Usage
}

func validateDockerStats(dockerStats *types.StatsJSON, containerEnabledRestartPolicy bool) error {
	if containerEnabledRestartPolicy && dockerStats.Read.IsZero() {
		return fmt.Errorf("invalid container statistics reported for container with restart policy enabled, %s",
			invalidStatZeroValueReadTimeMsg)
	}

	if config.CgroupV2 {
		// PercpuUsage is not available in cgroupv2
		if numCores == uint64(0) {
			return fmt.Errorf("invalid number of cores returned from runtime.NumCPU, numCores=0")
		}
	} else {
		// The length of PercpuUsage represents the number of cores in an instance.
		if len(dockerStats.CPUStats.CPUUsage.PercpuUsage) == 0 || numCores == uint64(0) {
			return fmt.Errorf("invalid container statistics reported, no cpu core usage reported")
		}
	}
	return nil
}

func getStorageStats(dockerStats *types.StatsJSON) (uint64, uint64) {
	// initialize block io and loop over stats to aggregate
	if dockerStats.BlkioStats.IoServiceBytesRecursive == nil {
		seelog.Debug("Storage stats not reported for container")
		return uint64(0), uint64(0)
	}
	storageReadBytes := uint64(0)
	storageWriteBytes := uint64(0)
	for _, blockStat := range dockerStats.BlkioStats.IoServiceBytesRecursive {
		switch op := blockStat.Op; op {
		case "Read":
			storageReadBytes += blockStat.Value
		case "Write":
			storageWriteBytes += blockStat.Value
		default:
			//ignoring "Async", "Total", "Sum", etc
			continue
		}
	}
	return storageReadBytes, storageWriteBytes
}
