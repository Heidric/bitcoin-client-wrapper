package v1

var details = []transactionDetails{
	{
		Account:   "accountstring",
		Address:   "1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd",
		Category:  "category",
		Amount:    0.235,
		Label:     "label",
		Vout:      2,
		Fee:       0.015,
		Abandoned: false,
	},
}

var getTransactionResult = transactionData{
	Amount:            0.09,
	Fee:               0.01,
	Confirmations:     3,
	Blockhash:         "blockhashstring",
	Blockindex:        523,
	Txid:              "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d",
	Time:              1560692350,
	Timereceived:      1560691350,
	Bip125Replaceable: "unknown",
	Details:           details,
	Hex:               "rawhexdata",
}

var getNewAddressResult = "1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd"

var sendToAddressResult = "1075db55d416d3ca199f55b6084e2115b9345e16c5cf302fc80e9d5fbf5d48d"
