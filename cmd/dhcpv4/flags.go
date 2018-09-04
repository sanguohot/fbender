/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv4

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/spf13/pflag"
)

type optionCodeSliceValue struct {
	value   *[]dhcpv4.OptionCode
	changed bool
}

// NewOptionCodeSliceValue creates a new option code slice value for pflag.
func NewOptionCodeSliceValue() pflag.Value {
	v := []dhcpv4.OptionCode{}
	return &optionCodeSliceValue{
		value:   &v,
		changed: false,
	}
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func (s *optionCodeSliceValue) Set(value string) error {
	values, err := readAsCSV(value)
	if err != nil {
		return err
	}
	optcodes := []dhcpv4.OptionCode{}
	for _, v := range values {
		optcode, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return err
		}
		optcodes = append(optcodes, dhcpv4.OptionCode(optcode))
	}
	if !s.changed {
		*s.value = optcodes
	} else {
		*s.value = append(*s.value, optcodes...)
	}
	s.changed = true
	return nil
}

func (s *optionCodeSliceValue) Type() string {
	return "optionCodeSlice"
}

func (s *optionCodeSliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, o := range *s.value {
		if name, ok := dhcpv4.OptionCodeToString[o]; ok {
			out[i] = name
		} else {
			out[i] = fmt.Sprintf("Unknown(%d)", o)
		}
	}
	return "[" + strings.Join(out, ",") + "]"
}

// GetOptionCodes returns an option code slice from a pflag set
func GetOptionCodes(f *pflag.FlagSet, name string) ([]dhcpv4.OptionCode, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("flag %s accessed but not defined", name)
	}
	return GetOptionCodesValue(flag.Value)
}

// GetOptionCodesValue returns an option code slice from a pflag value
func GetOptionCodesValue(v pflag.Value) ([]dhcpv4.OptionCode, error) {
	if optcodes, ok := v.(*optionCodeSliceValue); ok {
		return *optcodes.value, nil
	}
	return nil, fmt.Errorf("trying to get option codes value of flag of type %s", v.Type())
}