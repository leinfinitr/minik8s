package executor

import (
	"time"
)
type callback func()

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