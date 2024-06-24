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

package builtin_test

import (
	"fmt"

	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/dirs"
	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/apparmor"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/interfaces/udev"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/testutil"
)

type CpuLatencyQoSObserveInterfaceSuite struct {
	iface    interfaces.Interface
	slotInfo *snap.SlotInfo
	slot     *interfaces.ConnectedSlot
	plugInfo *snap.PlugInfo
	plug     *interfaces.ConnectedPlug
}

var _ = Suite(&CpuLatencyQoSObserveInterfaceSuite{
	iface: builtin.MustInterface("cpu-latency-qos-observe"),
})

const cpuLatencyQoSObserveConsumerYaml = `name: consumer
version: 0
apps:
 app:
  plugs: [cpu-latency-qos-observe]
`

const cpuLatencyQoSObserveCoreYaml = `name: core
version: 0
type: os
slots:
  cpu-latency-qos-observe:
`

func (s *CpuLatencyQoSObserveInterfaceSuite) SetUpTest(c *C) {
	s.plug, s.plugInfo = MockConnectedPlug(c, cpuLatencyQoSObserveConsumerYaml, nil, "cpu-latency-qos-observe")
	s.slot, s.slotInfo = MockConnectedSlot(c, cpuLatencyQoSObserveCoreYaml, nil, "cpu-latency-qos-observe")
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "cpu-latency-qos-observe")
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestSanitizeSlot(c *C) {
	c.Assert(interfaces.BeforePrepareSlot(s.iface, s.slotInfo), IsNil)
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestAppArmorSpec(c *C) {
	appSet, err := interfaces.NewSnapAppSet(s.plug.Snap(), nil)
	c.Assert(err, IsNil)
	spec := apparmor.NewSpecification(appSet)
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, s.slot), IsNil)
	c.Assert(spec.SecurityTags(), DeepEquals, []string{"snap.consumer.app"})
	c.Assert(spec.SnippetForTag("snap.consumer.app"), testutil.Contains, "# Description: Allow read access to the device node cpu_dma_latency,")
	c.Assert(spec.SnippetForTag("snap.consumer.app"), testutil.Contains, "# responsible for controlling the CPU latency QoS from userspace.")
	c.Assert(spec.SnippetForTag("snap.consumer.app"), testutil.Contains, "/dev/cpu_dma_latency r,")
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestUDevSpec(c *C) {
	appSet, err := interfaces.NewSnapAppSet(s.plug.Snap(), nil)
	c.Assert(err, IsNil)
	spec := udev.NewSpecification(appSet)
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, s.slot), IsNil)
	c.Assert(spec.Snippets(), HasLen, 2)
	c.Assert(spec.Snippets()[0], Equals, `# cpu-latency-qos-observe
SUBSYSTEM=="misc", KERNEL=="cpu_dma_latency", TAG+="snap_consumer_app"`)
	c.Assert(spec.Snippets(), testutil.Contains, fmt.Sprintf(`TAG=="snap_consumer_app", SUBSYSTEM!="module", SUBSYSTEM!="subsystem", RUN+="%v/snap-device-helper $env{ACTION} snap_consumer_app $devpath $major:$minor"`, dirs.DistroLibExecDir))
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestStaticInfo(c *C) {
	si := interfaces.StaticInfoOf(s.iface)
	c.Assert(si.ImplicitOnCore, Equals, true)
	c.Assert(si.ImplicitOnClassic, Equals, true)
	c.Assert(si.Summary, Equals, `allow read access to cpu_dma_latency device`)
	c.Assert(si.BaseDeclarationSlots, testutil.Contains, "cpu-latency-qos-observe")
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestAutoConnect(c *C) {
	c.Assert(s.iface.AutoConnect(s.plugInfo, s.slotInfo), Equals, true)
}

func (s *CpuLatencyQoSObserveInterfaceSuite) TestInterfaces(c *C) {
	c.Check(builtin.Interfaces(), testutil.DeepContains, s.iface)
}
