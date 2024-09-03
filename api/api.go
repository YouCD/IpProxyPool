package api

import (
	"IpProxyPool/middleware/config"
	"IpProxyPool/middleware/database"
	"IpProxyPool/middleware/storage"
	"IpProxyPool/util/iputil"
	"context"
	"encoding/json"
	"fmt"
	"github.com/youcd/toolkit/log"
	"net/http"
	"time"
)

// Run for request
func Run(setting *config.System) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/all", ProxyAllHandler)
	mux.HandleFunc("/http", ProxyHttpHandler)
	mux.HandleFunc("/https", ProxyHttpsHandler)
	mux.HandleFunc("/count", CountHandler)
	mux.HandleFunc("/del", ProxyDelHandler)
	server := &http.Server{
		Addr:           setting.HttpAddr + ":" + setting.HttpPort,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Infof("Server run at:")
	log.Infof("- Local:   http://localhost:%s ", setting.HttpPort)
	log.Infof("- Network: http://%s:%s ", iputil.GetLocalHost(), setting.HttpPort)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic("listen: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	errs := server.Shutdown(ctx)
	if errs != nil {
		log.Info("Server Shutdown:", errs)
		fmt.Println("Server Shutdown:", errs)
	}

	log.Info("Server exiting")
}

func ProxyDelHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		database.DeleteByIp(request.URL.Query().Get("ip"))
		writer.Write([]byte("ok"))
	}
}

func CountHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		writer.Header().Set("content-type", "application/json")
		b, _ := json.Marshal(database.Count())
		writer.Write(b)
	}
}

// IndexHandler .
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		apiMap := make(map[string]string, 0)
		apiMap["/"] = "api 指引"
		apiMap["/all"] = "获取随机的一个 http 或 https 类型的代理IP"
		apiMap["/http"] = "获取随机的一个 http 类型的代理IP"
		apiMap["/https"] = "获取随机的一个 https 类型的代理IP"
		apiMap["/count"] = "统计信息"
		apiMap["/del"] = "删除代理IP"
		b, _ := json.Marshal(apiMap)
		w.Write(b)
	}
}

// ProxyAllHandler .
func ProxyAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(storage.RandomProxy())
		if err != nil {
			return
		}
		w.Write(b)
	}
}

// ProxyHttpHandler .
func ProxyHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(storage.RandomByProxyType("http"))
		if err != nil {
			return
		}
		log.Debug("get http proxy: ", string(b))
		w.Write(b)
	}
}

// ProxyHttpsHandler .
func ProxyHttpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(storage.RandomByProxyType("https"))
		if err != nil {
			return
		}
		log.Debug("get https proxy: ", string(b))
		w.Write(b)
	}
}
