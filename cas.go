package main

func listItems(e *Envelope) {
	titleId, err := getKey(e.doc, "TitleId")
	if err != nil {
		e.Error(9, "Unable to obtain title.", err)
	}

	e.AddKVNode("ListResultTotalSize", "1")
	e.AddCustomType(Items{
		TitleId: titleId,
		Contents: ContentsMetadata{
			TitleIncluded: false,
			ContentIndex:  1,
		},
		Attributes: []Attributes{
			{
				Name:  "MaxUserInodes",
				Value: "10",
			},
			{
				Name:  "itemComment",
				Value: "Does not catch on fire.",
			},
		},
		Ratings: Ratings{
			Name:   "Testing",
			Rating: 1,
			Age:    13,
		},
		Prices: Prices{
			ItemId: 0,
			Price: Price{
				Amount:   1,
				Currency: "POINTS",
			},
			Limits:      LimitStruct(PR),
			LicenseKind: "PERMANENT",
		},
	})
}
