package config

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/type_and_version"
)

type ContractType string

var (
	EVM2EVMOnRamp  ContractType = "EVM2EVMOnRamp"
	EVM2EVMOffRamp ContractType = "EVM2EVMOffRamp"
	CommitStore    ContractType = "CommitStore"
	PriceRegistry  ContractType = "PriceRegistry"
	Unknown        ContractType = "Unknown" // 1.0.0 Contracts which have no TypeAndVersion
	ContractTypes               = mapset.NewSet[ContractType](
		EVM2EVMOffRamp,
		EVM2EVMOnRamp,
		CommitStore,
		PriceRegistry,
	)
)

func VerifyTypeAndVersion(addr common.Address, client bind.ContractBackend, expectedType ContractType) (semver.Version, error) {
	contractType, version, err := TypeAndVersion(addr, client)
	if err != nil {
		return semver.Version{}, fmt.Errorf("failed getting type and version %w", err)
	}
	if contractType != expectedType {
		return semver.Version{}, fmt.Errorf("wrong contract type %s", contractType)
	}
	return version, nil
}

func TypeAndVersion(addr common.Address, client bind.ContractBackend) (ContractType, semver.Version, error) {
	tv, err := type_and_version.NewITypeAndVersion(addr, client)
	if err != nil {
		return "", semver.Version{}, err
	}
	tvStr, err := tv.TypeAndVersion(nil)
	if err != nil {
		return "", semver.Version{}, fmt.Errorf("error calling typeAndVersion on addr: %s %w", addr.String(), err)
	}

	contractType, versionStr, err := ParseTypeAndVersion(tvStr)
	if err != nil {
		return "", semver.Version{}, err
	}
	v, err := semver.NewVersion(versionStr)
	if err != nil {
		return "", semver.Version{}, fmt.Errorf("failed parsing version %s: %w", versionStr, err)
	}

	if !ContractTypes.Contains(ContractType(contractType)) {
		return "", semver.Version{}, fmt.Errorf("unrecognized contract type %v", contractType)
	}
	return ContractType(contractType), *v, nil
}

// default version to use when TypeAndVersion is missing.
const defaultVersion = "1.0.0"

func ParseTypeAndVersion(tvStr string) (string, string, error) {
	if tvStr == "" {
		tvStr = string(Unknown) + " " + defaultVersion
	}
	typeAndVersionValues := strings.Split(tvStr, " ")

	if len(typeAndVersionValues) < 2 {
		return "", "", fmt.Errorf("invalid type and version %s", tvStr)
	}
	return typeAndVersionValues[0], typeAndVersionValues[1], nil
}
