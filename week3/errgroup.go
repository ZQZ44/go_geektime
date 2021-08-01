package main

import (
  "os"
  "fmt"
  "syscall"
  "context"
  "net/http"
  "os/signal"
  "github.com/pkg/errors"
  "golang.org/x/sync/errgroup"
)

func main() {

    //set handler
    http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
        _,err := w.Write([]byte("hello world"))
        if err != nil {
            fmt.Println("writer err", err)
                return
        }
    })
    //create server
    server := http.Server{
        Addr:  ":8080",
        Handler:  http.DefaultServeMux,
    }

    //create errgroup
    g,ctx := errgroup.WithContext(context.Background())

    //start server
    g.Go(func() error {
        return server.ListenAndServe()
    })

    //listen to ctx.Done
    g.Go(func() error {
        select{
        case <-ctx.Done():
            fmt.Println("ctx is Done: ", ctx.Err())
        }
        fmt.Println("shutdown server")
        return server.Shutdown(ctx)
    })

    //listen to system signal
    g.Go(func() error {
        ch := make(chan os.Signal, 0)
        signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

        select {
        case <-ctx.Done():
            return ctx.Err()
        case sig := <-ch:
            return errors.Errorf("recv sig: %v", sig)
        }
    })

    //waiting errgroup finish
    err := g.Wait()
    if err != nil{
        fmt.Println("errgroup exit: ", err)
    }
}

