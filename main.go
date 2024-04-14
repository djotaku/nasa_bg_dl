package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.nasa.gov/feeds/iotd-feed")
	imagesToGet := getImageMeta(*feed)
	for _, image := range imagesToGet {
		filename := strings.ReplaceAll(image[0], " ", "_")
		imageURL := image[1]
		DownloadFile(filename, imageURL)
		orientation := determineRatio(filename + ".jpg")
		fmt.Println(orientation)
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath + ".jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// getImageMeta takes in a feed and returns an array
// with 3 titles and URLs to fetch
func getImageMeta(feed gofeed.Feed) [3][2]string {
	items := feed.Items
	var images [3][2]string
	for pos, item := range items {
		images[pos] = [2]string{item.Title, item.Enclosures[0].URL}
		if pos == 2 {
			break
		}
	}
	return images
}

// determineRatio determines if a given image is
// landscape, portrait, or square
func determineRatio(imagePath string) string {
	file, er := os.Open(imagePath)
	if er != nil {
		fmt.Println("error opening file")
	}
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("error with image")
	}
	ratio := float64(image.Width) / float64(image.Height)
	switch {
	case 0.75 <= ratio && ratio <= 1.33333:
		return "square"
	case ratio > 1.33333:
		return "wide"
	default:
		return "tall"
	}
}
