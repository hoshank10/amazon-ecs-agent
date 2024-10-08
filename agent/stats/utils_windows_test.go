//go:build windows && unit
// +build windows,unit

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerStatsToContainerStatsZeroCoresGeneratesError(t *testing.T) {
	numCores = uint64(0)
	// not using windows_test_stats.json here to save file open/read time
	jsonStat := fmt.Sprintf(`
		{
			"cpu_stats":{
				"cpu_usage":{
					"total_usage":%d
				}
			}
		}`, 100)
	dockerStat := &types.StatsJSON{}
	json.Unmarshal([]byte(jsonStat), dockerStat)
	err := validateDockerStats(dockerStat, false)
	assert.Error(t, err, "expected error converting container stats with zero cpu cores")
}

func TestValidateDockerStatsZeroValueReadTime(t *testing.T) {
	testCases := []struct {
		name                          string
		containerEnabledRestartPolicy bool
	}{
		{
			name:                          "container does not have restart policy enabled",
			containerEnabledRestartPolicy: false,
		},
		{
			name:                          "container has restart policy enabled",
			containerEnabledRestartPolicy: true,
		},
	}
	inputJsonFile, _ := filepath.Abs("./windows_test_stats.json")
	jsonBytes, err := os.ReadFile(inputJsonFile)
	assert.NoError(t, err)
	dockerStat := &types.StatsJSON{}
	json.Unmarshal(jsonBytes, dockerStat)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set read time of docker stat to zero value of time.Time.
			dockerStat.Read = time.Time{}
			err = validateDockerStats(dockerStat, tc.containerEnabledRestartPolicy)
			if tc.containerEnabledRestartPolicy {
				assert.Error(t, err)
				assert.ErrorContains(t, err, invalidStatZeroValueReadTimeMsg)
			} else {
				assert.NoError(t, nil)
			}
		})
	}
}

func TestDockerStatsToContainerStats(t *testing.T) {
	numCores = 4
	inputJsonFile, _ := filepath.Abs("./windows_test_stats.json")
	jsonBytes, _ := ioutil.ReadFile(inputJsonFile)
	dockerStat := &types.StatsJSON{}
	json.Unmarshal([]byte(jsonBytes), dockerStat)
	containerStats, err := dockerStatsToContainerStats(dockerStat)
	assert.NoError(t, err, "converting container stats failed")
	require.NotNil(t, containerStats, "containerStats should not be nil")
	netStats := containerStats.networkStats
	assert.NotNil(t, netStats, "networkStats should not be nil")
	validateNetworkMetrics(t, netStats)
	assert.Equal(t, uint64(2500), containerStats.cpuUsage,
		"unexpected value for cpuUsage", containerStats.cpuUsage)
	assert.Equal(t, uint64(3), containerStats.storageReadBytes,
		"unexpected value for storageReadBytes", containerStats.storageReadBytes)
	assert.Equal(t, uint64(15), containerStats.storageWriteBytes,
		"Unexpected value for storageWriteBytes", containerStats.storageWriteBytes)

}
