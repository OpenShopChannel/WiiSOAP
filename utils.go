//	Copyright (C) 2018-2020 CornierKhan1
//
//	WiiSOAP is SOAP Server Software, designed specifically to handle Wii Shop Channel SOAP.
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as published
//    by the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see http://www.gnu.org/licenses/.

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var namespaceParse = regexp.MustCompile(`^urn:(.{3})\.wsapi\.broadon\.com/(.*)$`)

// parseAction interprets contents along the lines of "urn:ecs.wsapi.broadon.com/CheckDeviceStatus",
//where "CheckDeviceStatus" is the action to be performed.
func parseAction(original string) (string, string) {
	// Intended to return the original string, the service's name and the name of the action.
	matches := namespaceParse.FindStringSubmatch(original)
	if len(matches) != 3 {
		// It seems like the passed action was not matched properly.
		return "", ""
	}

	service := matches[1]
	action := matches[2]
	return service, action
}

// NewEnvelope returns a new Envelope with proper attributes initialized.
func NewEnvelope(service string, action string, body []byte) (*Envelope, error) {
	// Get a sexy new timestamp to use.
	timestampNano := fmt.Sprint(time.Now().UTC().UnixNano())[0:13]

	// Tidy up parsed document for easier usage going forward.
	doc, err := normalise(service, action, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	// Return an envelope with properly set defaults to respond with.
	e := Envelope{
		SOAPEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		XSD:     "http://www.w3.org/2001/XMLSchema",
		XSI:     "http://www.w3.org/2001/XMLSchema-instance",
		Body: Body{
			Response: Response{
				XMLName: xml.Name{Local: action + "Response"},
				XMLNS:   "urn:" + service + ".wsapi.broadon.com",

				TimeStamp: timestampNano,
			},
		},
		doc: doc,
	}

	// Obtain common request values.
	err = e.ObtainCommon()
	if err != nil {
		return nil, err
	}

	return &e, nil
}

// Timestamp returns a shared timestamp for this request.
func (e *Envelope) Timestamp() string {
	return e.Body.Response.TimeStamp
}

// DeviceId returns the Device ID for this request.
func (e *Envelope) DeviceId() int {
	return e.Body.Response.DeviceId
}

// Region returns the region for this request. It should be only used in IAS-related requests.
func (e *Envelope) Region() string {
	return e.region
}

// Country returns the region for this request. It should be only used in IAS-related requests.
func (e *Envelope) Country() string {
	return e.country
}

// Language returns the region for this request. It should be only used in IAS-related requests.
func (e *Envelope) Language() string {
	return e.language
}

// AccountId returns the account ID for this request. It should be only used in authenticated requests.
// If for whatever reason AccountId is not present (such as in SyncRegistration),
// it will return an error. Please design in a way so that this is not an issue.
func (e *Envelope) AccountId() (int64, error) {
	accountId, err := getKey(e.doc, "AccountId")
	if err != nil {
		return 0, nil
	}

	return strconv.ParseInt(accountId, 10, 64)
}

// ObtainCommon interprets a given node, and updates the envelope with common key values.
func (e *Envelope) ObtainCommon() error {
	var err error
	doc := e.doc

	// These fields are common across all requests.
	e.Body.Response.Version, err = getKey(doc, "Version")
	if err != nil {
		return err
	}
	deviceIdString, err := getKey(doc, "DeviceId")
	if err != nil {
		return err
	}
	e.Body.Response.DeviceId, err = strconv.Atoi(deviceIdString)
	if err != nil {
		return err
	}
	e.Body.Response.MessageId, err = getKey(doc, "MessageId")
	if err != nil {
		return err
	}

	// These are as well, but we do not need to send them back in our response.
	e.region, err = getKey(doc, "Region")
	if err != nil {
		return err
	}
	e.country, err = getKey(doc, "Country")
	if err != nil {
		return err
	}
	e.language, err = getKey(doc, "Language")
	if err != nil {
		return err
	}

	return nil
}

// AddKVNode adds a given key by name to a specified value, such as <key>value</key>.
func (e *Envelope) AddKVNode(key string, value string) {
	e.Body.Response.CustomFields = append(e.Body.Response.CustomFields, KVField{
		XMLName: xml.Name{Local: key},
		Value:   value,
	})
}

// AddCustomType adds a given key by name to a specified structure.
func (e *Envelope) AddCustomType(customType interface{}) {
	e.Body.Response.CustomFields = append(e.Body.Response.CustomFields, customType)
}

// becomeXML marshals the Envelope object, returning the intended boolean state on success.
// ..there has to be a better way to do this, TODO.
func (e *Envelope) becomeXML() (bool, string) {
	// Non-zero error codes indicate a failure.
	intendedStatus := e.Body.Response.ErrorCode == 0

	var contents []byte
	var err error

	// If we're in debug mode, pretty print XML.
	if isDebug {
		contents, err = xml.MarshalIndent(e, "", "  ")
	} else {
		contents, err = xml.Marshal(e)
	}

	if err != nil {
		return false, "an error occurred marshalling XML: " + err.Error()
	} else {
		// Add XML header on top of existing contents.
		result := xml.Header + string(contents)
		return intendedStatus, result
	}
}

// Error sets the necessary keys for this SOAP response to reflect the given error.
func (e *Envelope) Error(errorCode int, reason string, err error) {
	e.Body.Response.ErrorCode = errorCode

	// Ensure all additional fields are empty to avoid conflict.
	e.Body.Response.CustomFields = nil

	e.AddKVNode("UserReason", reason)
	e.AddKVNode("ServerReason", err.Error())
}

// normalise parses a document, returning a document with only the request type's child nodes, stripped of prefix.
func normalise(service string, action string, reader io.Reader) (*xmlquery.Node, error) {
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return nil, err
	}

	// Find the keys for this element named after the action.
	result := doc.SelectElement("//" + service + ":" + action)
	if result == nil {
		return nil, errors.New("missing root node")
	}
	stripNamespace(result)

	return result, nil
}

// stripNamespace removes a prefix from nodes, changing a key from "ias:Version" to "Version".
// It is based off of https://github.com/antchfx/xmlquery/issues/15#issuecomment-567575075.
func stripNamespace(node *xmlquery.Node) {
	node.Prefix = ""

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		stripNamespace(child)
	}
}

// getKey returns the value for a child key from a node, if documented.
func getKey(doc *xmlquery.Node, key string) (string, error) {
	node := xmlquery.FindOne(doc, "//"+key)

	if node == nil {
		return "", errors.New("missing mandatory key named " + key)
	} else {
		return node.InnerText(), nil
	}
}

// Derived from https://stackoverflow.com/a/31832326, adding numbers
const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func debugPrint(v ...interface{}) {
	if !isDebug {
		return
	}

	log.Print(v...)
}
