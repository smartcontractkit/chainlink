// Copyright 2021 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/alecthomas/units"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (in *PhysicalMachineChaosSpec) Default(root interface{}, field *reflect.StructField) {
	if in == nil {
		return
	}

	if len(in.UID) == 0 {
		in.UID = uuid.New().String()
	}

	for i := range in.Address {
		// add http prefix for address
		if !strings.HasPrefix(in.Address[i], "http") {
			in.Address[i] = fmt.Sprintf("http://%s", in.Address[i])
		}
	}
}

func (in *PhysicalMachineChaosSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// make sure the configuration corresponding to action is not empty
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(in)
	if err != nil {
		allErrs = append(allErrs,
			field.Invalid(path.Child("spec"), in, err.Error()))
	}

	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		allErrs = append(allErrs,
			field.Invalid(path.Child("spec"), in, err.Error()))
	}

	skipConfigCheck := false
	if _, ok := inInterface[string(in.Action)]; !ok {
		skipConfigCheck = true
		allErrs = append(allErrs,
			field.Invalid(path.Child("spec"), in,
				"the configuration corresponding to action is required"))
	}

	if len(in.Address) == 0 && in.Selector.Empty() {
		allErrs = append(allErrs,
			field.Invalid(path.Child("address"), in.Address, "one of address or selector should be specified"))
	}
	if len(in.Address) != 0 && !in.Selector.Empty() {
		allErrs = append(allErrs,
			field.Invalid(path.Child("address"), in.Address, "only one of address or selector could be specified"))
	}
	// make sure address is not empty
	for _, address := range in.Address {
		if len(address) == 0 {
			allErrs = append(allErrs,
				field.Invalid(path.Child("address"), in.Address, "the address is required"))
		}
	}

	if skipConfigCheck {
		return allErrs
	}

	var validateConfigErr error
	switch in.Action {
	case PMStressCPUAction:
		validateConfigErr = validateStressCPUAction(in.StressCPU)
	case PMStressMemAction:
		validateConfigErr = validateStressMemAction(in.StressMemory)
	case PMDiskWritePayloadAction:
		validateConfigErr = validateDiskPayloadAction(in.DiskWritePayload)
	case PMDiskReadPayloadAction:
		validateConfigErr = validateDiskPayloadAction(in.DiskReadPayload)
	case PMDiskFillAction:
		validateConfigErr = validateDiskFillAction(in.DiskFill)
	case PMNetworkCorruptAction:
		validateConfigErr = validateNetworkCorruptAction(in.NetworkCorrupt)
	case PMNetworkDuplicateAction:
		validateConfigErr = validateNetworkDuplicateAction(in.NetworkDuplicate)
	case PMNetworkLossAction:
		validateConfigErr = validateNetworkLossAction(in.NetworkLoss)
	case PMNetworkDelayAction:
		validateConfigErr = validateNetworkDelayAction(in.NetworkDelay)
	case PMNetworkPartitionAction:
		validateConfigErr = validateNetworkPartitionAction(in.NetworkPartition)
	case PMNetworkBandwidthAction:
		validateConfigErr = validateNetworkBandwidthAction(in.NetworkBandwidth)
	case PMNetworkDNSAction:
		validateConfigErr = validateNetworkDNSAction(in.NetworkDNS)
	case PMProcessAction:
		validateConfigErr = validateProcessAction(in.Process)
	case PMJVMExceptionAction:
		validateConfigErr = validateJVMExceptionAction(in.JVMException)
	case PMJVMGCAction:
		validateConfigErr = validateJVMGCAction(in.JVMGC)
	case PMJVMLatencyAction:
		validateConfigErr = validateJVMLatencyAction(in.JVMLatency)
	case PMJVMReturnAction:
		validateConfigErr = validateJVMReturnAction(in.JVMReturn)
	case PMJVMStressAction:
		validateConfigErr = validateJVMStressAction(in.JVMStress)
	case PMJVMRuleDataAction:
		validateConfigErr = validateJVMRuleDataAction(in.JVMRuleData)
	case PMClockAction:
		validateConfigErr = validateClockAction(in.Clock)
	default:
	}

	if validateConfigErr != nil {
		allErrs = append(allErrs,
			field.Invalid(path.Child("spec"), in,
				validateConfigErr.Error()))
	}

	return allErrs
}

