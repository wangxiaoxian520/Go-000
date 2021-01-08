package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	g := new(errgroup.Group)
	sigs := make(chan os.Signal, 1)
	stop := make(chan struct{})

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	//signal
	g.Go(func() error {
		sig := <-sigs
		fmt.Println(sig)
		if sig == syscall.SIGINT {
			close(stop)
		}
		fmt.Println("信号监听结束.")
		return errors.New("sigal error")
	})
	server := http.Server{Addr: "127.0.0.1:8080"}
	g.Go(func() error {
		go func() {
			<-stop
			sigs <- syscall.SIGTERM
			fmt.Println("服务中断")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()
			err := server.Shutdown(ctx)
			fmt.Printf("服务关闭原因：%v\n", err)
		}()
		fmt.Println("服务开始")
		return server.ListenAndServe()
	})

	go func() {
		fmt.Println("延迟五秒")
		time.Sleep(time.Second * 5)
		close(stop)
	}()

	//上周作业修改
	// 	g, ctx := errgroup.WithContext(context.Background())
	// s := http.Server{Addr: "127.0.0.1:8080"}
	// g.Go(func() error {
	// 	g.Go(func() error {
	// 		<-ctx.Done()
	// 		fmt.Println("server closed")
	// 		return s.Shutdown(context.TODO()) //
	// 	})
	// 	return s.ListenAndServe()
	// })

	// sigs := make(chan os.Signal)
	// signal.Notify(sigs, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
	// g.Go(func() error {
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("ctx cancel")
	// 		return ctx.Err()
	// 	case <-sigs:
	// 		fmt.Println(<-sigs)
	// 		return errors.New("singal closed")
	// 	}

	// })
	if err := g.Wait(); err != nil {
		fmt.Printf("gorutine退出原因:%v\n", err)
	}
}
