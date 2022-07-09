package postgres

import "github.com/brsuite/btcwallet/walletdb"

type Fixture interface {
	DB() walletdb.DB
	Dump() (map[string]interface{}, error)
}
