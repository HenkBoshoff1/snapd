// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
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

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/interfaces/ifacetest"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/snap/snaptest"
)

type utilsSuite struct {
	iface      interfaces.Interface
	slotOS     *interfaces.Slot
	slotApp    *interfaces.Slot
	slotGadget *interfaces.Slot
}

var _ = Suite(&utilsSuite{
	iface:      &ifacetest.TestInterface{InterfaceName: "iface"},
	slotOS:     &interfaces.Slot{SlotInfo: &snap.SlotInfo{Snap: &snap.Info{Type: snap.TypeOS}}},
	slotApp:    &interfaces.Slot{SlotInfo: &snap.SlotInfo{Snap: &snap.Info{Type: snap.TypeApp}}},
	slotGadget: &interfaces.Slot{SlotInfo: &snap.SlotInfo{Snap: &snap.Info{Type: snap.TypeGadget}}},
})

func (s *utilsSuite) TestSanitizeSlotReservedForOS(c *C) {
	errmsg := "iface slots are reserved for the core snap"
	c.Assert(builtin.SanitizeSlotReservedForOS(s.iface, s.slotOS), IsNil)
	c.Assert(builtin.SanitizeSlotReservedForOS(s.iface, s.slotApp), ErrorMatches, errmsg)
	c.Assert(builtin.SanitizeSlotReservedForOS(s.iface, s.slotGadget), ErrorMatches, errmsg)
}

func (s *utilsSuite) TestSanitizeSlotReservedForOSOrGadget(c *C) {
	errmsg := "iface slots are reserved for the core and gadget snaps"
	c.Assert(builtin.SanitizeSlotReservedForOSOrGadget(s.iface, s.slotOS), IsNil)
	c.Assert(builtin.SanitizeSlotReservedForOSOrGadget(s.iface, s.slotApp), ErrorMatches, errmsg)
	c.Assert(builtin.SanitizeSlotReservedForOSOrGadget(s.iface, s.slotGadget), IsNil)
}

func (s *utilsSuite) TestSanitizeSlotReservedForOSOrApp(c *C) {
	errmsg := "iface slots are reserved for the core and app snaps"
	c.Assert(builtin.SanitizeSlotReservedForOSOrApp(s.iface, s.slotOS), IsNil)
	c.Assert(builtin.SanitizeSlotReservedForOSOrApp(s.iface, s.slotApp), IsNil)
	c.Assert(builtin.SanitizeSlotReservedForOSOrApp(s.iface, s.slotGadget), ErrorMatches, errmsg)
}

func MockPlug(c *C, yaml string, si *snap.SideInfo, plugName string) *interfaces.Plug {
	info := snaptest.MockInfo(c, yaml, si)
	if plugInfo, ok := info.Plugs[plugName]; ok {
		return &interfaces.Plug{PlugInfo: plugInfo}
	}
	panic(fmt.Sprintf("cannot find plug %q in snap %q", plugName, info.Name()))
}

func MockSlot(c *C, yaml string, si *snap.SideInfo, slotName string) *interfaces.Slot {
	info := snaptest.MockInfo(c, yaml, si)
	if slotInfo, ok := info.Slots[slotName]; ok {
		return &interfaces.Slot{SlotInfo: slotInfo}
	}
	panic(fmt.Sprintf("cannot find slot %q in snap %q", slotName, info.Name()))
}
