package scheduler

import (
	"fmt"
	"github.com/carlescere/scheduler"
)

func InitScheduler() {
	_, _ = scheduler.Every(30).Minutes().Run(schedulerTest)
}

func schedulerTest() {
	fmt.Println("scheduler test")
}