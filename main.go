package main

import (
	"block-scan/scan"
	"context"
)

func main() {
	scan, _ := scan.NewEthScan()
	scan.Run(context.Background())

	//fmt.Println(123123)
	//for {
	//	ethscan, _ := scan.NewEthScan()
	//	err := ethscan.Run(context.Background())
	//	fmt.Println("-----", err)
	//	time.Sleep(time.Second * 5)
	//}
}
