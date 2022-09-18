// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package accelerators

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/stats"

	"k8s.io/klog/v2"
)

type nvidiaManager struct{}

func NewNvidiaManager(includedMetrics container.MetricSet) stats.Manager {
	if !includedMetrics.Has(container.AcceleratorUsageMetrics) {
		klog.V(2).Info("NVIDIA GPU metrics disabled")
		return &stats.NoopManager{}
	}

	manager := &nvidiaManager{}
	return manager
}

// Destroy shuts down NVML.
func (nm *nvidiaManager) Destroy() {
}

// GetCollector returns a collector that can fetch NVIDIA gpu metrics for NVIDIA devices
// present in the devices.list file in the given devicesCgroupPath.
func (nm *nvidiaManager) GetCollector(devicesCgroupPath string) (stats.Collector, error) {
	return &stats.NoopCollector{}, nil
}

// parseDevicesCgroup parses the devices cgroup devices.list file for the container
// and returns a list of minor numbers corresponding to NVIDIA GPU devices that the
// container is allowed to access. In cases where the container has access to all
// devices or all NVIDIA devices but the devices are not enumerated separately in
// the devices.list file, we return an empty list.
// This is defined as a variable to help in testing.
var parseDevicesCgroup = func(devicesCgroupPath string) ([]int, error) {
	// Always return a non-nil slice
	nvidiaMinorNumbers := []int{}

	devicesList := filepath.Join(devicesCgroupPath, "devices.list")
	f, err := os.Open(devicesList)
	if err != nil {
		return nvidiaMinorNumbers, fmt.Errorf("error while opening devices cgroup file %q: %v", devicesList, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	// See https://www.kernel.org/doc/Documentation/cgroup-v1/devices.txt for the file format
	for s.Scan() {
		text := s.Text()

		fields := strings.Fields(text)
		if len(fields) != 3 {
			return nvidiaMinorNumbers, fmt.Errorf("invalid devices cgroup entry %q: must contain three whitespace-separated fields", text)
		}

		// Split the second field to find out major:minor numbers
		majorMinor := strings.Split(fields[1], ":")
		if len(majorMinor) != 2 {
			return nvidiaMinorNumbers, fmt.Errorf("invalid devices cgroup entry %q: second field should have one colon", text)
		}

		// NVIDIA graphics devices are character devices with major number 195.
		// https://github.com/torvalds/linux/blob/v4.13/Documentation/admin-guide/devices.txt#L2583
		if fields[0] == "c" && majorMinor[0] == "195" {
			minorNumber, err := strconv.Atoi(majorMinor[1])
			if err != nil {
				return nvidiaMinorNumbers, fmt.Errorf("invalid devices cgroup entry %q: minor number is not integer", text)
			}
			// We don't want devices like nvidiactl (195:255) and nvidia-modeset (195:254)
			if minorNumber < 128 {
				nvidiaMinorNumbers = append(nvidiaMinorNumbers, minorNumber)
			}
			// We are ignoring the "195:*" case
			// where the container has access to all NVIDIA devices on the machine.
		}
		// We are ignoring the "*:*" case
		// where the container has access to all devices on the machine.
	}
	return nvidiaMinorNumbers, nil
}
