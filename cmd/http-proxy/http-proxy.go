package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type ProxyItem struct {
	Prefix string `yaml:"Prefix"`
	Target string `yaml:"Target"`
}

type Config struct {
	Listen     string      `yaml:"Listen"`
	StaticFs   string      `yaml:"StaticFs"`
	ProxyItems []ProxyItem `yaml:"ProxyItems"`
}

func proxyIt(proxyItem *ProxyItem, writer http.ResponseWriter, request *http.Request) {
	remote, err := url.Parse(proxyItem.Target)
	if err != nil {
		log.Println("parse", proxyItem.Target, "failed:", err)

		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	request.URL.Path = strings.TrimPrefix(request.URL.Path, proxyItem.Prefix)
	proxy.ServeHTTP(writer, request)
}

func main() {
	var cfg Config

	d, err := os.ReadFile("http_proxy_config.yaml")
	if err != nil {
		log.Panicln("no config file:", err.Error())
	}

	err = yaml.Unmarshal(d, &cfg)
	if err != nil {
		log.Panicln("invalid config file:", err.Error())
	}

	if cfg.StaticFs == "" {
		log.Panicln("no static fs config")
	}

	r := mux.NewRouter()

	for idx := 0; idx < len(cfg.ProxyItems); idx++ {
		item := &cfg.ProxyItems[idx]

		r.PathPrefix(item.Prefix).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			proxyIt(item, writer, request)
		})

		log.Println("proxy item: " + item.Prefix + " => " + item.Target)
	}

	log.Println("dest dir:", cfg.StaticFs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticFs)))

	log.Println("Listening on :" + cfg.Listen + "...")

	server := &http.Server{
		Addr:        cfg.Listen,
		Handler:     r,
		ReadTimeout: time.Second,
	}
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
