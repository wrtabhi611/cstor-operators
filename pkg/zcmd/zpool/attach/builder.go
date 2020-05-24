/*
Copyright 2019 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pattach

import (
	"fmt"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/openebs/cstor-operators/pkg/zcmd/bin"
	"github.com/pkg/errors"
)

const (
	// Operation defines type of zfs operation
	Operation = "attach"
)

//PoolAttach defines structure for pool 'Attach' operation
type PoolAttach struct {
	//list of property
	Property []string

	//forcefully attach
	Forcefully bool

	//device name
	Device string

	//new device name
	NewDevice string

	//pool name
	Pool string

	// command string
	Command string

	// checks is list of predicate function used for validating object
	checks []PredicateFunc

	// error
	err error

	// Executor is to execute the commands
	Executor bin.Executor
}

// NewPoolAttach returns new instance of object PoolAttach
func NewPoolAttach() *PoolAttach {
	return &PoolAttach{}
}

// WithCheck add given check to checks list
func (p *PoolAttach) WithCheck(check ...PredicateFunc) *PoolAttach {
	p.checks = append(p.checks, check...)
	return p
}

// WithProperty method fills the Property field of PoolAttach object.
func (p *PoolAttach) WithProperty(key, value string) *PoolAttach {
	p.Property = append(p.Property, fmt.Sprintf("%s=%s", key, value))
	return p
}

// WithForcefully method fills the Forcefully field of PoolAttach object.
func (p *PoolAttach) WithForcefully(Forcefully bool) *PoolAttach {
	p.Forcefully = Forcefully
	return p
}

// WithDevice method fills the Device field of PoolAttach object.
func (p *PoolAttach) WithDevice(Device string) *PoolAttach {
	p.Device = Device
	return p
}

// WithNewDevice method fills the NewDevice field of PoolAttach object.
func (p *PoolAttach) WithNewDevice(NewDevice string) *PoolAttach {
	p.NewDevice = NewDevice
	return p
}

// WithPool method fills the Pool field of PoolAttach object.
func (p *PoolAttach) WithPool(Pool string) *PoolAttach {
	p.Pool = Pool
	return p
}

// WithCommand method fills the Command field of PoolAttach object.
func (p *PoolAttach) WithCommand(Command string) *PoolAttach {
	p.Command = Command
	return p
}

// WithExecutor method fills the Executor field of PoolDump object.
func (p *PoolAttach) WithExecutor(executor bin.Executor) *PoolAttach {
	p.Executor = executor
	return p
}

// Validate is to validate generated PoolAttach object by builder
func (p *PoolAttach) Validate() *PoolAttach {
	for _, check := range p.checks {
		if !check(p) {
			p.err = errors.Wrapf(p.err, "validation failed {%v}", runtime.FuncForPC(reflect.ValueOf(check).Pointer()).Name())
		}
	}
	return p
}

// Execute is to execute generated PoolAttach object
func (p *PoolAttach) Execute() ([]byte, error) {
	p, err := p.Build()
	if err != nil {
		return nil, err
	}

	if IsExecutorSet()(p) {
		return p.Executor.Execute(p.Command)
	}

	// execute command here
	// #nosec
	return exec.Command(bin.BASH, "-c", p.Command).CombinedOutput()
}

// Build returns the PoolAttach object generated by builder
func (p *PoolAttach) Build() (*PoolAttach, error) {
	var c strings.Builder
	p = p.Validate()
	p.appendCommand(&c, bin.ZPOOL)
	p.appendCommand(&c, fmt.Sprintf(" %s ", Operation))

	if IsForcefullySet()(p) {
		p.appendCommand(&c, fmt.Sprintf(" -f "))
	}

	if IsPropertySet()(p) {
		for _, v := range p.Property {
			p.appendCommand(&c, fmt.Sprintf(" -o %s ", v))
		}
	}

	p.appendCommand(&c, p.Pool)

	p.appendCommand(&c, fmt.Sprintf(" %s ", p.Device))
	p.appendCommand(&c, fmt.Sprintf(" %s ", p.NewDevice))

	p.Command = c.String()
	return p, p.err
}

// appendCommand append string to given string builder
func (p *PoolAttach) appendCommand(c *strings.Builder, cmd string) {
	_, err := c.WriteString(cmd)
	if err != nil {
		p.err = errors.Wrapf(p.err, "Failed to append cmd{%s} : %s", cmd, err.Error())
	}
}
