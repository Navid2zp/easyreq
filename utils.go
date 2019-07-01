package easyreq

import (
	"encoding/json"
	"encoding/xml"
)

// Check if string is a json
func IsJson(data []byte) bool {
	var temp map[string]interface{}
	return json.Unmarshal(data, &temp) == nil
}

// Check if string is a XML
func IsXML(data []byte) bool {
	var temp map[string]interface{}
	return xml.Unmarshal(data, &temp) == nil
}
