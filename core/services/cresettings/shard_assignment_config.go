package cresettings

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

type ShardAssignmentConfig struct {
	StaticDefaultAssignment []uint32
	DisabledShards          []uint32
	PerOwnerAssignment      map[string][]uint32
	PerOrgAssignment        map[string][]uint32
	HashedDefaultAssignment bool
	HashedOwnerAssignment   map[string]bool
}

const (
	ConfigTypeSettings        = "settings"
	ConfigTypeShardAssignment = "shard_assignment"
)

func ParseShardAssignmentConfig(raw string) (*ShardAssignmentConfig, error) {
	cfg := &ShardAssignmentConfig{
		PerOwnerAssignment:    make(map[string][]uint32),
		PerOrgAssignment:      make(map[string][]uint32),
		HashedOwnerAssignment: make(map[string]bool),
	}

	if strings.TrimSpace(raw) == "" {
		return cfg, nil
	}

	tree, err := toml.Load(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shard_assignment TOML: %w", err)
	}

	if v := tree.Get("static_default_assignment"); v != nil {
		vals, ok := v.([]any)
		if !ok {
			return nil, errors.New("static_default_assignment must be an array of integers")
		}
		for _, iv := range vals {
			shardID, err := toShardID(iv)
			if err != nil {
				return nil, fmt.Errorf("invalid shard ID in static_default_assignment: %w", err)
			}
			cfg.StaticDefaultAssignment = append(cfg.StaticDefaultAssignment, shardID)
		}
	}

	if v := tree.Get("disabled_shards"); v != nil {
		vals, ok := v.([]any)
		if !ok {
			return nil, errors.New("disabled_shards must be an array of integers")
		}
		for _, iv := range vals {
			shardID, err := toShardID(iv)
			if err != nil {
				return nil, fmt.Errorf("invalid shard ID in disabled_shards: %w", err)
			}
			cfg.DisabledShards = append(cfg.DisabledShards, shardID)
		}
	}

	if v := tree.Get("hashed_default_assignment"); v != nil {
		b, ok := v.(bool)
		if !ok {
			return nil, errors.New("hashed_default_assignment must be a boolean")
		}
		cfg.HashedDefaultAssignment = b
	}

	if perOwner := tree.Get("per_owner_assignment"); perOwner != nil {
		tbl, ok := perOwner.(*toml.Tree)
		if !ok {
			return nil, errors.New("per_owner_assignment must be a table")
		}
		for _, key := range tbl.Keys() {
			normalized, err := normalizeOwnerHex(key)
			if err != nil {
				return nil, fmt.Errorf("invalid owner address %q in per_owner_assignment: %w", key, err)
			}
			v := tbl.Get(key)
			vals, ok := v.([]any)
			if !ok {
				return nil, fmt.Errorf("per_owner_assignment[%q] must be an array of integers", key)
			}
			var shards []uint32
			for _, iv := range vals {
				shardID, err := toShardID(iv)
				if err != nil {
					return nil, fmt.Errorf("invalid shard ID in per_owner_assignment[%q]: %w", key, err)
				}
				shards = append(shards, shardID)
			}
			if len(shards) == 0 {
				return nil, fmt.Errorf("per_owner_assignment[%q] must not be empty", key)
			}
			cfg.PerOwnerAssignment[normalized] = shards
		}
	}

	if perOrg := tree.Get("per_org_assignment"); perOrg != nil {
		tbl, ok := perOrg.(*toml.Tree)
		if !ok {
			return nil, errors.New("per_org_assignment must be a table")
		}
		for _, key := range tbl.Keys() {
			if strings.TrimSpace(key) == "" {
				return nil, errors.New("per_org_assignment keys must not be empty")
			}
			v := tbl.Get(key)
			vals, ok := v.([]any)
			if !ok {
				return nil, fmt.Errorf("per_org_assignment[%q] must be an array of integers", key)
			}
			var shards []uint32
			for _, iv := range vals {
				shardID, err := toShardID(iv)
				if err != nil {
					return nil, fmt.Errorf("invalid shard ID in per_org_assignment[%q]: %w", key, err)
				}
				shards = append(shards, shardID)
			}
			if len(shards) == 0 {
				return nil, fmt.Errorf("per_org_assignment[%q] must not be empty", key)
			}
			cfg.PerOrgAssignment[key] = shards
		}
	}

	if hashedOwners := tree.Get("hashed_owner_assignment"); hashedOwners != nil {
		vals, ok := hashedOwners.([]any)
		if !ok {
			return nil, errors.New("hashed_owner_assignment must be an array of strings")
		}
		for _, iv := range vals {
			s, ok := iv.(string)
			if !ok {
				return nil, errors.New("hashed_owner_assignment entries must be strings")
			}
			normalized, err := normalizeOwnerHex(s)
			if err != nil {
				return nil, fmt.Errorf("invalid owner address %q in hashed_owner_assignment: %w", s, err)
			}
			if _, exists := cfg.PerOwnerAssignment[normalized]; exists {
				return nil, fmt.Errorf("owner %q appears in both per_owner_assignment and hashed_owner_assignment", normalized)
			}
			cfg.HashedOwnerAssignment[normalized] = true
		}
	}

	return cfg, nil
}

func toShardID(v any) (uint32, error) {
	switch val := v.(type) {
	case int64:
		if val < 0 || val > 4294967295 {
			return 0, fmt.Errorf("shard ID must be in range [0, 4294967295], got %d", val)
		}
		return uint32(val), nil
	case int:
		if val < 0 || val > 4294967295 {
			return 0, fmt.Errorf("shard ID must be in range [0, 4294967295], got %d", val)
		}
		return uint32(val), nil
	default:
		return 0, fmt.Errorf("shard ID must be an integer, got %T", v)
	}
}

func normalizeOwnerHex(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	s = strings.ToLower(s)
	if len(s) != 40 {
		return "", fmt.Errorf("owner address must be 40 hex chars (20 bytes), got %d chars", len(s))
	}
	if _, err := hex.DecodeString(s); err != nil {
		return "", fmt.Errorf("owner address is not valid hex: %w", err)
	}
	return s, nil
}
