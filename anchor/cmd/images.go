package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// ImagesListOutput used to interact with template
type ImagesListOutput struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Tag         string `json:"Tag"`
	Size        string `json:"Size"`
	CreatedTime string `json:"CreatedTime"`
}

// ImagesList return images
func ImagesList() ([]ImagesListOutput, error) {
	images, err := DockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("")
	}
	imagesListOutput := []ImagesListOutput{}
	for _, image := range images {
		var name, tag string
		if len(image.RepoTags) > 0 {
			nametags := strings.Split(image.RepoTags[0], ":")
			if len(nametags) >= 2 {
				name = nametags[0]
				tag = nametags[1]
			} else {
				name = image.RepoTags[0]
				tag = ""
			}
		} else {
			name = ""
			tag = ""
		}
		imagesListOutput = append(imagesListOutput, ImagesListOutput{
			ID:          image.ID,
			Name:        name,
			Tag:         tag,
			Size:        strconv.FormatFloat(float64(image.Size)/1024/1024, 'f', 3, 64),
			CreatedTime: time.Unix(image.Created, 0).Format("2006-01-02 15:04"),
		})
	}
	return imagesListOutput, nil
}

// ImageGet return image
func ImageGet(mid string) (types.ImageInspect, []byte, error) {
	return DockerClient.ImageInspectWithRaw(context.Background(), mid)
}

// ImageDelete delete an image
func ImageDelete(mid string) ([]types.ImageDeleteResponseItem, error) {
	return DockerClient.ImageRemove(context.Background(), mid, types.ImageRemoveOptions{})
}
