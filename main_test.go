package main

import (
	"github.com/mmcdole/gofeed"
	"os"
	"testing"
)

func TestSquareRatio(t *testing.T) {
	orientation := determineRatio("test_images/square.jpg")
	if orientation != "square" {
		t.Fatalf(`determineRatio("test_images/square.jpg") = %q, want "square"`, orientation)

	}
}

func TestWideRatio(t *testing.T) {
	orientation := determineRatio("test_images/wide.jpg")
	if orientation != "wide" {
		t.Fatalf(`determineRatio("test_images/wide.jpg") = %q, want "wide"`, orientation)

	}
}

func TestTallRatio(t *testing.T) {
	orientation := determineRatio("test_images/tall.jpg")
	if orientation != "tall" {
		t.Fatalf(`determineRatio("test_images/tall.jpg") = %q, want "tall"`, orientation)

	}
}

func TestGetImageMeta(t *testing.T) {
	file, _ := os.Open("test.xml")
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	images := getImageMeta(*feed)
	if len(images) != 3 {
		t.Fatalf(`len(images) = %q, want 3`, len(images))
	}
	if images[0][0] != "The First Space Shuttle" {
		t.Fatalf(`images[0][0] == %q, want "The First Space Shuttle"`, images[0][0])
	}
}