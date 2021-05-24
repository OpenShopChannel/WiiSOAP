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
	"github.com/antchfx/xmlquery"
)

/////////////////////
// SOAP STRUCTURES //
/////////////////////
// The structures may seem repetitive and redundant, but blame WSC's inconsistent SOAP requests.

// Config - WiiSOAP Configuration data.
type Config struct {
	XMLName xml.Name `xml:"Config"`

	Address string `xml:"Address"`
	BaseURL string `xml:"BaseURL"`

	SQLAddress string `xml:"SQLAddress"`
	SQLUser    string `xml:"SQLUser"`
	SQLPass    string `xml:"SQLPass"`
	SQLDB      string `xml:"SQLDB"`

	Debug bool `xml:"Debug"`
}

// Envelope represents the root element of any response, soapenv:Envelope.
type Envelope struct {
	XMLName string `xml:"soapenv:Envelope"`
	SOAPEnv string `xml:"xmlns:soapenv,attr"`
	XSD     string `xml:"xmlns:xsd,attr"`
	XSI     string `xml:"xmlns:xsi,attr"`

	// Represents a soapenv:Body within.
	Body Body

	// Used for internal state tracking.
	doc *xmlquery.Node

	// Common IAS values.
	region   string
	country  string
	language string
}

// Body represents the nested soapenv:Body element as a child on the root element,
// containing the response intended for the action being handled.
type Body struct {
	XMLName string `xml:"soapenv:Body"`

	// Represents the actual response inside
	Response Response
}

// Response describes the inner response format, along with common fields across requests.
type Response struct {
	XMLName xml.Name
	XMLNS   string `xml:"xmlns,attr"`

	// These common fields are persistent across all requests.
	Version            string `xml:"Version"`
	DeviceId           int    `xml:"DeviceId"`
	MessageId          string `xml:"MessageId"`
	TimeStamp          string `xml:"TimeStamp"`
	ErrorCode          int
	ServiceStandbyMode bool `xml:"ServiceStandbyMode"`

	// Allows for <name>[dynamic content]</name> situations.
	CustomFields []interface{}
}

// KVField represents an individual node in form of <XMLName>Contents</XMLName>.
type KVField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// Balance represents a common XML structure.
type Balance struct {
	XMLName  xml.Name `xml:"Balance"`
	Amount   int      `xml:"Amount"`
	Currency string   `xml:"Currency"`
}

type LimitKinds int

const (
	// PR is presumably "purchased".
	PR LimitKinds = 0
	TR            = 1
	DR            = 2
	SR            = 3
	LR            = 4
	AT            = 10000
)

// LimitStruct returns a Limits struct filled for the given kind.
func LimitStruct(kind LimitKinds) Limits {
	names := map[LimitKinds]string{
		PR: "PR",
		TR: "TR",
		DR: "DR",
		SR: "SR",
		LR: "LR",
		AT: "AT",
	}

	return Limits{
		Limits:    kind,
		LimitKind: names[kind],
	}
}

// Limits represents a common XML structure for transaction information.
type Limits struct {
	XMLName   xml.Name   `xml:"Limits"`
	Limits    LimitKinds `xml:"Limits"`
	LimitKind string     `xml:"LimitKind"`
}

// Transactions represents a common XML structure.
type Transactions struct {
	XMLName        xml.Name `xml:"Transactions"`
	TransactionId  string   `xml:"TransactionId"`
	Date           string   `xml:"Date"`
	Type           string   `xml:"Type"`
	TotalPaid      string   `xml:"TotalPaid"`
	Currency       string   `xml:"Currency"`
	ItemId         string   `xml:"ItemId"`
	ItemPricing    string   `xml:"ItemPricing"`
	Limits         Limits   `xml:"Limits"`
	TitleId        string   `xml:"TitleId,omitempty"`
	ItemCode       int      `xml:"ItemCode,omitempty"`
	ReferenceId    int      `xml:"ReferenceId,omitempty"`
	ReferenceValue string   `xml:"ReferenceValue,omitempty"`
}

// Tickets represents the format to inform a console of available titles for its consumption.
type Tickets struct {
	XMLName      xml.Name `xml:"Tickets"`
	TicketId     string   `xml:"TicketId"`
	TitleId      string   `xml:"TitleId"`
	RevokeDate   int      `xml:"RevokeDate"`
	Version      int      `xml:"Version"`
	MigrateCount int      `xml:"MigrateCount"`
	MigrateLimit int      `xml:"MigrateLimit"`
}

// Attributes represents a common structure of the same name.
type Attributes struct {
	XMLName xml.Name `xml:"Attributes"`
	Name    string   `xml:"Name"`
	Value   string   `xml:"Value"`
}

// ContentsMetadata describes data about contents within a title.
type ContentsMetadata struct {
	XMLName       xml.Name `xml:"Contents"`
	TitleIncluded bool     `xml:"TitleIncluded"`
	ContentIndex  int      `xml:"ContentIndex"`
}

// Price holds the price for a title.
type Price struct {
	XMLName  xml.Name `xml:"Price"`
	Amount   int      `xml:"Amount"`
	Currency string   `xml:"Currency"`
}

// Prices describes a common structure for listing prices within a title.
type Prices struct {
	ItemId      int `xml:"ItemId"`
	Price       Price
	Limits      Limits `xml:"Limits"`
	LicenseKind string `xml:"LicenseKind"`
}

// Items allows specifying an overview of a title's contents.
type Items struct {
	XMLName    xml.Name         `xml:"Items"`
	TitleId    string           `xml:"TitleId"`
	Contents   ContentsMetadata `xml:"Contents"`
	Attributes []Attributes     `xml:"Attribute,omitempty"`
	Ratings    Ratings          `xml:"Ratings,omitempty"`
	Prices     Prices           `xml:"Prices,omitempty"`
}

// Ratings allows specifying the rating of an item across multiple properties.
type Ratings struct {
	XMLName xml.Name `xml:"Ratings"`
	Name    string   `xml:"Name"`
	Rating  int      `xml:"Rating"`
	Age     int      `xml:"Age"`
}
