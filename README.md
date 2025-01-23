i.MX Serial over USB driver
===========================

This Go package implements communication through Serial over USB (CDC-ACM) on
NXP i.MX SoCs, to be used with `GOOS=tamago GOARCH=arm` as supported by the
[TamaGo](https://github.com/usbarmory/tamago) framework for bare metal
Go on ARM SoCs.

This CDC-ACM Serial over USB driver has been tested on Linux as follows:

```
picocom -b 115200 -eb /dev/ttyACM0 --imap lfcrlf
```

Authors
=======

Andrea Barisani  
andrea@inversepath.com  

Documentation
=============

The package API documentation can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/usbarmory/imx-usbserial).


For more information about TamaGo see its
[repository](https://github.com/usbarmory/tamago) and
[project wiki](https://github.com/usbarmory/tamago/wiki).

License
=======

tamago | https://github.com/usbarmory/imx-usbserial  
Copyright (c) WithSecure Corporation

These source files are distributed under the BSD-style license found in the
[LICENSE](https://github.com/usbarmory/imx-usbnet/blob/master/LICENSE) file.
