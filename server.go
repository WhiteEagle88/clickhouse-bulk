package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server - main server object
type Server struct {
	Listen    string
	Collector *Collector
	Debug     bool
	echo      *echo.Echo
}

// Status - response status struct
type Status struct {
	Status    string                       `json:"status"`
	SendQueue int                          `json:"send_queue,omitempty"`
	Servers   map[string]*ClickhouseServer `json:"servers,omitempty"`
	Tables    map[string]*Table            `json:"tables,omitempty"`
}

// NewServer - create server
func NewServer(listen string, collector *Collector, debug bool) *Server {
	return &Server{listen, collector, debug, echo.New()}
}

func (server *Server) readHandler(c echo.Context) error {
	s := ""

	qs := getQueryStringWithAuthentication(c)
	params, content, insert := server.Collector.ParseQuery(qs, s)
	if insert {
		go server.Collector.Push(params, content)
		return c.String(http.StatusOK, "")
	}
	resp, status, _ := server.Collector.Sender.SendQuery(&ClickhouseRequest{Params: qs, Content: s})
	return c.String(status, resp)
}

func (server *Server) writeHandler(c echo.Context) error {
	q, _ := ioutil.ReadAll(c.Request().Body)
	s := string(q)

	qs := getQueryStringWithAuthentication(c)
	params, content, insert := server.Collector.ParseQuery(qs, s)
	if insert {
		go server.Collector.Push(params, content)
		return c.String(http.StatusOK, "")
	}
	resp, status, _ := server.Collector.Sender.SendQuery(&ClickhouseRequest{Params: qs, Content: s})
	return c.String(status, resp)
}

func (server *Server) statusHandler(c echo.Context) error {
	return c.JSON(200, Status{Status: "ok"})
}

// Start - start http server
func (server *Server) Start() error {
	return server.echo.Start(server.Listen)
}

// Shutdown - stop http server
func (server *Server) Shutdown(ctx context.Context) error {
	return server.echo.Shutdown(ctx)
}

// InitServer - run server
func InitServer(listen string, collector *Collector, debug bool) *Server {
	server := NewServer(listen, collector, debug)
	server.echo.GET("/", server.readHandler)
	server.echo.POST("/", server.writeHandler)
	server.echo.GET("/status", server.statusHandler)
	server.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return server
}

// SafeQuit - safe prepare to quit
func SafeQuit(collect *Collector, sender Sender) {
	collect.FlushAll()
	if count := sender.Len(); count > 0 {
		log.Printf("Sending %+v tables\n", count)
	}
	for !sender.Empty() && !collect.Empty() {
		collect.WaitFlush()
	}
	collect.WaitFlush()
}

// RunServer - run all
func RunServer(cnf Config) {
	InitMetrics()
	dumper := NewDumper(cnf.DumpDir)
	//dumperDebug := NewDumperDebug("debug")
	sender := NewClickhouse(cnf.Clickhouse.DownTimeout, cnf.Clickhouse.ConnectTimeout, cnf.Debug)
	sender.Dumper = dumper
	//sender.DumperDebug = dumperDebug
	for _, url := range cnf.Clickhouse.Servers {
		sender.AddServer(url)
	}

	collect := NewCollector(sender, cnf.FlushCount, cnf.FlushInterval, cnf.Cache)

	// send collected data on SIGTERM and exit
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	srv := InitServer(cnf.Listen, collect, cnf.Debug)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		for {
			_ = <-signals
			log.Printf("STOP signal\n")
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Shutdown error %+v\n", err)
				SafeQuit(collect, sender)
				os.Exit(1)
			}
		}
	}()

	if cnf.DumpCheckInterval >= 0 {
		dumper.Listen(sender, cnf.DumpCheckInterval)
	}

	err := srv.Start()
	if err != nil {
		log.Printf("ListenAndServe: %+v\n", err)
		SafeQuit(collect, sender)
		os.Exit(1)
	}
}

func getQueryStringWithAuthentication(c echo.Context) string {
	qs := c.QueryString()
	user, password, ok := c.Request().BasicAuth()
	userFromHeader := url.QueryEscape(c.Request().Header.Get("X-ClickHouse-User"))
	passwordFromHeader := url.QueryEscape(c.Request().Header.Get("X-ClickHouse-Key"))
	if ok {
		if qs == "" {
			qs = "user=" + user + "&password=" + password
		} else {
			qs = "user=" + user + "&password=" + password + "&" + qs
		}
	} else if userFromHeader != "" || passwordFromHeader != "" {
		if qs == "" {
			qs = "user=" + userFromHeader + "&password=" + passwordFromHeader
		} else {
			qs = "user=" + userFromHeader + "&password=" + passwordFromHeader + "&" + qs
		}
	}

	return qs
}