func validateStressCPUAction(spec *StressCPUSpec) error {
	if spec.Load == 0 {
		return errors.New("load can't be 0")
	}

	if spec.Workers == 0 {
		return errors.New("workers can't be 0")
	}

	return nil
}

func validateStressMemAction(spec *StressMemorySpec) error {
	if len(spec.Size) == 0 {
		return errors.New("size is required")
	}

	if _, err := ParseUnit(spec.Size); err != nil {
		return err
	}

	return nil
}

func validateDiskPayloadAction(spec *DiskPayloadSpec) error {
	if spec.PayloadProcessNum == 0 {
		return errors.New("payload-process-num can't be 0")
	}

	if len(spec.Size) == 0 {
		return errors.New("size is required")
	}

	if _, err := ParseUnit(spec.Size); err != nil {
		return err
	}

	return nil
}

func validateDiskFillAction(spec *DiskFillSpec) error {
	if len(spec.Size) == 0 {
		return errors.New("size is required")
	}

	if _, err := ParseUnit(spec.Size); err != nil {
		return err
	}

	return nil
}

func validateNetworkCommon(spec *NetworkCommonSpec) error {
	if !CheckPercent(spec.Correlation, true) {
		return errors.Errorf("correlation %s is invalid", spec.Correlation)
	}

	if len(spec.Device) == 0 {
		return errors.New("device is required")
	}

	if len(spec.IPAddress) == 0 && len(spec.Hostname) == 0 {
		return errors.New("one of ip-address and hostname is required")
	}

	return nil
}

func validateNetworkCorruptAction(spec *NetworkCorruptSpec) error {
	if err := validateNetworkCommon(&spec.NetworkCommonSpec); err != nil {
		return err
	}

	if !CheckPercent(spec.Percent, false) {
		return errors.New("percent is invalid")
	}

	return nil
}

func validateNetworkDuplicateAction(spec *NetworkDuplicateSpec) error {
	if err := validateNetworkCommon(&spec.NetworkCommonSpec); err != nil {
		return err
	}

	if !CheckPercent(spec.Percent, false) {
		return errors.New("percent is invalid")
	}

	return nil
}

func validateNetworkLossAction(spec *NetworkLossSpec) error {
	if err := validateNetworkCommon(&spec.NetworkCommonSpec); err != nil {
		return err
	}

	if !CheckPercent(spec.Percent, false) {
		return errors.New("percent is invalid")
	}

	return nil
}

func validateNetworkDelayAction(spec *NetworkDelaySpec) error {
	if err := validateNetworkCommon(&spec.NetworkCommonSpec); err != nil {
		return err
	}

	if len(spec.Latency) == 0 {
		return errors.New("latency is invalid")
	}

	return nil
}

func validateNetworkPartitionAction(spec *NetworkPartitionSpec) error {
	if len(spec.Device) == 0 {
		return errors.New("device is required")
	}

	if len(spec.IPAddress) == 0 && len(spec.Hostname) == 0 {
		return errors.New("one of ip-address and hostname is required")
	}

	if spec.Direction != "to" && spec.Direction != "from" {
		return errors.New("direction should be one of 'to' and 'from'")
	}

	if len(spec.AcceptTCPFlags) > 0 && spec.IPProtocol != "tcp" {
		return errors.New("protocol should be 'tcp' when set accept-tcp-flags")
	}

	return nil
}

func validateNetworkBandwidthAction(spec *NetworkBandwidthSpec) error {
	if len(spec.Device) == 0 {
		return errors.New("device is required")
	}

	if len(spec.Rate) == 0 || spec.Limit == 0 || spec.Buffer == 0 {
		return errors.Errorf("rate, limit and buffer both are required when action is bandwidth")
	}

	return nil
}

func validateNetworkDNSAction(spec *NetworkDNSSpec) error {
	if (len(spec.DNSDomainName) != 0 && len(spec.DNSIp) == 0) || (len(spec.DNSDomainName) == 0 && len(spec.DNSIp) != 0) {
		return errors.Errorf("DNS host %s must match a DNS ip %s", spec.DNSDomainName, spec.DNSIp)
	}

	return nil
}

