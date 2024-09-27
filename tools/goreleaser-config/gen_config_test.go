package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
	"gopkg.in/yaml.v3"
)

// TestGenerate tests the Generate function.
func TestGenerate(t *testing.T) {
	fixtureRaw, err := os.ReadFile("./testdata/.goreleaser.develop.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var fixture config.Project
	err = yaml.Unmarshal(fixtureRaw, &fixture)
	if err != nil {
		t.Fatal(err)
	}
	project := Generate("develop")
	// diff := cmp.Diff(fixture, project)
	diff := ""
	for _, docker := range project.Dockers {
		// find matching docker.ID in fixture
		var fixtureDocker config.Docker
		for _, fixtureDocker = range fixture.Dockers {
			if docker.ID == fixtureDocker.ID {
				break
			}
		}
		if fixtureDocker.ID == "" {
			t.Errorf("Generate() missing docker.ID %s", docker.ID)
		}
		diff += cmp.Diff(fixtureDocker, docker)
	}

	for _, manifest := range project.DockerManifests {
		// find matching manifest.ID in fixture
		var fixtureManifest config.DockerManifest
		for _, fixMan := range fixture.DockerManifests {
			if manifest.ID == fixMan.ID {
				fixtureManifest = fixMan
				break
			}
		}
		if fixtureManifest.ID == "" {
			t.Errorf("Generate() missing manifest.ID %s", manifest.ID)
		}
		diff += cmp.Diff(fixtureManifest, manifest)
	}

	if diff != "" {
		t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
	}


	// diff rest of fields excluding Dockers and DockerManifests
	project.Dockers = nil
	project.DockerManifests = nil
	fixture.Dockers = nil
	fixture.DockerManifests = nil
	diff = cmp.Diff(fixture, project)
	if diff != "" {
		t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
	}
}
