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
	"context"
	"encoding/xml"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// SharedChallenge represents a static value to this nonsensical challenge response system.
	// The given challenge must be 11 characters or less. Contents do not matter.
	SharedChallenge = "NintyWhyPls"
)

var baseUrl string
var pool *pgxpool.Pool
var ctx = context.Background()
var isDebug = false

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
	readConfig := Config{}
	err = xml.Unmarshal(ioconfig, &readConfig)
	checkError(err)

	fmt.Println("[i] Initializing core...")
	isDebug = readConfig.Debug

	// Start SQL.
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", readConfig.SQLUser, readConfig.SQLPass, readConfig.SQLAddress, readConfig.SQLDB)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(ctx, dbConf)
	checkError(err)

	// Ensure this PostgreSQL connection is valid.
	defer pool.Close()
	checkError(err)

	baseUrl = readConfig.BaseURL

	// Start the HTTP server.
	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", readConfig.Address)

	r := NewRoute()
	ecs := r.HandleGroup("ecs")
	{
		ecs.Authenticated("CheckDeviceStatus", checkDeviceStatus)
		ecs.Authenticated("NotifyETicketsSynced", notifyETicketsSynced)
		ecs.Authenticated("ListETickets", listETickets)
		ecs.Authenticated("GetETickets", getETickets)
		ecs.Authenticated("PurchaseTitle", purchaseTitle)
		ecs.Unauthenticated("GetECConfig", getECConfig)
		ecs.Authenticated("ListPurchaseHistory", listPurchaseHistory)
	}

	ias := r.HandleGroup("ias")
	{
		ias.Unauthenticated("CheckRegistration", checkRegistration)
		ias.Unauthenticated("GetChallenge", getChallenge)
		ias.Authenticated("GetRegistrationInfo", getRegistrationInfo)
		ias.Unauthenticated("SyncRegistration", syncRegistration)
		ias.Unauthenticated("Register", register)
		ias.Authenticated("Unregister", unregister)
	}

	cas := r.HandleGroup("cas")
	{
		cas.Authenticated("ListItems", listItems)
	}
	log.Fatal(http.ListenAndServe(readConfig.Address, r.Handle()))

	// From here on out, all special cool things should go into their respective handler function.
}
