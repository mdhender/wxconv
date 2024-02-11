// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package adapters

import (
	"bytes"
	"encoding/xml"
	"github.com/mdhender/wxconv/models/wxml173"
)

type WXML interface {
	BaseVersion() string
}

// UTF8ToWXML converts the data to the associated version of WXML.
func UTF8ToWXML(data []byte) (WXML, error) {
	// verify the xml header
	xmlHeader := []byte("<?xml version='1.0' encoding='utf-16'?>\n")
	if !bytes.HasPrefix(data, xmlHeader) {
		return nil, ErrMissingXMLHeader
	}

	// remove the xml header
	data = data[len(xmlHeader):]

	// read the version from the xml data
	var version struct {
		Version string `xml:"version,attr"`
	}
	if err := xml.Unmarshal(data, &version); err != nil {
		return nil, err
	}
	// log.Printf("UTF8ToWXML: read version %+v\n", version)

	switch version.Version {
	case "1.73":
		srcMap := &wxml173.Map{}
		// convert from xml to a structure that's built just for the conversion
		if err := xml.Unmarshal(data, &srcMap); err != nil {
			return nil, err
		}
		return srcMap, nil
	}
	return nil, ErrUnsupportedVersion
}
