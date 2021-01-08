package main

import (
	"Go-000/Week04/internal/service"
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		l, err := net.Listen("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		//当发生报错时，退出
		g.Go(func() error {
			<-ctx.Done()
			s.GracefulStop()
			return ctx.Err()
		})
		service.RegisterAPI(s)
		err = s.Serve(l)
		if err != nil {
			log.Fatalf("failed to server: %v", err)
		}
		return
	})
	// syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt
	ch := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt}
	signal.Notify(ch, sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ch:
				return errors.New("singal err")
			}
		}
	})
	if err := g.Wait(); err != nil {
		log.Fatalf("quit reason: %v", err)
	}
}
