package scan

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)


type EthClient struct {
	client *ethclient.Client
}

func NewEthClient(url string) (*EthClient, error) {
	scan := &EthClient{}
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)	
		return nil, err
	}
	scan.client = client
	return scan, nil
}

func (e *EthClient) BlockByNumber(ctx context.Context, height int64) (*types.Block, error) {
	block, err := e.client.BlockByNumber(ctx, big.NewInt(height))
	if err != nil {
		return nil, err
	}
	return block, err
}

func (e *EthClient) BlockNumber(ctx context.Context) (uint64, error) {
	height, err := e.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return height, err
}

func (e *EthClient) TransactionReceipt(ctx context.Context, hash common.Hash) (uint64, error) {
	receipt, err := e.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return 0, err
	}
	return receipt.GasUsed, err
}