func validateProcessAction(spec *ProcessSpec) error {
	if len(spec.Process) == 0 {
		return errors.New("process is required")
	}

	if spec.Signal == 0 {
		return errors.New("signal is required")
	}

	return nil
}

func validateJVMClassMethod(spec *JVMClassMethodSpec) error {
	if len(spec.Class) == 0 {
		return errors.New("class is required")
	}

	if len(spec.Method) == 0 {
		return errors.New("method is required")
	}

	return nil
}

func validateJVMExceptionAction(spec *JVMExceptionSpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if err := validateJVMClassMethod(&spec.JVMClassMethodSpec); err != nil {
		return err
	}

	if len(spec.ThrowException) == 0 {
		return errors.New("exception is required")
	}

	return nil
}

func validateJVMGCAction(spec *JVMGCSpec) error {
	return CheckPid(spec.Pid)
}

func validateJVMLatencyAction(spec *JVMLatencySpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if err := validateJVMClassMethod(&spec.JVMClassMethodSpec); err != nil {
		return err
	}

	if spec.LatencyDuration == 0 {
		return errors.New("latency is required")
	}

	return nil
}

func validateJVMReturnAction(spec *JVMReturnSpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if err := validateJVMClassMethod(&spec.JVMClassMethodSpec); err != nil {
		return err
	}

	if len(spec.ReturnValue) == 0 {
		return errors.New("value is required")
	}

	return nil
}

func validateJVMStressAction(spec *JVMStressSpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if spec.CPUCount == 0 && len(spec.MemoryType) == 0 {
		return errors.New("one of cpu-count and mem-type is required")
	}

	if spec.CPUCount > 0 && len(spec.MemoryType) > 0 {
		return errors.New("inject stress on both CPU and memory is not support")
	}

	return nil
}

func validateJVMRuleDataAction(spec *JVMRuleDataSpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if len(spec.RuleData) == 0 {
		return errors.New("rule-data is required")
	}

	return nil
}

func validateClockAction(spec *ClockSpec) error {
	if err := CheckPid(spec.Pid); err != nil {
		return err
	}

	if len(spec.TimeOffset) == 0 {
		return errors.New("time-offset is required")
	}

	return nil
}

func CheckPid(pid int) error {
	if pid == 0 {
		return errors.New("pid is required")
	}

	if pid < 0 {
		return errors.New("pid is invalid")
	}

	return nil
}

func CheckPercent(p string, allowZero bool) bool {
	if len(p) == 0 {
		return allowZero
	}

	v, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return false
	}

	if v == 0 && !allowZero {
		return false
	}

	if v < 0 || v > 100 {
		return false
	}

	return true
}

var (
	// See https://en.wikipedia.org/wiki/Binary_prefix
	shortBinaryUnitMap = units.MakeUnitMap("", "c", 1024)
	binaryUnitMap      = units.MakeUnitMap("iB", "c", 1024)
	decimalUnitMap     = units.MakeUnitMap("B", "c", 1000)
)

// ParseUnit parse a digit with unit such as "K" , "KiB", "KB", "c", "MiB", "MB", "M".
// If input string is a digit without unit ,
// it will be regarded as a digit with unit M(1024*1024 bytes).
func ParseUnit(s string) (uint64, error) {
	if _, err := strconv.Atoi(s); err == nil {
		s += "B"
	}
	if n, err := units.ParseUnit(s, shortBinaryUnitMap); err == nil {
		return uint64(n), nil
	}

	if n, err := units.ParseUnit(s, binaryUnitMap); err == nil {
		return uint64(n), nil
	}

	if n, err := units.ParseUnit(s, decimalUnitMap); err == nil {
		return uint64(n), nil
	}
	return 0, errors.Wrapf(errInvalidValue, "unknown unit %s", s)
}

func (in *NetworkBandwidthSpec) Validate(root interface{}, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(in.Rate) == 0 {
		allErrs = append(allErrs,
			field.Invalid(path.Child("rate"), in.Rate, "rate is required"))
	}

	return allErrs
}
