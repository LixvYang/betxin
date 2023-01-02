package cron

import (
	"fmt"
	"time"

	"github.com/lixvyang/betxin/internal/utils"
	"github.com/lixvyang/betxin/internal/utils/email"
	"github.com/lixvyang/betxin/pkg/timewheel"
	"github.com/shirou/gopsutil/mem"
)

var usedPercent int

func HeathCheck() {
	timewheel.Every(10*time.Second, func() {
		// 检查内存使用情况
		u, _ := mem.VirtualMemory()
		usedPercent = int(u.UsedPercent)
		if utils.AppMode == "release" {
			if usedPercent > 80 {
				go email.NotifyHandler(fmt.Sprintf("服务器内存使用过高: %d", usedPercent))
			}
		}
	})
}
