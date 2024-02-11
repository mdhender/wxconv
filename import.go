// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package wxconv

import (
	"fmt"
	"github.com/mdhender/wxconv/adapters"
	"github.com/mdhender/wxconv/models/wxx"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ImportJSONFile(path string, debug bool, debugOutpathPath string) (*wxx.Map, error) {
	return nil, fmt.Errorf("not implemented")
}

func ImportWXXFile(path string, debug bool, debugOutputPath string) (*wxx.Map, error) {
	started := time.Now()
	step := started

	var inputUtf8xml, inputUtf16xml string
	if debugOutputPath != "" {
		inputUtf8xml = filepath.Join(debugOutputPath, "input-utf-8.xml")
		inputUtf16xml = filepath.Join(debugOutputPath, "input-utf-16.xml")
	}

	// read input
	step = time.Now()
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if debug {
		log.Printf("debug: read      input             in %v\n", time.Now().Sub(step))
	}

	// unzip input and convert to UTF-16
	step = time.Now()
	src, err = adapters.GZipToUTF16(src)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if debug {
		log.Printf("debug: converted gzip to utf-16    in %v\n", time.Now().Sub(step))
	}

	if inputUtf16xml != "" {
		// save UTF-16 input
		step = time.Now()
		if err = os.WriteFile(inputUtf16xml, src, 0644); err != nil {
			return nil, err
		}
		if debug {
			log.Printf("debug: created   input-utf-16.xml  in %v\n", time.Now().Sub(step))
		}
	}

	// convert input from UTF-16 to UTF-8
	step = time.Now()
	src, err = adapters.UTF16ToUTF8(src)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if debug {
		log.Printf("debug: converted utf-16 to utf-8   in %v\n", time.Now().Sub(step))
	}

	if inputUtf8xml != "" {
		// save UTF-8 input
		step = time.Now()
		if err = os.WriteFile(inputUtf8xml, src, 0644); err != nil {
			return nil, err
		}
		if debug {
			log.Printf("debug: created   input-utf-8.xml   in %v\n", time.Now().Sub(step))
		}
	}

	// convert UTF-8 to WXML
	step = time.Now()
	wxml, err := adapters.UTF8ToWXML(src)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if debug {
		log.Printf("debug: converted utf-8 to wxml     in %v\n", time.Now().Sub(step))
	}

	// convert the WXML to WMAP
	step = time.Now()
	wmap, err := adapters.WXMLToWXX(wxml)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if debug {
		log.Printf("debug: converted wxml to wmap      in %v\n", time.Now().Sub(step))
	}

	if debug {
		log.Printf("debug: completed import            in %v\n", time.Now().Sub(started))
	}

	return wmap, nil
}
