package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

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
		fmt.Println(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	request.URL.Path = strings.TrimPrefix(request.URL.Path, proxyItem.Prefix)
	proxy.ServeHTTP(writer, request)
}

func main() {
	var cfg Config

	d, err := os.ReadFile("http_proxy_config.yaml")
	if err != nil {
		panic("no config file:" + err.Error())
	}

	err = yaml.Unmarshal(d, &cfg)
	if err != nil {
		panic("invalid config file:" + err.Error())
	}

	if cfg.StaticFs == "" {
		panic("no static fs config")
	}

	r := mux.NewRouter()

	for idx := 0; idx < len(cfg.ProxyItems); idx++ {
		item := &cfg.ProxyItems[idx]

		r.PathPrefix(item.Prefix).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			proxyIt(item, writer, request)
		})
	}

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticFs)))

	log.Print("Listening on :" + cfg.Listen + "...")
	err = http.ListenAndServe(cfg.Listen, r)
	if err != nil {
		log.Fatal(err)
	}
}
