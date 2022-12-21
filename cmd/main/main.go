package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lixvyang/betxin/internal/router"
	"github.com/lixvyang/betxin/internal/service"
	"github.com/lixvyang/betxin/internal/service/dailycurrency"
	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/model"

	betxinredis "github.com/lixvyang/betxin/internal/utils/redis"
)

func main() {
	signalch := make(chan os.Signal, 1)
	utils.InitSetting.Do(utils.Init)
	service.InitMixin.Do(service.InitMixinClient)
	ctx := context.Background()
	model.InitDb()
	betxinredis.NewRedisClient(ctx)
	go dailycurrency.DailyCurrency(ctx)
	go service.Worker(ctx)
	go router.InitRouter(signalch)

	//attach signal
	signal.Notify(signalch, os.Interrupt, syscall.SIGTERM)
	signalType := <-signalch
	signal.Stop(signalch)
	//cleanup before exit
	log.Printf("On Signal <%s>", signalType)
	log.Println("Exit command received. Exiting...")
}
