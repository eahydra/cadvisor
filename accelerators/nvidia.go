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
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/stats"

	"k8s.io/klog/v2"
)

type nvidiaManager struct {
}

var sysFsPCIDevicesPath = "/sys/bus/pci/devices/"

const nvidiaVendorID = "0x10de"

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
