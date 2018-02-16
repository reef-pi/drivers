# drivers

Drivers for various ancillary hardware used in reef-pi, based on
[reef-pi/rpi](https://github.com/reef-pi/rpi)

[![Build Status](https://travis-ci.org/reef-pi/drivers.png?branch=master)](https://travis-ci.org/reef-pi/drivers)
[![Coverage Status](https://codecov.io/gh/reef-pi/drivers/branch/master/graph/badge.svg)](https://codecov.io/gh/reef-pi/drivers)
[![Go Report Card](https://goreportcard.com/badge/reef-pi/drivers)](https://goreportcard.com/report/reef-pi/drivers)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/reef-pi/drivers/blob/master/LICENSE.txt)
[![GoDoc](https://godoc.org/github.com/reef-pi/drivers?status.svg)](https://godoc.org/github.com/reef-pi/drivers)

## Introduction

This repository contains a set of drivers for ancillay hardware  used in [reef-pi](http://reef-pi.com) project. They are intended to
have minimal dependencies (most cases only reef-pi/rpi). These drivers API are not stable, and subjected to change as per reef-pi's
requirement.

## Currently available drivers

- PWM: PCA9685
- LED Display: HT16k33
- pH probe: Atlas scientific ezo ph circuit

## License

Copyright:: Copyright (c) 2018 Ranjib Dey.
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
