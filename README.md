# drivers

Drivers for various ancillary hardware used in reef-pi, based on
[reef-pi/hal](https://github.com/reef-pi/hal)

[![Build Status](https://github.com/reef-pi/drivers/workflows/go/badge.svg?branch=master)](https://github.com/reef-pi/drivers/actions)
[![Coverage Status](https://codecov.io/gh/reef-pi/drivers/branch/master/graph/badge.svg)](https://codecov.io/gh/reef-pi/drivers)
[![Go Report Card](https://goreportcard.com/badge/reef-pi/drivers)](https://goreportcard.com/report/reef-pi/drivers)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/reef-pi/drivers/blob/master/LICENSE.txt)
[![GoDoc](https://godoc.org/github.com/reef-pi/drivers?status.svg)](https://godoc.org/github.com/reef-pi/drivers)

## Introduction

This repository contains a set of drivers for ancillay hardware used
in [reef-pi](http://reef-pi.com) project. They are intended to have
minimal dependencies (most cases only reef-pi/rpi). These drivers API
are not stable, and subjected to change as per reef-pi's requirement.

## Currently available drivers

- Kasa hs300, hs303, hs103, hs110 smart switches and power strips
- Digital Loggers [web power switch](https://dlidirect.com/products/new-pro-switch)
- Tasmota based smart outlets
- reef-pi open source ph_board: ADS1115 based pH circuits
- PCA9685 PWM driver
- ADS1x15 Analog to digital converter
- Atlas Scientific ezo ph circuit
- Blue acro pico-board: ATSAMD10 pH adapter for the blueAcro Pico board



## License

Copyright:: Copyright (c) 2025 Ranjib Dey.
License:: Apache License, Version 2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
