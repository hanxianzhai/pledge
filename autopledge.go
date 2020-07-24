package main

import (
	"fmt"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/urfave/cli/v2"
	"time"
)
// lotus-miner auto pledge --interval=10 --count=5 --batch=60 --total=200
// interval任务间隔时间，count任务数，batch批次数，total总任务数
type taskFunc func() error
var sectorAutoCmd = &cli.Command{
	Name:  "auto",
	Usage: "automatic execute sectors pledge",
	Subcommands: []*cli.Command{
		sectorsAutoPledgeCmd,
	},
}

var sectorsAutoPledgeCmd = &cli.Command{
	Name:  "pledge",
	Usage: "taskinterval taskcount batchinterval totaltasknumber",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "interval",
			Usage: "task interval:minutes",
			Value: 10,
		},
		&cli.IntFlag{
			Name:  "count",
			Usage: "the count task in a batch pledge",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "batch",
			Usage: "batch interval:minutes",
			Value: 60,
		},
		&cli.IntFlag{
			Name:  "total",
			Usage: "total task counts",
			Value: 200,
		},
	},
	Action: func(cctx *cli.Context) error {
		nodeApi, closer, err := lcli.GetStorageMinerAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := lcli.ReqContext(cctx)

		interval,count,batch,total:=getParam(cctx)
		sche(interval,count,batch,total,func()error{
			return nodeApi.PledgeSector(ctx)
		})

		return nil
	},
}
func getParam(cctx *cli.Context)(int,int,int,int){
	return cctx.Int("interval"),cctx.Int("count"),cctx.Int("batch"),cctx.Int("total")
}
func sche(interval,count,batch,total int,cc taskFunc) error{
	totalCount:=1
	for {
		for j := 0; j < count; j++ {
			var err = cc()
			if err!=nil{
				fmt.Println("任务执行出错:s%",err)
			}
			fmt.Println("任务执行第d%次",totalCount)
			totalCount++
			time.Sleep(time.Second * 60* time.Duration(interval))
		}
		if (totalCount-total > 0){
			break
		}
		time.Sleep(time.Second * 60 *time.Duration(batch))
	}
}