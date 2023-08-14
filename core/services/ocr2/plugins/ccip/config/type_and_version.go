package config

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	type_and_version "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/type_and_version_interface_wrapper"
)

type ContractType string

var (
	EVM2EVMOnRamp  ContractType = "EVM2EVMOnRamp"
	EVM2EVMOffRamp ContractType = "EVM2EVMOffRamp"
	CommitStore    ContractType = "CommitStore"
	Router         ContractType = "Router"
	ContractTypes               = map[ContractType]struct{}{
		EVM2EVMOffRamp: {},
		EVM2EVMOnRamp:  {},
		CommitStore:    {},
	}
)

func VerifyTypeAndVersion(addr common.Address, client bind.ContractBackend, expectedType ContractType) error {
	contractType, _, err := typeAndVersion(addr, client)
	if err != nil {
		return errors.Errorf("failed getting type and version %v", err)
	}
	if contractType != expectedType {
		return errors.Errorf("Wrong contract type %s", contractType)
	}
	return nil
}

func typeAndVersion(addr common.Address, client bind.ContractBackend) (ContractType, semver.Version, error) {
	tv, err := type_and_version.NewTypeAndVersionInterface(addr, client)
	if err != nil {
		return "", semver.Version{}, errors.Wrap(err, "failed creating a type and version")
	}
	tvStr, err := tv.TypeAndVersion(nil)
	if err != nil {
		return "", semver.Version{}, errors.Wrap(err, "failed to call type and version")
	}
	typeAndVersionValues := strings.Split(tvStr, " ")

	if len(typeAndVersionValues) < 2 {
		return "", semver.Version{}, fmt.Errorf("invalid type and version %s", tvStr)
	}
	contractType, version := typeAndVersionValues[0], typeAndVersionValues[1]
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", semver.Version{}, err
	}
	if _, ok := ContractTypes[ContractType(contractType)]; !ok {
		return "", semver.Version{}, errors.Errorf("unrecognized contract type %v", contractType)
	}
	return ContractType(contractType), *v, nil
}
