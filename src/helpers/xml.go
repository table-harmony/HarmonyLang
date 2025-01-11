package helpers

import (
	"encoding/xml"
	"fmt"
)

func MapToXml(root string, data map[string]interface{}) ([]byte, error) {
	type KV struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}

	type MapXML struct {
		XMLName xml.Name
		Items   []KV `xml:",any"`
	}

	kvItems := make([]KV, 0, len(data))
	for k, v := range data {
		kvItems = append(kvItems, KV{XMLName: xml.Name{Local: k}, Value: fmt.Sprint(v)})
	}

	mapStruct := MapXML{
		XMLName: xml.Name{Local: root},
		Items:   kvItems,
	}

	return xml.MarshalIndent(mapStruct, "", "  ")
}
