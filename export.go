// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package wxconv

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/wxconv/adapters"
	"github.com/mdhender/wxconv/models/wxx"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ExportJSONFile(m *wxx.Map, path string, debug bool, debugOutputPath string) error {
	step := time.Now()
	if b, err := json.MarshalIndent(m, "", "\t"); err != nil {
		return err
	} else if err = os.WriteFile(path, b, 0644); err != nil {
		return err
	}
	if debug {
		log.Printf("debug: exported  json              in %v\n", time.Now().Sub(step))
	}
	return nil
}

func ExportWXXFile(m *wxx.Map, path string, debug bool, debugOutputPath string) error {
	started := time.Now()

	var outputUtf8xml, outputUtf16Xml string
	if debugOutputPath != "" {
		outputUtf8xml = filepath.Join(debugOutputPath, "input-utf-8.xml")
		outputUtf16Xml = filepath.Join(debugOutputPath, "input-utf-16.xml")
	}

	step := time.Now()
	tmap, err := adapters.WMAPToTMAPv173(m)
	if err != nil {
		return fmt.Errorf("wmapToTmap: %w", err)
	}
	if debug {
		log.Printf("debug: converted wmap to tmap      in %v\n", time.Now().Sub(step))
	}

	step = time.Now()
	data, err := tmap.Encode()
	if err != nil {
		return fmt.Errorf("tmapEncode: %w", err)
	}
	if debug {
		log.Printf("debug: converted tmap to xml       in %v\n", time.Now().Sub(step))
	}

	if outputUtf8xml != "" {
		if err = os.WriteFile(outputUtf8xml, data, 0644); err != nil {
			return err
		}
		if debug {
			log.Printf("debug: created   output-utf-8.xml  in %v\n", time.Now().Sub(started))
		}
	}

	step = time.Now()
	data, err = adapters.UTF8ToUTF16(data)
	if err != nil {
		return fmt.Errorf("utf8Toutf16: %w", err)
	}
	if debug {
		log.Printf("debug: converted utf-8 to utf-16   in %v\n", time.Now().Sub(step))
	}

	if outputUtf16Xml != "" {
		step = time.Now()
		if err = os.WriteFile(outputUtf16Xml, data, 0644); err != nil {
			return err
		}
		if debug {
			log.Printf("debug: created   output-utf-16.xml in %v\n", time.Now().Sub(step))
		}
	}

	step = time.Now()
	data, err = adapters.UTF16ToGZip(data)
	if err != nil {
		return fmt.Errorf("utf1tToGzip: %w", err)
	}
	if debug {
		log.Printf("debug: converted utf-16 to gzip    in %v\n", time.Now().Sub(step))
	}

	step = time.Now()
	if err = os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	if debug {
		log.Printf("debug: exported  wxx file          in %v\n", time.Now().Sub(step))
	}

	if debug {
		log.Printf("debug: completed                   in %v\n", time.Now().Sub(started))
	}

	return nil
}
