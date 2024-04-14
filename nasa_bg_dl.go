package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/adrg/xdg"
	"github.com/mmcdole/gofeed"
)

type directories struct {
	Tmp    string
	Wide   string
	Square string
	Tall   string
}

func getDirectories() directories {
	configFilePath, err := xdg.ConfigFile("nasa_bg_dl/settings.json")
	if err != nil {
		fmt.Println("error opening settings.json")
	}
	fmt.Printf("Looking for settings.jon. The file should be at the following path: %s\n", configFilePath)
	settingsJson, err := os.Open(configFilePath)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println("Unable to open the config file. Did you place it in the right spot?")
	}
	defer func(settingsJson *os.File) {
		err := settingsJson.Close()
		if err != nil {
			errorString := fmt.Sprintf("Couldn't close the settings file. Error: %s", err)
			fmt.Println(errorString)

		}
	}(settingsJson)
	byteValue, _ := io.ReadAll(settingsJson)
	var outputDirectories *directories
	err = json.Unmarshal(byteValue, &outputDirectories)
	if err != nil {
		fmt.Println("Check that you do not have errors in your JSON file.")
		errorString := fmt.Sprintf("Could not unmashal json: %s\n", err)
		fmt.Println(errorString)
	}
	return *outputDirectories
}

// getImage downloads the image and puts it in the directory where it it should end up.
func getImage(image [3]string, outputDirectories directories, wg *sync.WaitGroup) {
	defer wg.Done()
	folder := outputDirectories.Tmp
	filename := image[2]
	filename += strings.ReplaceAll(image[0], " ", "_")
	fileSuffix := ".jpg"
	imageURL := image[1]
	pathFile := folder + filename + fileSuffix
	DownloadFile(pathFile, imageURL)
	orientation := determineRatio(pathFile)
	var newPathFile string
	switch {
	case orientation == "square":
		newPathFile = outputDirectories.Square + filename + fileSuffix

	case orientation == "wide":
		newPathFile = outputDirectories.Wide + filename + fileSuffix

	case orientation == "tall":
		newPathFile = outputDirectories.Tall + filename + fileSuffix
	}
	err := os.Rename(pathFile, newPathFile)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Moved image with title %s to %s\n", image[0], newPathFile)
	}
}

func main() {

	var wg sync.WaitGroup

	outputDirectories := getDirectories()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.nasa.gov/feeds/iotd-feed")
	imagesToGet := getImageMeta(*feed)
	for _, image := range imagesToGet {
		wg.Add(1)
		go getImage(image, outputDirectories, &wg)
	}
	wg.Wait()

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
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
func getImageMeta(feed gofeed.Feed) [3][3]string {
	items := feed.Items
	var images [3][3]string
	for pos, item := range items {
		location, _ := time.LoadLocation("GMT")
		RFC1123NoSeconds := "Mon, 02 Jan 2006 15:04 MST"
		itemTime, err := time.ParseInLocation(RFC1123NoSeconds, item.Published, location)
		if err != nil {
			fmt.Printf("Error parsing time. Error is %s\n", err)
		}
		itemDateString := itemTime.Format(time.DateOnly) + "_"
		images[pos] = [3]string{item.Title, item.Enclosures[0].URL, itemDateString}
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
