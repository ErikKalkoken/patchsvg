// Patchsvg is a command line tool for patching SVG files, so they can be used with the Fyne GUI toolkit.
//
// The SVG library used by [Fyne GUI toolkit] supports a limited subset of SVG only.
// This tool can patch the following known issues:
//
// - Icon has viewPort, but missing width and height (e.g. icons from the website [Pictogrammers Material Design Icons])
//
// [Pictogrammers Material Design Icons]: https://pictogrammers.com/library/mdi/
// [Fyne GUI toolkit]: https://fyne.io/
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	commandName = "patchsvg"
	version     = "0.1.0"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	XMLns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Path    Path     `xml:"path"`
}

type Path struct {
	D string `xml:"d,attr"`
}

func main() {
	flag.Usage = myUsage
	versionFlag := flag.Bool("v", false, "show the current version")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("%s %s\n", commandName, version)
		os.Exit(0)
	}
	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}
	files, err := identifyFiles(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Println("No matching SVG files found.")
		return
	}
	for _, path := range files {
		dat, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		dat2, err := modifyFile(dat)
		if err != nil {
			fmt.Printf("SKIPPED: %s: %s", path, err)
			continue
		}
		if err := os.WriteFile(path, dat2, 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Fixed: %s\n", path)
	}
}

func myUsage() {
	s := "Usage: " + commandName + " [options] <glob pattern>\n\n" +
		"A tool for patching SVG files so they can be used in the Fyne GUI toolkit.\n\n" +
		"Examples:\n" +
		commandName + " icon.svg\n" +
		commandName + " \"resources/*\"\n\n" +
		"Options:\n"
	fmt.Fprint(flag.CommandLine.Output(), s)
	flag.PrintDefaults()
}

func identifyFiles(pattern string) ([]string, error) {
	files1, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	files2 := make([]string, 0)
	for _, path := range files1 {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if fileInfo.IsDir() {
			continue
		}
		if filepath.Ext(path) != ".svg" {
			continue
		}
		files2 = append(files2, path)
	}
	return files2, nil
}

func modifyFile(dat []byte) ([]byte, error) {
	var g SVG
	if err := xml.Unmarshal(dat, &g); err != nil {
		return nil, err
	}
	if err := fixSVG(&g); err != nil {
		return nil, err
	}
	return xml.Marshal(g)
}

func fixSVG(g *SVG) error {
	width, height, err := parseViewBox(g)
	if err != nil {
		return err
	}
	g.Width = strconv.Itoa(width)
	g.Height = strconv.Itoa(height)
	return nil
}

func parseViewBox(g *SVG) (int, int, error) {
	if g.ViewBox == "" {
		return 0, 0, fmt.Errorf("No viewBox defined")
	}
	parts := strings.Split(g.ViewBox, " ")
	if len(parts) != 4 {
		return 0, 0, fmt.Errorf("viewBox not probably defined")
	}
	width, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.Atoi(parts[3])
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}
