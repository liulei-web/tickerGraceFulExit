# tickerGraceFulExit

```sh
go get -u github.com/liulei-web/tickerGraceFulExit
```

```go
package main

import (
	"context"
	"fmt"
	"github.com/liulei-web/tickerGraceFulExit"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var job = tickerGraceFulExit.TickerManage(tickerGraceFulExit.SetMaxTask(100))

func handle(ctx context.Context, goroutineNo int, sec int64) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	for {
		select {
		case <-ticker.C:
			go func(t int) {
				if err := job.Add(); err == nil {
					fmt.Printf("goroutine %v handle start \n", t)
					time.Sleep(3 * time.Second)
					fmt.Printf("goroutine %v handle end  \n", t)
					job.Done()
				} else {
					ticker.Stop()
					return
				}
			}(goroutineNo)
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("context done")
			return
		}
	}

}

func main() {

	ctx := context.Background()
	ctxCancel, cancelFun := context.WithTimeout(ctx, 15*time.Second)

	go func() {

		for i := 1; i <= 100; i++ {
			go handle(ctxCancel, i, 1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c

	fmt.Println("wait all goroutine end")

	job.WaitGraceFulExit()
	cancelFun()
	fmt.Println("end")
	time.Sleep(5 * time.Second)
	fmt.Println("exit")

}

```
