package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gopkg.in/alecthomas/kingpin.v2"

	sxgeo "github.com/tenrok/go-sxgeo"
)

const VERSION = "1.0.0"

var (
	KpApp       = kingpin.New("ip2geo", "A tool with HTTP server to get the geodata by IP.")
	KpFlgDbPath = KpApp.Flag("dbpath", "Path to SxGeoCity.dat file.").Short('d').Default("./SxGeoCity.dat").String()
	KpCmdFind   = KpApp.Command("find", "Find IP address in database.")
	KpArgIP     = KpCmdFind.Arg("ip", "IP address.").Required().String()
	KpCmdServe  = KpApp.Command("serve", "Start HTTP server, listen and serve.")
	KpFlgAddr   = KpCmdServe.Flag("addr", "Address to listen for HTTP requests on.").Short('a').Default(":8080").String()
)

var Geo sxgeo.SxGEO

// Use a single instance of Validate
var Validate = validator.New()

func main() {
	KpApp.HelpFlag.Short('h')
	KpApp.Version(VERSION).VersionFlag.Short('v')

	// Parse command line
	command := kingpin.MustParse(KpApp.Parse(os.Args[1:]))

	Geo = sxgeo.New(*KpFlgDbPath)

	switch command {
	case KpCmdFind.FullCommand():
		cmdFind()
	case KpCmdServe.FullCommand():
		cmdServe()
	}
}

func cmdFind() {
	if err := Validate.Var(*KpArgIP, "ipv4"); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	city, err := Geo.GetCityFull(*KpArgIP)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	j, err := json.MarshalIndent(city, "", "\t")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println(string(j))
}

func cmdServe() {
	router := gin.Default()
	router.GET("/sxgeo/:ip", func(c *gin.Context) {
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		ip := c.Param("ip")
		if err := Validate.Var(ip, "required,ipv4"); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Give me an IP, please"})
			return
		}
		city, err := Geo.GetCityFull(ip)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, city)
	})

	srv := &http.Server{Addr: *KpFlgAddr, Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for range interrupt {
			log.Println("Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
