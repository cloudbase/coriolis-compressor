package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cloudbase/coriolis-compressor/routers"
)

const (
	defaultBindAddress = "/var/run/coriolis-compressor.sock"
)

func getListener(host string) (net.Listener, error) {
	idx := strings.LastIndex(host, ":")
	if idx == -1 {
		return net.Listen("unix", host)
	}
	return net.Listen("tcp", host)
}

func main() {
	bindMsg := "Address to bind to. This can be a ip:port combination, or a unix socket"
	bind := flag.String("bind", defaultBindAddress, bindMsg)
	flag.Parse()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	listener, err := getListener(*bind)
	if err != nil {
		log.Fatal(err)
	}

	r := routers.GetRouter()
	srv := &http.Server{
		Handler: r,
	}

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	<-stop

	log.Print("shutting down http server")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Print(fmt.Sprintf("failed to shutdown web server: %q", err))
	}
}
