// Copyright 2021 Compass Systems
// SPDX-License-Identifier: LGPL-3.0-only

package blockstore

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "blockstore")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bs, err := New(dir, "10")
	if err != nil {
		t.Fatal(err)
	}
	// Load non-existent dir/file
	block, err := bs.TryLoadLatestBlock()
	if err != nil {
		t.Fatal(err)
	}

	if block.Uint64() != uint64(0) {
		t.Fatalf("Expected: %d got: %d", 0, block.Uint64())
	}

	// Mos block number
	block = big.NewInt(999)
	err = bs.StoreBlock(block)
	if err != nil {
		t.Fatal(err)
	}

	// Load block number
	latest, err := bs.TryLoadLatestBlock()
	if err != nil {
		t.Fatal(err)
	}

	if block.Uint64() != latest.Uint64() {
		t.Fatalf("Expected: %d got: %d", block.Uint64(), latest.Uint64())
	}

	// Mos block number again
	block = big.NewInt(1234)
	err = bs.StoreBlock(block)
	if err != nil {
		t.Fatal(err)
	}

	// Load block number
	latest, err = bs.TryLoadLatestBlock()
	if err != nil {
		t.Fatal(err)
	}

	if block.Uint64() != latest.Uint64() {
		t.Fatalf("Expected: %d got: %d", block.Uint64(), latest.Uint64())
	}
}
