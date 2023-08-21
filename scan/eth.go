package scan

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type EthScan struct {
	height uint64
	queue  *fifo
	client *EthClient
}

func NewEthScan() (*EthScan, error) {
	client, err := NewEthClient("https://goerli.infura.io/v3/6c5e7022939547d68f0825405b7f3186")
	if err != nil {
		return nil, err
	}
	scan := &EthScan{
		client: client,
		queue:  FifoNew(50),
	}
	return scan, nil
}

func (e *EthScan) Run(ctx context.Context) error {
	var err error
	if e.height == 0 {
		e.height, err = e.client.BlockNumber(ctx)
		if err != nil {
			return err
		}
	}
	for {
		block, err := e.client.BlockByNumber(ctx, int64(e.height))
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			e.dealBlock(ctx, block)
			e.queue.Print()
			e.height += 1
		}
	}
	return nil
}

func (e *EthScan) dealBlock(ctx context.Context, block *types.Block) {
	for _, tx := range block.Transactions() {
		totalFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
		v := GasValue{
			GasUsed:   totalFee,
			GasLimit:  tx.Gas(),
			GasTipCap: tx.GasTipCap().Int64(),
			GasFeeCap: tx.GasFeeCap().Int64(),
			GasBase:   block.Header().BaseFee.Uint64(),
			GasPrice:  tx.GasPrice().Int64(),
		}
		avg := e.queue.Avg()
		if avg > 0 && (v.GasLimit/avg) > 5 {
			if gas, err := e.client.TransactionReceipt(ctx, tx.Hash()); err == nil && gas > 0 {
				v.GasLimit = gas
				v.GasUsed = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(gas)))
			}
		}
		e.queue.Set(v)
	}
	gas := e.queue.Gas()
	fmt.Println("current gas", decimal.NewFromBigInt(gas, -18))
}
