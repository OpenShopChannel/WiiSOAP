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
	"fmt"
	"log"
)

var ownedTitles *sql.Stmt

func ecsInitialize() {
	var err error
	ownedTitles, err = db.Prepare(`SELECT o.ticket_id, o.title_id, s.version, o.revocation_date
		FROM owned_titles o JOIN shop_titles s
		WHERE o.title_id = s.title_id AND o.account_id = ?`)
	if err != nil {
		log.Fatalf("ecs initialize: error preparing statement: %v\n", err)
	}
}

func checkDeviceStatus(e *Envelope) {
	e.AddCustomType(Balance{
		Amount:   2018,
		Currency: "POINTS",
	})
	e.AddKVNode("ForceSyncTime", "0")
	e.AddKVNode("ExtTicketTime", e.Timestamp())
	e.AddKVNode("SyncTime", e.Timestamp())
}

func notifyETicketsSynced(e *Envelope) {
	// TODO: Implement handling of synchronization timing
}

func listETickets(e *Envelope) {
	rows, err := ownedTitles.Query(e.AccountId())
	if err != nil {
		e.Error(2, "that's all you've got for me? ;3", err)
		return
	}

	// Add all available titles for this account.
	defer rows.Close()
	for rows.Next() {
		var ticketId string
		var titleId string
		var version int
		var revocationDate int
		err = rows.Scan(&ticketId, &titleId, &version, &revocationDate)
		if err != nil {
			e.Error(2, "that's all you've got for me? ;3", err)
			return
		}

		e.AddCustomType(Tickets{
			TicketId:   ticketId,
			TitleId:    titleId,
			Version:    version,
			RevokeDate: revocationDate,

			// We do not support migration.
			MigrateCount: 0,
			MigrateLimit: 0,
		})
	}

	e.AddKVNode("ForceSyncTime", "0")
	e.AddKVNode("ExtTicketTime", e.Timestamp())
	e.AddKVNode("SyncTime", e.Timestamp())
}

func getETickets(e *Envelope) {
	e.AddKVNode("ForceSyncTime", "0")
	e.AddKVNode("ExtTicketTime", e.Timestamp())
	e.AddKVNode("SyncTime", e.Timestamp())
}

func purchaseTitle(e *Envelope) {
	e.AddCustomType(Balance{
		Amount:   2018,
		Currency: "POINTS",
	})
	e.AddCustomType(Transactions{
		TransactionId: "00000000",
		Date:          e.Timestamp(),
		Type:          "PURCHGAME",
	})
	e.AddKVNode("SyncTime", e.Timestamp())
	e.AddKVNode("Certs", "00000000")
	e.AddKVNode("TitleId", "00000000")
	e.AddKVNode("ETickets", "00000000")
}

// genServiceUrl returns a URL with the given service against a configured URL.
// Given a baseUrl of example.com and genServiceUrl("ias", "IdentityAuthenticationSOAP"),
// it would return http://ias.example.com/ias/services/ias/IdentityAuthenticationSOAP.
func genServiceUrl(service string, path string) string {
	return fmt.Sprintf("http://%s.%s/%s/services/%s", service, baseUrl, service, path)
}

func getECConfig(e *Envelope) {
	contentUrl := fmt.Sprintf("http://ccs.%s/ccs/download", baseUrl)
	e.AddKVNode("ContentPrefixURL", contentUrl)
	e.AddKVNode("UncachedContentPrefixURL", contentUrl)
	e.AddKVNode("SystemContentPrefixURL", contentUrl)
	e.AddKVNode("SystemUncachedContentPrefixURL", contentUrl)

	e.AddKVNode("EcsURL", genServiceUrl("ecs", "ECommerceSOAP"))
	e.AddKVNode("IasURL", genServiceUrl("ias", "IdentityAuthenticationSOAP"))
	e.AddKVNode("CasURL", genServiceUrl("cas", "CatalogingSOAP"))
	e.AddKVNode("NusURL", genServiceUrl("nus", "NetUpdateSOAP"))
}
