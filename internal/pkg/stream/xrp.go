package stream

import "time"

type LedgerResp struct {
	Result struct {
		Ledger struct {
			AccountHash         string    `json:"account_hash"`
			CloseFlags          int       `json:"close_flags"`
			CloseTime           int       `json:"close_time"`
			CloseTimeHuman      string    `json:"close_time_human"`
			CloseTimeIso        time.Time `json:"close_time_iso"`
			CloseTimeResolution int       `json:"close_time_resolution"`
			Closed              bool      `json:"closed"`
			LedgerHash          string    `json:"ledger_hash"`
			LedgerIndex         string    `json:"ledger_index"`
			ParentCloseTime     int       `json:"parent_close_time"`
			ParentHash          string    `json:"parent_hash"`
			TotalCoins          string    `json:"total_coins"`
			TransactionHash     string    `json:"transaction_hash"`
		} `json:"ledger"`
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
		Status      string `json:"status"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
}

type LedgerAccountTx struct {
	Result struct {
		Account        string     `json:"account"`
		LedgerIndexMax int        `json:"ledger_index_max"`
		LedgerIndexMin int        `json:"ledger_index_min"`
		Limit          int        `json:"limit"`
		Status         string     `json:"status"`
		Transactions   []LedgerTx `json:"transactions"`
		Validated      bool       `json:"validated"`
	} `json:"result"`
}

type LedgerTx struct {
	Meta struct {
		AffectedNodes []struct {
			ModifiedNode struct {
				FinalFields struct {
					Account    string `json:"Account"`
					Balance    string `json:"Balance"`
					Flags      int    `json:"Flags"`
					OwnerCount int    `json:"OwnerCount"`
					Sequence   int    `json:"Sequence"`
				} `json:"FinalFields"`
				LedgerEntryType string `json:"LedgerEntryType"`
				LedgerIndex     string `json:"LedgerIndex"`
				PreviousFields  struct {
					Balance string `json:"Balance"`
				} `json:"PreviousFields"`
				PreviousTxnID     string `json:"PreviousTxnID"`
				PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
			} `json:"ModifiedNode"`
		} `json:"AffectedNodes"`
		TransactionIndex  int    `json:"TransactionIndex"`
		TransactionResult string `json:"TransactionResult"`
		DeliveredAmount   string `json:"delivered_amount"`
	} `json:"meta"`
	Tx struct {
		Account            string `json:"Account"`
		Amount             string `json:"Amount"`
		DeliverMax         string `json:"DeliverMax"`
		Destination        string `json:"Destination"`
		Fee                string `json:"Fee"`
		LastLedgerSequence int    `json:"LastLedgerSequence"`
		Memos              []struct {
			Memo struct {
				MemoData   string `json:"MemoData"`
				MemoType   string `json:"MemoType"`
				MemoFormat string `json:"MemoFormat"`
			} `json:"Memo"`
		} `json:"Memos"`
		Sequence        int    `json:"Sequence"`
		SigningPubKey   string `json:"SigningPubKey"`
		TransactionType string `json:"TransactionType"`
		TxnSignature    string `json:"TxnSignature"`
		Date            int    `json:"date"`
		Hash            string `json:"hash"`
		InLedger        int    `json:"inLedger"`
		LedgerIndex     int    `json:"ledger_index"`
		Topics          string `json:"topics"`
		LogData         []byte `json:"logData"`
	} `json:"tx"`
	Validated bool `json:"validated"`
}
