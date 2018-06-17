package main

import (
	"encoding/hex"
	"flag"
	"fmt"
)

func main() {
	dbdir := flag.String("dbdir", "/root/.bitcoin/chainstate", "utxo or blockindex database dir")
	prefix := flag.Int("prefix", 67, "please input leveldb key prefix")
	obfuscate := flag.String("obfuscate", "", "leveldb obfuscate key used in utxo database")
	flag.Parse()

	dboption := &DBOption{
		FilePath:  *dbdir,
		CacheSize: 1 << 20,
	}
	if obfuscate != nil {
		dboption.DontObfuscate = true
	}

	dbw, err := NewDBWrapper(dboption)
	if err != nil {
		panic(err)
	}
	if obfuscate != nil {
		bs, err := hex.DecodeString(*obfuscate)
		if err != nil {
			panic(err)
		}
		dbw.obfuscateKey = bs
	}

	iter := dbw.Iterator()
	defer iter.Close()

	var utxoCount int
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		if *prefix == -1 {
			utxoCount++
		} else if int(iter.GetKey()[0]) == *prefix {
			utxoCount++
		}
	}

	fmt.Println("utxo set count:", utxoCount)
}
