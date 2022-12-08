package dive

import (
	"fmt"
	"github.com/wagoodman/dive/dive/image"
	"github.com/wagoodman/dive/dive/image/docker"
	"github.com/wagoodman/dive/dive/image/podman"
	"github.com/wagoodman/dive/dive/image/nerdctl"
	"net/url"
	"strings"
)

const (
	SourceUnknown ImageSource = iota
	SourceDockerEngine
	SourcePodmanEngine
	SourceDockerArchive
	SourceNerdctl
)

type ImageSource int

var ImageSources = []string{SourceDockerEngine.String(), SourcePodmanEngine.String(), SourceDockerArchive.String(), SourceNerdctl.String()}

func (r ImageSource) String() string {
	return [...]string{"unknown", "docker", "podman", "docker-archive", "nerdctl"}[r]
}

func ParseImageSource(r string) ImageSource {
	switch r {
	case SourceDockerEngine.String():
		return SourceDockerEngine
	case SourcePodmanEngine.String():
		return SourcePodmanEngine
	case SourceDockerArchive.String():
		return SourceDockerArchive
	case SourceNerdctl.String():
		return SourceNerdctl
	case "docker-tar":
		return SourceDockerArchive
	default:
		return SourceUnknown
	}
}

func DeriveImageSource(image string) (ImageSource, string) {
	u, err := url.Parse(image)
	if err != nil {
		return SourceUnknown, ""
	}

	imageSource := strings.TrimPrefix(image, u.Scheme+"://")

	switch u.Scheme {
	case SourceDockerEngine.String():
		return SourceDockerEngine, imageSource
	case SourcePodmanEngine.String():
		return SourcePodmanEngine, imageSource
	case SourceDockerArchive.String():
		return SourceDockerArchive, imageSource
	case SourceNerdctl.String():
		return SourceNerdctl, imageSource
	case "docker-tar":
		return SourceDockerArchive, imageSource

	}
	return SourceUnknown, ""
}

func GetImageResolver(r ImageSource) (image.Resolver, error) {
	switch r {
	case SourceDockerEngine:
		return docker.NewResolverFromEngine(), nil
	case SourcePodmanEngine:
		return podman.NewResolverFromEngine(), nil
	case SourceNerdctl:
		return nerdctl.NewResolverFromEngine(), nil
	case SourceDockerArchive:
		return docker.NewResolverFromArchive(), nil
	}

	return nil, fmt.Errorf("unable to determine image resolver")
}
