// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements an application to import and export map data.
package main

import (
	"flag"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/mdhender/wxconv"
	"github.com/mdhender/wxconv/models/wxx"
	"log"
	"os"
	"path/filepath"
)

var (
	mv = semver.Version{
		Major:      0,
		Minor:      1,
		Patch:      0,
		PreRelease: "alpha",
	}
)

func main() {
	log.SetFlags(log.Ltime)

	var debug bool
	flag.BoolVar(&debug, "debug", debug, "show debug output")

	var showVersion bool
	flag.BoolVar(&showVersion, "version", showVersion, "show version and quit")

	var importWXXFile string
	flag.StringVar(&importWXXFile, "import", importWXXFile, ".wxx file to load")

	var importJSONFile string
	flag.StringVar(&importJSONFile, "import-json", importJSONFile, ".json file to load")

	var exportWXXFile string
	flag.StringVar(&exportWXXFile, "export", exportWXXFile, ".wxx file to create")

	var exportJSONFile string
	flag.StringVar(&exportJSONFile, "export-json", exportJSONFile, ".json file to create")

	var debugOutputPath string
	flag.StringVar(&debugOutputPath, "debug-output-path", debugOutputPath, "path to create debug files in")

	flag.Parse()
	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(2)
	}

	if showVersion {
		fmt.Printf("%s\n", mv.String())
		os.Exit(0)
	}

	hasJSONImport, hasJSONExport := importJSONFile != "", exportJSONFile != ""
	hasWXXImport, hasWXXExport := importWXXFile != "", exportWXXFile != ""

	if hasJSONImport && hasWXXImport {
	} else if hasJSONImport {
		// input must exist and be a regular file
		sb, err := os.Stat(importJSONFile)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		} else if sb.IsDir() {
			log.Fatalf("error: %s: is a directory\n", importJSONFile)
		} else if !sb.Mode().IsRegular() {
			log.Fatalf("error: %s: is not a file\n", importJSONFile)
		}
	} else if hasWXXImport {
		// input must exist and be a regular file
		sb, err := os.Stat(importWXXFile)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		} else if sb.IsDir() {
			log.Fatalf("error: %s: is a directory\n", importWXXFile)
		} else if !sb.Mode().IsRegular() {
			log.Fatalf("error: %s: is not a file\n", importWXXFile)
		}
	} else {
		log.Fatalf("error: you must specify a file to import\n")
	}

	if debugOutputPath != "" {
		if s, err := filepath.Abs(debugOutputPath); err != nil {
			log.Fatalf("output: %v\n", err)
		} else {
			debugOutputPath = s
		}
		sb, err := os.Stat(debugOutputPath)
		if err != nil {
			log.Fatalf("%s: %v", debugOutputPath, err)
		} else if !sb.IsDir() {
			log.Fatalf("%s: is not a directory\n", debugOutputPath)
		}
	}

	var m *wxx.Map
	var err error

	if hasJSONImport {
		log.Fatalf("error: json import not implemented yet\n")
	} else if hasWXXImport {
		m, err = wxconv.ImportWXXFile(importWXXFile, debug, debugOutputPath)
		if err != nil {
			log.Printf("import: %s\n", importWXXFile)
			log.Fatalf("import: %v", err)
		}
	} else {
		log.Fatalf("error: you must specify a file to import\n")
	}

	if hasJSONExport {
		if err = wxconv.ExportJSONFile(m, exportJSONFile, true, debugOutputPath); err != nil {
			log.Printf("export: %s", exportWXXFile)
			log.Fatalf("export: %v", err)
		}
		log.Printf("created %s\n", exportJSONFile)
	}

	if hasWXXExport {
		err := wxconv.ExportWXXFile(m, exportWXXFile, debug, debugOutputPath)
		if err != nil {
			log.Printf("export: %s", exportWXXFile)
			log.Fatalf("export: %v", err)
		}
		log.Printf("created %s\n", exportWXXFile)
	}
}
