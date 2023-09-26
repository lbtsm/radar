// Copyright 2021 Compass Systems
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type Conner interface {
	Client() *ethclient.Client
}

type Connection struct {
	endpoint string
	http     bool
	conn     *ethclient.Client
}

// NewConn returns an uninitialized connection, must call Connection.Connect() before using.
func NewConn(endpoint string) *Connection {
	return &Connection{
		endpoint: endpoint,
		http:     true,
	}
}

// Connect starts the ethereum WS connection
func (c *Connection) Connect() error {
	var (
		err       error
		rpcClient *rpc.Client
	)
	rpcClient, err = rpc.DialHTTP(c.endpoint)
	if err != nil {
		return err
	}
	c.conn = ethclient.NewClient(rpcClient)

	return nil
}

func (c *Connection) Client() *ethclient.Client {
	return c.conn
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (*big.Int, error) {
	bnum, err := c.conn.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetUint64(bnum), nil
}

// Close terminates the client connection and stops any running routines
func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
