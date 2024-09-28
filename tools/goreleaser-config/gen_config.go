package main

import (
	"fmt"
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

// Generate creates the goreleaser configuration based on the environment.
var validEnvironments = []string{"devspace", "develop", "prod"}

func Generate(environment string) config.Project {
	checkEnvironments(environment)

	project := config.Project{
		ProjectName: "chainlink",
		Version:     2,
		Env:         commonEnv(),
		Before: config.Before{
			Hooks: []config.Hook{
				{
					Cmd: "go mod tidy",
				},
				{
					Cmd: "./tools/bin/goreleaser_utils before_hook",
				},
			},
		},
		Builds:          builds(environment),
		Dockers:         dockers(environment),
		DockerManifests: dockerManifests(environment),
		DockerSigns:     dockerSigns(),
		Checksum: config.Checksum{
			NameTemplate: "checksums.txt",
		},
		Snapshot: config.Snapshot{
			VersionTemplate: "{{ .Env.VERSION }}-{{ .ShortCommit }}",
		},
		Nightly: config.Nightly{
			VersionTemplate: "{{ .Env.VERSION }}-{{ .Env.IMAGE_TAG }}",
		},
		Partial: config.Partial{
			By: "target",
		},
		Release: config.Release{
			Disable: "true",
		},
		Changelog: config.Changelog{
			Disable: "true",
			Sort:    "asc",
			Filters: config.Filters{
				Exclude: []string{
					"^docs:",
					"^test:",
				},
			},
		},
		Archives: []config.Archive{
			{
				Format: "binary",
			},
		},
	}

	// Add SBOMs if needed
	if environment == "prod" {
		project.Changelog.Disable = "false"
		project.Archives = []config.Archive{
			{
				Format: "tar.gz",
			},
		}
		project.SBOMs = []config.SBOM{
			{
				Artifacts: "archive",
			},
		}
	}

	return project
}

func checkEnvironments(environment string) {
	valid := false
	for _, env := range validEnvironments {
		if environment == env {
			valid = true
			break
		}
	}
	if !valid {
		panic(fmt.Sprintf("invalid environment: %s, valid environments are %v", environment, validEnvironments))
	}
}

// commonEnv returns the common environment variables used across environments.
func commonEnv() []string {
	return []string{
		`IMAGE_PREFIX={{ if index .Env "IMAGE_PREFIX"  }}{{ .Env.IMAGE_PREFIX }}{{ else }}localhost:5001{{ end }}`,
		`IMAGE_TAG={{ if index .Env "IMAGE_TAG" }}{{ .Env.IMAGE_TAG }}{{ else }}develop{{ end }}`,
		`VERSION={{ if index .Env "CHAINLINK_VERSION" }}{{ .Env.CHAINLINK_VERSION }}{{ else }}v0.0.0-local{{ end }}`,
		`IMAGE_LABEL_DESCRIPTION="node of the decentralized oracle network, bridging on and off-chain computation"`,
		`IMAGE_LABEL_LICENSES="MIT"`,
		`IMAGE_LABEL_SOURCE="https://github.com/smartcontractkit/{{ .ProjectName }}"`,
	}
}

// builds returns the build configurations based on the environment.
func builds(environment string) []config.Build {
	switch environment {
	case "devspace":
		return []config.Build{
			build(true),
		}
	case "develop", "prod":
		return []config.Build{
			build(false),
		}

	default:
		return nil
	}
}

// build creates a build configuration.
func build(isDevspace bool) config.Build {
	ldflags := []string{
		"-s -w -r=$ORIGIN/libs",
		"-X github.com/smartcontractkit/chainlink/v2/core/static.Version={{ .Env.VERSION }}",
		"-X github.com/smartcontractkit/chainlink/v2/core/static.Sha={{ .FullCommit }}",
	}
	if isDevspace {
		ldflags[2] = "-X github.com/smartcontractkit/chainlink/v2/core/static.Version={{ .Version }}"
	}

	return config.Build{
		Binary:          "chainlink",
		NoUniqueDistDir: "true",
		Targets:         []string{"go_first_class"},
		Hooks: config.BuildHookConfig{
			Post: []config.Hook{
				{Cmd: "./tools/bin/goreleaser_utils build_post_hook {{ dir .Path }}"},
			},
		},
		BuildDetails: config.BuildDetails{
			Flags:   []string{"-trimpath", "-buildmode=pie"},
			Ldflags: ldflags,
		},
	}
}

// dockers returns the docker configurations based on the environment.
func dockers(environment string) []config.Docker {
	var dockers []config.Docker
	switch environment {
	case "devspace":
		dockers = []config.Docker{
			docker("linux-amd64", "linux", "amd64", environment, true),
		}

	case "develop", "prod":
		architectures := []string{"amd64", "arm64"}
		var imageNames []string
		if environment == "develop" {
			imageNames = []string{"chainlink", "ccip"}
		} else {
			// on prod, we have the ECR prefix for the image
			imageNames = []string{"chainlink/chainlink", "chainlink/ccip"}
		}

		for _, imageName := range imageNames {
			for _, arch := range architectures {
				id := fmt.Sprintf("linux-%s-%s", arch, imageName)
				pluginId := fmt.Sprintf("%s-plugins", id)

				dockers = append(dockers, docker(id, "linux", arch, environment, false))
				dockers = append(dockers, docker(pluginId, "linux", arch, environment, false))
			}
		}
	}
	return dockers
}

// docker creates a docker configuration.
func docker(id, goos, goarch, environment string, isDevspace bool) config.Docker {
	isCCIP := strings.Contains(id, "ccip")
	isPlugins := strings.Contains(id, "plugins")
	extraFiles := []string{"tmp/libs"}
	if isPlugins || isDevspace {
		extraFiles = append(extraFiles, "tmp/plugins")
	}
	if isCCIP {
		extraFiles = append(extraFiles, "ccip/config")
	}

	buildFlagTemplates := []string{
		fmt.Sprintf("--platform=%s/%s", goos, goarch),
		"--pull",
		"--build-arg=CHAINLINK_USER=chainlink",
		"--build-arg=COMMIT_SHA={{ .FullCommit }}",
	}

	if strings.Contains(id, "ccip") {
		buildFlagTemplates = append(buildFlagTemplates,
			"--build-arg=CL_CHAIN_DEFAULTS=/chainlink/ccip-config")
	}

	if strings.Contains(id, "plugins") || isDevspace {
		buildFlagTemplates = append(buildFlagTemplates,
			"--build-arg=CL_MEDIAN_CMD=chainlink-feeds",
			"--build-arg=CL_MERCURY_CMD=chainlink-mercury",
			"--build-arg=CL_SOLANA_CMD=chainlink-solana",
			"--build-arg=CL_STARKNET_CMD=chainlink-starknet",
		)
	}

	buildFlagTemplates = append(buildFlagTemplates,
		`--label=org.opencontainers.image.created={{ .Date }}`,
		`--label=org.opencontainers.image.description={{ .Env.IMAGE_LABEL_DESCRIPTION }}`,
		`--label=org.opencontainers.image.licenses={{ .Env.IMAGE_LABEL_LICENSES }}`,
		`--label=org.opencontainers.image.revision={{ .FullCommit }}`,
		`--label=org.opencontainers.image.source={{ .Env.IMAGE_LABEL_SOURCE }}`,
		`--label=org.opencontainers.image.title={{ .ProjectName }}`,
		`--label=org.opencontainers.image.version={{ .Env.VERSION }}`,
		`--label=org.opencontainers.image.url={{ .Env.IMAGE_LABEL_SOURCE }}`,
	)

	dockerConfig := config.Docker{
		ID:                 id,
		Dockerfile:         "core/chainlink.goreleaser.Dockerfile",
		Use:                "buildx",
		Goos:               goos,
		Goarch:             goarch,
		Files:              extraFiles,
		BuildFlagTemplates: buildFlagTemplates,
	}

	if environment == "devspace" {
		dockerConfig.ImageTemplates = []string{"{{ .Env.IMAGE }}"}
	} else {
		var base string
		if isCCIP {
			base = "{{ .Env.IMAGE_PREFIX }}/ccip"
		} else {
			base = "{{ .Env.IMAGE_PREFIX }}/chainlink"
		}

		imageTemplates := []string{}
		if strings.Contains(id, "plugins") {
			taggedBase := fmt.Sprintf("%s:{{ .Env.IMAGE_TAG }}-plugins", base)
			// We have a default, non-arch specific image for plugins that defaults to amd64
			if goarch == "amd64" {
				imageTemplates = append(imageTemplates, taggedBase)
			}
			imageTemplates = append(imageTemplates,
				fmt.Sprintf("%s-%s", taggedBase, archSuffix(id)),
				fmt.Sprintf("%s:sha-{{ .ShortCommit }}-plugins-%s", base, archSuffix(id)))
		} else {
			taggedBase := fmt.Sprintf("%s:{{ .Env.IMAGE_TAG }}", base)
			// We have a default, non-arch specific image for plugins that defaults to amd64
			if goarch == "amd64" {
				imageTemplates = append(imageTemplates, taggedBase)
			}
			imageTemplates = append(imageTemplates,
				fmt.Sprintf("%s-%s", taggedBase, archSuffix(id)),
				fmt.Sprintf("%s:sha-{{ .ShortCommit }}-%s", base, archSuffix(id)))
		}

		dockerConfig.ImageTemplates = imageTemplates
	}

	return dockerConfig
}

// archSuffix returns the architecture suffix for image tags.
func archSuffix(id string) string {
	if strings.Contains(id, "arm64") {
		return "arm64"
	}
	return "amd64"
}

// dockerManifests returns the docker manifest configurations based on the environment.
func dockerManifests(environment string) []config.DockerManifest {
	if environment == "devspace" {
		return []config.DockerManifest{
			{
				NameTemplate:   "{{ .Env.IMAGE }}",
				ImageTemplates: []string{"{{ .Env.IMAGE }}"},
			},
		}
	}

	// Define the image names based on the environment
	imageNames := []string{"chainlink", "ccip"}

	var manifests []config.DockerManifest

	for _, imageName := range imageNames {
		fullImageName := fmt.Sprintf("{{ .Env.IMAGE_PREFIX }}/%s", imageName)

		manifestConfigs := []struct {
			ID     string
			Suffix string
		}{
			{ID: "tagged", Suffix: ":{{ .Env.IMAGE_TAG }}"},
			{ID: "sha", Suffix: ":sha-{{ .ShortCommit }}"},
			{ID: "tagged-plugins", Suffix: ":{{ .Env.IMAGE_TAG }}-plugins"},
			{ID: "sha-plugins", Suffix: ":sha-{{ .ShortCommit }}-plugins"},
		}
		for _, cfg := range manifestConfigs {
			nameTemplate := fmt.Sprintf("%s%s", fullImageName, cfg.Suffix)
			manifests = append(manifests, config.DockerManifest{
				ID:             fmt.Sprintf("%s-%s", cfg.ID, imageName),
				NameTemplate:   nameTemplate,
				ImageTemplates: manifestImages(nameTemplate),
			})
		}
	}

	return manifests
}

// manifestImages generates image templates for docker manifests.
func manifestImages(imageName string) []string {
	architectures := []string{"amd64", "arm64"}
	var images []string
	// Add the default image for tagged images
	if !strings.Contains(imageName, "sha") {
		images = append(images, imageName)
	}
	for _, arch := range architectures {
		images = append(images, fmt.Sprintf("%s-%s", imageName, arch))
	}
	return images
}

// dockerSigns returns the docker sign configurations.
func dockerSigns() []config.Sign {
	return []config.Sign{
		{
			Artifacts: "all",
			Args: []string{
				"sign",
				"${artifact}",
				"--yes",
			},
		},
	}
}
