package executor

import (
	"time"
)
type callback func()
type condCallback func() bool

func ExecuteInPeriod(delay time.Duration,timeGap []time.Duration,callback callback){
	if len(timeGap) == 0 {
		return
	}
	<-time.After(delay)
	for {
		for _,gap := range timeGap {
			<-time.After(gap)
			callback()
		}
	}
}

func CondExecuteInPeriod(delay time.Duration,timeGap []time.Duration,callback condCallback){
	if len(timeGap) == 0 {
		return
	}
	<-time.After(delay)
	for {
		for _,gap := range timeGap {
			<-time.After(gap)
			if callback() {
				return;
			}
		}
	}
}