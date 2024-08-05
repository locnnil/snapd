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

const ptracePoCSummary = `PoC of seccomp bypass using ptrace`

const ptracePoCBaseDeclarationSlots = `
  ptrace-poc:
    allow-installation:
      slot-snap-type:
        - core
    deny-auto-connection: true
`

const ptracePoCConnectedPlugAppArmor = `
# Description: Allow tracing our own processes.
# Note, this allows seccomp sandbox escape on kernels < 4.8
capability sys_ptrace,
`

const ptracePoCConnectedPlugSecComp = `
# Description: PoC of seccomp bypass using ptracePoC.
ptrace
`

func init() {
	registerIface(&commonInterface{
		name:                  "ptrace-poc",
		summary:               ptracePoCSummary,
		implicitOnCore:        true,
		implicitOnClassic:     true,
		baseDeclarationSlots:  ptracePoCBaseDeclarationSlots,
		connectedPlugAppArmor: ptracePoCConnectedPlugAppArmor,
		connectedPlugSecComp:  ptracePoCConnectedPlugSecComp,
	})
}
