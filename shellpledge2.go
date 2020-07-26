package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"
)
type tefunc func() (error)
func main() {
	// shellpledge --interval=10 --count=5 --batch=60 --total=200 --start=20
	// seconds
	interval := flag.Int("interval", 10, "value of interval")
	count := flag.Int("count", 5, "value of count")
	batch := flag.Int("batch", 60, "value of batch")
	total := flag.Int("total", 100, "value of total")
	start := flag.Int("start", 100, "value of start")
	flag.Parse()
	log.Println(fmt.Sprintf("interval=%d,count=%d,batch=%d,total=%d,start=%d",*interval,*count,*batch,*total,*start))
	log.Println("任务将在",*start,"分钟后开始执行")
	time.Sleep(time.Second * 60 * time.Duration(*start))

	schedle(*interval,*count,*batch,*total,func()(error){
		cmd := exec.Command("/usr/local/bin/lotus-miner","sectors","pledge")

		buf, err :=cmd.Output()
		log.Println(fmt.Sprintf("任务执行结果返回：%s",buf))
		return err
	})

}
func schedle(interval,count,batch,total int,cc tefunc) {
	totalCount:=1
	for {
		for j := 0; j < count; j++ {
			var err = cc()
			if err!=nil{
				log.Println("任务执行出错:",err)
			}
			log.Println(fmt.Sprintf("任务执行第%d次",totalCount))
			totalCount++
			if (totalCount-total > 0){
				return
			}
			time.Sleep(time.Second * time.Duration(interval))
		}
		time.Sleep(time.Second * time.Duration(batch))
	}
}