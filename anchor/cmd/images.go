package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/yinwoods/anchor/anchor/util"
	"golang.org/x/net/context"
)

const (
	dockerImageURL = "http://localhost:8089/api/image"
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

	resp, _ := util.HTTPGet(dockerImageURL)
	var images []types.ImageSummary
	json.Unmarshal(resp, &images)
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

// ImageCreate pull an image
func ImageCreate(name string) (io.ReadCloser, error) {
	return DockerClient.ImagePull(context.Background(), name, types.ImagePullOptions{})
}

// ImageGet return image
func ImageGet(id string) (types.ImageInspect, error) {
	resp, err := util.HTTPGet(fmt.Sprintf("%s/%s/json", dockerImageURL, id))
	if err != nil {
		return types.ImageInspect{}, err
	}
	var image types.ImageInspect
	json.Unmarshal(resp, &image)
	return image, nil
}

// ImageDelete delete an image
func ImageDelete(id string) error {
	_, err := util.HTTPDelete(fmt.Sprintf("%s/%s", dockerImageURL, id))
	return err
}
