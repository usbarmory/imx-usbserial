// Serial over USB driver
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package usbserial

import (
	"github.com/usbarmory/tamago/soc/nxp/usb"
)

// MaxPacketSize represents the USB data interface endpoint maximum packet size
var MaxPacketSize uint16 = 512

func addControlInterface(device *usb.Device, serial *UART) (iface *usb.InterfaceDescriptor) {
	iface = &usb.InterfaceDescriptor{}
	iface.SetDefaults()

	iface.NumEndpoints = 1
	iface.InterfaceClass = usb.COMMUNICATION_INTERFACE_CLASS
	iface.InterfaceSubClass = usb.ACM_SUBCLASS
	iface.InterfaceProtocol = usb.AT_COMMAND_PROTOCOL

	iInterface, _ := device.AddString(`CDC Abstract Control Model (ACM) UART`)
	iface.Interface = iInterface

	// Set IAD to be inserted before first interface, to support multiple
	// functions in this same configuration.
	iface.IAD = &usb.InterfaceAssociationDescriptor{}
	iface.IAD.SetDefaults()
	// alternate settings do not count
	iface.IAD.InterfaceCount = 1
	iface.IAD.FunctionClass = iface.InterfaceClass
	iface.IAD.FunctionSubClass = iface.InterfaceSubClass

	iFunction, _ := device.AddString(`CDC`)
	iface.IAD.Function = iFunction

	header := &usb.CDCHeaderDescriptor{}
	header.SetDefaults()

	iface.ClassDescriptors = append(iface.ClassDescriptors, header.Bytes())

	acm := &usb.CDCAbstractControlManagementDescriptor{}
	acm.SetDefaults()

	iface.ClassDescriptors = append(iface.ClassDescriptors, acm.Bytes())

	union := &usb.CDCUnionDescriptor{}
	union.SetDefaults()

	numInterfaces := 1 + len(device.Configurations[0].Interfaces)
	union.MasterInterface = uint8(numInterfaces - 1)
	union.SlaveInterface0 = uint8(numInterfaces)

	iface.ClassDescriptors = append(iface.ClassDescriptors, union.Bytes())

	cm := &usb.CDCCallManagementDescriptor{}
	cm.SetDefaults()

	iface.ClassDescriptors = append(iface.ClassDescriptors, cm.Bytes())

	ep1IN := &usb.EndpointDescriptor{}
	ep1IN.SetDefaults()
	ep1IN.EndpointAddress = 0x81
	ep1IN.Attributes = 3
	ep1IN.MaxPacketSize = 64
	ep1IN.Interval = 11
	ep1IN.Function = serial.Control

	iface.Endpoints = append(iface.Endpoints, ep1IN)

	device.Configurations[0].AddInterface(iface)

	return
}

func addDataInterfaces(device *usb.Device, serial *UART) {
	iface1 := &usb.InterfaceDescriptor{}
	iface1.SetDefaults()

	iface1.NumEndpoints = 2
	iface1.InterfaceClass = usb.DATA_INTERFACE_CLASS

	iInterface, _ := device.AddString(`CDC Data`)
	iface1.Interface = iInterface

	ep2IN := &usb.EndpointDescriptor{}
	ep2IN.SetDefaults()
	ep2IN.EndpointAddress = 0x82
	ep2IN.Attributes = 2
	ep2IN.MaxPacketSize = MaxPacketSize
	ep2IN.Function = serial.Tx

	iface1.Endpoints = append(iface1.Endpoints, ep2IN)

	ep2OUT := &usb.EndpointDescriptor{}
	ep2OUT.SetDefaults()
	ep2OUT.EndpointAddress = 0x02
	ep2IN.MaxPacketSize = MaxPacketSize
	ep2OUT.Attributes = 2
	ep2OUT.Function = serial.Rx

	iface1.Endpoints = append(iface1.Endpoints, ep2OUT)

	device.Configurations[0].AddInterface(iface1)

	return
}