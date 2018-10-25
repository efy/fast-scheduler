package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/ddo/go-fast.v0"
)

var gauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "internet_download_speed",
	Help: "download speed measured in Kbps using the fast.com service",
})

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go measure()

	http.ListenAndServe(":2112", nil)
}

func measure() {
	for {
		f := fast.New()

		err := f.Init()
		if err != nil {
			log.Fatal(err)
		}

		urls, err := f.GetUrls()
		if err != nil {
			log.Fatal(err)
		}

		kbps := make(chan float64)

		go func() {
			var kk float64
			var kc int
			for kb := range kbps {
				log.Printf("kbps: %.2f", kb)
				kk += kb
				kc++
			}

			avg := kk / float64(kc)
			log.Println(avg, "Kbps (avg)")
			gauge.Set(avg)
		}()

		err = f.Measure(urls, kbps)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(30 * time.Second)
	}
}
