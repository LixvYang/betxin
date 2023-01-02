package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/lixvyang/betxin/internal/model"
	"github.com/lixvyang/betxin/internal/router"
	"github.com/lixvyang/betxin/internal/service"
	"github.com/lixvyang/betxin/internal/service/dailycurrency"
	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/cron"

	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
)

func main() {
	signalch := make(chan os.Signal, 1)
	utils.InitSetting.Do(utils.Init)
	service.InitMixin.Do(service.InitMixinClient)
	ctx := context.Background()
	model.InitDb()
	betxinredis.NewRedisClient(ctx)
	go cron.HeathCheck()
	go dailycurrency.DailyCurrency(ctx)
	go service.Worker(ctx)
	go router.InitRouter(signalch)
	// tm := time.Tick(time.Second)
	// go http.ListenAndServe("0.0.0.0:8888", nil)

	// for range tm {
	// 	go fmt.Println(runtime.NumGoroutine())
	// }
	// for {
	// 	select {
	// 	case <-tm:
	// 		go fmt.Println(runtime.NumGoroutine())
	// 	}
	// }

	//attach signal
	signal.Notify(signalch, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalch
	signal.Stop(signalch)
	//cleanup before exit
	log.Printf("On Signal <%s>", signalType)
	log.Println("Exit command received. Exiting...")
}
