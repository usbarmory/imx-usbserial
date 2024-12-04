// Serial over USB driver
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

// Package usbserial implements Serial over USB (CDC-ACM) on i.MX6 SoCs.
//
// This package is only meant to be used with `GOOS=tamago GOARCH=arm` as
// supported by the TamaGo framework for bare metal Go, see
// https://github.com/usbarmory/tamago.
package usbserial

import (
	"bytes"
	"sync"

	"github.com/usbarmory/tamago/soc/nxp/usb"
)

// UART represents a Serial over USB interface instance.
type UART struct {
	sync.Mutex

	// Device is the USB device associated to Serial over USB interface.
	Device *usb.Device

	// Rx is endpoint 1 OUT function, set by Init() to ACMRx if not
	// already defined.
	Rx func([]byte, error) ([]byte, error)

	// Tx is endpoint 1 IN function, set by Init() to ACMTx if not already
	// defined.
	Tx func([]byte, error) ([]byte, error)

	// Control is endpoint 2 IN function
	Control func([]byte, error) ([]byte, error)

	buf bytes.Buffer
}

// Init initializes a Serial over USB interface.
func (serial *UART) Init() (err error) {
	if serial.Device == nil {
		serial.Device = &usb.Device{}
	}

	if serial.Rx == nil {
		serial.Rx = serial.ACMRx
	}

	if serial.Tx == nil {
		serial.Tx = serial.ACMTx
	}

	if serial.Control == nil {
		serial.Control = serial.ACMControl
	}

	addControlInterface(serial.Device, serial)
	addDataInterfaces(serial.Device, serial)

	return
}

// ACMControl implements the endpoint 1 IN function.
func (serial *UART) ACMControl(_ []byte, lastErr error) (in []byte, err error) {
	// ignore for now
	return
}

// ACMRx implements the endpoint 2 OUT function, used to receive serial
// communication from host to device.
func (serial *UART) ACMRx(out []byte, lastErr error) (_ []byte, err error) {
	// disabled by default
	return
}

// ACMTx implements the endpoint 2 IN function, used to transmit serial
// communication from device to host.
func (serial *UART) ACMTx(_ []byte, lastErr error) (in []byte, err error) {
	serial.Lock()
	defer serial.Unlock()

	in = serial.buf.Bytes()
	serial.buf.Reset()

	return
}

// Write appends the contents of p to the serial transmission buffer.
func (serial *UART) Write(p []byte) (n int, err error) {
	serial.Lock()
	defer serial.Unlock()

	return serial.buf.Write(p)
}

// WriteByte appends the byte c to the serial transmission buffer.
func (serial *UART) WriteByte(c byte) (err error) {
	serial.Lock()
	defer serial.Unlock()

	return serial.buf.WriteByte(c)
}
