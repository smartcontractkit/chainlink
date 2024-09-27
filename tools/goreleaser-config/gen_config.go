package main

import (
	"fmt"
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

// Generate creates the goreleaser configuration based on the variation.
func Generate(variation string) config.Project {
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
		Builds:          builds(variation),
		Dockers:         dockers(variation),
		DockerManifests: dockerManifests(variation),
		DockerSigns:     dockerSigns(),
		Checksum: config.Checksum{
			NameTemplate: "checksums.txt",
		},
		Snapshot: config.Snapshot{
			VersionTemplate: "{{ .Env.CHAINLINK_VERSION }}-{{ .ShortCommit }}",
		},
		Partial: config.Partial{
			By: "target",
		},
		Release: config.Release{
			Disable: "true",
		},
		Changelog: config.Changelog{
			Sort: "asc",
			Filters: config.Filters{
				Exclude: []string{
					"^docs:",
					"^test:",
				},
			},
		},
	}

	// Add SBOMs if needed
	if variation == "ci" || variation == "prod" {
		project.SBOMs = []config.SBOM{
			{
				Artifacts: "archive",
			},
		}
	}

	// Add Archives if needed
	if variation == "prod" {
		project.Archives = []config.Archive{
			{
				Format: "binary",
			},
		}
	}

	return project
}

// commonEnv returns the common environment variables used across variations.
func commonEnv() []string {
	return []string{
		`IMAGE_PREFIX={{ if index .Env "IMAGE_PREFIX"  }}{{ .Env.IMAGE_PREFIX }}{{ else }}localhost:5001{{ end }}`,
		`IMAGE_NAME={{ if index .Env "IMAGE_NAME" }}{{ .Env.IMAGE_NAME }}{{ else }}chainlink{{ end }}`,
		`IMAGE_TAG={{ if index .Env "IMAGE_TAG" }}{{ .Env.IMAGE_TAG }}{{ else }}develop{{ end }}`,
		`IMAGE_LABEL_DESCRIPTION="node of the decentralized oracle network, bridging on and off-chain computation"`,
		`IMAGE_LABEL_LICENSES="MIT"`,
		`IMAGE_LABEL_SOURCE="https://github.com/smartcontractkit/{{ .ProjectName }}"`,
	}
}

// builds returns the build configurations based on the variation.
func builds(variation string) []config.Build {
	switch variation {
	case "devspace":
		return []config.Build{
			build(true),
		}
	case "develop", "ci", "prod":
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
		"-X github.com/smartcontractkit/chainlink/v2/core/static.Version={{ .Env.CHAINLINK_VERSION }}",
		"-X github.com/smartcontractkit/chainlink/v2/core/static.Sha={{ .FullCommit }}",
	}
	if isDevspace {
		ldflags[2] = "-X github.com/smartcontractkit/chainlink/v2/core/static.Version={{ .Version }}"
	}

	return config.Build{
		Binary:  "chainlink",
		Targets: []string{"go_first_class"},
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

// dockers returns the docker configurations based on the variation.
func dockers(variation string) []config.Docker {
	var dockers []config.Docker
	switch variation {
	case "devspace":
		dockers = []config.Docker{
			docker("linux-amd64", "linux", "amd64", variation, true),
		}

	case "develop", "ci", "prod":
		architectures := []string{"amd64", "arm64"}
		for _, arch := range architectures {
			dockers = append(dockers, docker("linux-"+arch, "linux", arch, variation, false))
			dockers = append(dockers, docker("linux-"+arch+"-plugins", "linux", arch, variation, false))
		}
	}
	return dockers
}

// docker creates a docker configuration.
func docker(id, goos, goarch, variation string, isDevspace bool) config.Docker {
	extraFiles := []string{"tmp/libs"}
	if strings.Contains(id, "plugins") || isDevspace {
		extraFiles = append(extraFiles, "tmp/plugins")
	}

	buildFlagTemplates := []string{
		fmt.Sprintf("--platform=%s/%s", goos, goarch),
		"--pull",
		"--build-arg=CHAINLINK_USER=chainlink",
		"--build-arg=COMMIT_SHA={{ .FullCommit }}",
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
		`--label=org.opencontainers.image.version={{ .Env.CHAINLINK_VERSION }}`,
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

	if variation == "devspace" {
		dockerConfig.ImageTemplates = []string{"{{ .Env.IMAGE }}"}
	} else {
		base := "{{ .Env.IMAGE_PREFIX }}/{{ .Env.IMAGE_NAME }}"

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

// dockerManifests returns the docker manifest configurations based on the variation.
func dockerManifests(variation string) []config.DockerManifest {
	if variation == "devspace" {
		return []config.DockerManifest{
			{
				NameTemplate:   "{{ .Env.IMAGE }}",
				ImageTemplates: []string{"{{ .Env.IMAGE }}"},
			},
		}
	}

	imageName := "{{ .Env.IMAGE_PREFIX }}/{{ .Env.IMAGE_NAME }}"

	name1 := fmt.Sprintf("%s:{{ .Env.IMAGE_TAG }}", imageName)
	name2 := fmt.Sprintf("%s:sha-{{ .ShortCommit }}", imageName)
	name3 := fmt.Sprintf("%s:{{ .Env.IMAGE_TAG }}-plugins", imageName)
	name4 := fmt.Sprintf("%s:sha-{{ .ShortCommit }}-plugins", imageName)
	return []config.DockerManifest{
		{
			ID:             "tagged",
			NameTemplate:   name1,
			ImageTemplates: manifestImages(name1),
		},
		{
			ID:             "sha",
			NameTemplate:   name2,
			ImageTemplates: manifestImages(name2),
		},
		{
			ID:             "tagged-plugins",
			NameTemplate:   name3,
			ImageTemplates: manifestImages(name3),
		},
		{
			ID:             "sha-plugins",
			NameTemplate:   name4,
			ImageTemplates: manifestImages(name4),
		},
	}
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
