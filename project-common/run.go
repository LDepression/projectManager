/**
 * @Author: lenovo
 * @Description:
 * @File:  run
 * @Version: 1.0.0
 * @Date: 2023/07/11 15:22
 */

package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName, addr string, stop func()) {
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		log.Printf("%s listening on %s", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//syscall.SIGINT 也就是kill -2
	//syscall.SIGTERM 也就是结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutting Down project %s...\n", srvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("%s shutdone ,cause by :%v\n", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout")
	}
	log.Printf("%s stop success...\n", srvName)
}
