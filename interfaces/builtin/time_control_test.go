// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
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
	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/apparmor"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/interfaces/udev"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/testutil"
)

type TimeControlInterfaceSuite struct {
	iface interfaces.Interface
	slot  *interfaces.Slot
	plug  *interfaces.Plug
}

var _ = Suite(&TimeControlInterfaceSuite{
	iface: builtin.MustInterface("time-control"),
	slot: &interfaces.Slot{
		SlotInfo: &snap.SlotInfo{
			Snap:      &snap.Info{SuggestedName: "core", Type: snap.TypeOS},
			Name:      "time-control",
			Interface: "time-control",
		},
	},
	plug: nil,
})

const timectlConsumerYaml = `name: consumer
apps:
 app:
  plugs: [time-control]
`

const timectlCoreYaml = `name: core
type: os
slots:
  time-control:
`

func (s *TimeControlInterfaceSuite) SetUpTest(c *C) {
	s.plug = MockPlug(c, timectlConsumerYaml, nil, "time-control")
	s.slot = MockSlot(c, timectlCoreYaml, nil, "time-control")
}

func (s *TimeControlInterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "time-control")
}

func (s *TimeControlInterfaceSuite) TestSanitizeSlot(c *C) {
	c.Assert(s.slot.Sanitize(s.iface), IsNil)
	slot := &interfaces.Slot{SlotInfo: &snap.SlotInfo{
		Snap:      &snap.Info{SuggestedName: "some-snap"},
		Name:      "time-control",
		Interface: "time-control",
	}}
	c.Assert(slot.Sanitize(s.iface), ErrorMatches,
		"time-control slots are reserved for the core snap")
}

func (s *TimeControlInterfaceSuite) TestSanitizePlug(c *C) {
	c.Assert(s.plug.Sanitize(s.iface), IsNil)
}

func (s *TimeControlInterfaceSuite) TestAppArmorSpec(c *C) {
	spec := &apparmor.Specification{}
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, nil, s.slot, nil), IsNil)
	c.Assert(spec.SecurityTags(), DeepEquals, []string{"snap.consumer.app"})
	c.Check(spec.SnippetForTag("snap.consumer.app"), testutil.Contains, "org/freedesktop/timedate1")
}

func (s *TimeControlInterfaceSuite) TestUDevSpec(c *C) {
	spec := &udev.Specification{}
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, nil, s.slot, nil), IsNil)
	c.Assert(spec.Snippets(), HasLen, 1)
	c.Assert(spec.Snippets()[0], testutil.Contains, `KERNEL=="/dev/rtc0", TAG+="snap_consumer_app"`)
}

func (s *TimeControlInterfaceSuite) TestStaticInfo(c *C) {
	si := interfaces.StaticInfoOf(s.iface)
	c.Assert(si.ImplicitOnCore, Equals, true)
	c.Assert(si.ImplicitOnClassic, Equals, true)
	c.Assert(si.Summary, Equals, `allows setting system date and time`)
	c.Assert(si.BaseDeclarationSlots, testutil.Contains, "time-control")
}

func (s *TimeControlInterfaceSuite) TestAutoConnect(c *C) {
	c.Assert(s.iface.AutoConnect(s.plug, s.slot), Equals, true)
}

func (s *TimeControlInterfaceSuite) TestInterfaces(c *C) {
	c.Check(builtin.Interfaces(), testutil.DeepContains, s.iface)
}
