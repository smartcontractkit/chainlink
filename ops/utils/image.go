package utils

import (
	"fmt"

	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Image implements the struct for fetching images
type Image struct {
	Name  string
	Tag   string
	Img   *docker.RemoteImage
	Local *docker.Image
}

// Pull retrieves the specified container image
func (i *Image) Pull(ctx *pulumi.Context) error {
	msg := LogStatus(fmt.Sprintf("Pulling %s", i.Name))
	img, err := docker.NewRemoteImage(ctx, i.Name, &docker.RemoteImageArgs{
		Name:        pulumi.String(i.Tag),
		KeepLocally: pulumi.BoolPtr(true),
	})
	i.Img = img
	return msg.Check(err)
}

// Build the image for the specified dockerfile in the YAML config
func (i *Image) Build(ctx *pulumi.Context, context string, dockerfile string) error {
	// build local image
	msg := LogStatus(fmt.Sprintf("Building %s", i.Name))
	img, err := docker.NewImage(ctx, i.Name, &docker.ImageArgs{
		ImageName: pulumi.String(i.Tag),
		// LocalImageName: pulumi.String(i.Tag),
		SkipPush: pulumi.Bool(true),
		Registry: docker.ImageRegistryArgs{},
		Build: docker.DockerBuildArgs{
			Context:    pulumi.String(context),
			Dockerfile: pulumi.String(dockerfile),
		},
	})
	i.Local = img
	return msg.Check(err)
}
