package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"soapi"

	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

const appBanner = `
███████╗ ██████╗ ██████╗ ██╗██╗  ██╗███████╗
██╔════╝██╔═══██╗██╔══██╗██║██║ ██╔╝██╔════╝
███████╗██║   ██║██████╔╝██║█████╔╝ █████╗  
╚════██║██║   ██║██╔══██╗██║██╔═██╗ ██╔══╝  
███████║╚██████╔╝██████╔╝██║██║  ██╗███████╗
╚══════╝ ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝

`

func runserver(port string) {
	router := httprouter.New()
	router.GET("/", soapi.RedirectHome)
	router.GET("/api/", soapi.Index)
	router.GET("/api/address/:addr", soapi.Address)
	router.GET("/api/bike/:lng/:lat", soapi.NearbyBike)
	router.ServeFiles("/web/*filepath", http.Dir("src/webpage/"))

	log.WithFields(log.Fields{
		"Port": port,
	}).Info("Listen server on: 0.0.0.0")

	http.ListenAndServe(":"+port, router)
}

func displayHelp() {
	fmt.Fprintf(os.Stderr, `sobike - 搜索身边的共享单车
		
Version: v0.1

Usage: sobike-web [-h] [-p port]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	var port string
	var help bool

	flag.BoolVar(&help, "h", false, "print this help")
	flag.StringVar(&port, "p", "8080", "port")
	flag.Usage = displayHelp

	flag.Parse()

	color.HiCyan(appBanner)
	if help {
		flag.Usage()
	} else {
		runserver(port)
	}
}
