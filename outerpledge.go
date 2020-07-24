package main

import (
	"fmt"
	"os/exec"
	"time"
)
type tefunc func() (error)
func main() {
	var interval,count,batch,total int
	fmt.Printf("Please enter your config 'interval count batch total' like '10 5 60 200': ")
	fmt.Scanln(&interval, &count,&batch,&total)
	fmt.Println("interval=d%,count=d%,batch=d%,total=d%",interval,count,batch,total)

	schedle(interval,count,batch,total,func()(error){
		cmd := exec.Command("lotus-miner sectors pledge")
		buf, err :=cmd.Output()
		fmt.Println("任务执行结果返回：s%",buf)
		return err
	})

}
func schedle(interval,count,batch,total int,cc tefunc) error{
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