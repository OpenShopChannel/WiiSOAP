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
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// SharedChallenge represents a static value to this nonsensical challenge response system.
	// The given challenge must be 11 characters or less. Contents do not matter.
	SharedChallenge = "NintyWhyPls"
)

var db *sql.DB

// checkError makes error handling not as ugly and inefficient.
func checkError(err error) {
	if err != nil {
		log.Fatalf("WiiSOAP forgot how to drive and suddenly crashed! Reason: %v\n", err)
	}
}

func main() {
	// Initial Start.
	fmt.Println("WiiSOAP 0.2.6 Kawauso\n[i] Reading the Config...")

	// Check the Config.
	ioconfig, err := ioutil.ReadFile("./config.xml")
	checkError(err)
	CON := Config{}
	err = xml.Unmarshal(ioconfig, &CON)
	checkError(err)

	fmt.Println("[i] Initializing core...")

	// Start SQL.
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", CON.SQLUser, CON.SQLPass, CON.SQLAddress, CON.SQLDB))
	checkError(err)

	// Close SQL after everything else is done.
	defer db.Close()
	err = db.Ping()
	checkError(err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Initialize handlers.
	ecsInitialize()
	iasInitialize()
	routeInitialize()

	// Start the HTTP server.
	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", CON.Address)

	r := NewRoute()
	ecs := r.HandleGroup("ecs")
	{
		ecs.Authenticated("CheckDeviceStatus", checkDeviceStatus)
		ecs.Authenticated("NotifyETicketsSynced", notifyETicketsSynced)
		ecs.Authenticated("ListETickets", listETickets)
		ecs.Authenticated("GetETickets", getETickets)
		ecs.Authenticated("PurchaseTitle", purchaseTitle)

	}

	ias := r.HandleGroup("ias")
	{
		ias.Unauthenticated("CheckRegistration", checkRegistration)
		ias.Unauthenticated("GetChallenge", getChallenge)
		ias.Authenticated("GetRegistrationInfo", getRegistrationInfo)
		ias.Unauthenticated("Register", register)
		ias.Authenticated("Unregister", unregister)
	}
	log.Fatal(http.ListenAndServe(CON.Address, r.Handle()))

	// From here on out, all special cool things should go into their respective handler function.
}
