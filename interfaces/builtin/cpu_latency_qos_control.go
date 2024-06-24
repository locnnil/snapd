// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2024 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin

const cpuLatencyQoSControlSummary = `allow write access to cpu_dma_latency device`

const cpuLatencyQoSControlBaseDeclarationSlots = `
  cpu-latency-qos-control:
    allow-installation:
      slot-snap-type:
        - core
    deny-auto-connection: true
`
const cpuLatencyQoSControlConnectedPlugAppArmor = `
# Description: Allow write access to the device node cpu_dma_latency,
# responsible for controlling the CPU latency QoS from userspace.

/dev/cpu_dma_latency rw,
`

var cpuLatencyQoSControlConnectedPlugUDev = []string{
	`SUBSYSTEM=="misc", KERNEL=="cpu_dma_latency"`,
}

func init() {
	registerIface(&commonInterface{
		name:                  "cpu-latency-qos-control",
		summary:               cpuLatencyQoSControlSummary,
		implicitOnCore:        true,
		implicitOnClassic:     true,
		baseDeclarationSlots:  cpuLatencyQoSControlBaseDeclarationSlots,
		connectedPlugAppArmor: cpuLatencyQoSControlConnectedPlugAppArmor,
		connectedPlugUDev:     cpuLatencyQoSControlConnectedPlugUDev,
	})
}
