package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(ctx)
	srv := &http.Server{Addr: ":8080"}
	srv2 := &http.Server{Addr: ":8082"}
	g.Go(func() error {
		return StartHttpServer(srv)
	})
	g.Go(func() error {
		return StartHttpServer2(srv2)
	})

	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	g.Go(func() error {
		<-chanel
		fmt.Println("-----chanel-----")
		cancel()
		fmt.Println(srv.Shutdown(errCtx))
		fmt.Println(srv2.Shutdown(errCtx))
		return errors.New("kill signal ")

	})
	if err := g.Wait(); err != nil {
		fmt.Println("group error: ", err)
	}
	fmt.Println("all group done!")
}

func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/getA", getA)
	fmt.Println("http 8080 server start")
	err := srv.ListenAndServe()
	return err
}

func StartHttpServer2(srv *http.Server) error {
	http.HandleFunc("/getB", getB)
	fmt.Println("http 8082 server start")
	err := srv.ListenAndServe()
	return err
}

func getA(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("p1") == "panic" {
		panic("test quit")
	}
	w.Write([]byte("hello getA http server"))
}

func getB(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello getB http server"))
}
