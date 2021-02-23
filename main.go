package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getUptime(hostname string, username string, password string) float64 {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	method := "GET"
	url := fmt.Sprintf("http://%s/RgSwInfo.asp", hostname)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		html := string(bodyBytes)

		// <tr><td><b>System Up Time</b></td><td>0 days 03h:17m:08s</td></tr>
		re := regexp.MustCompile("(\\d+) days (\\d+)h:(\\d+)m:(\\d+)s")
		parts := re.FindStringSubmatch(html)

		// Parse the timestamp into seconds
		days, _ := strconv.ParseFloat(parts[1], 32)
		hours, _ := strconv.ParseFloat(parts[2], 32)
		minutes, _ := strconv.ParseFloat(parts[3], 32)
		seconds, _ := strconv.ParseFloat(parts[4], 32)

		var secs float64 = days*86400 + hours*3600 + minutes*60 + seconds
		return secs
	}
	return 0.0
}

var (
	namespace = "smartrg"

	// Metric - Duration of last scrape
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "uptime_seconds"),
		"The uptime, in seconds, of the SmartRG 808AC device.",
		[]string{"collector"}, nil,
	)
)

// Collect struct
type Collect struct {
	Up           bool
	ResponseTime string
}

type Exporter struct {
	hostname, username, password string
	collect                      Collect
	up                           prometheus.Gauge
	uptime                       prometheus.Gauge
}

// Describe implements prometheus.Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)

	ch <- e.up
	// ch <- e.errorDesc
	ch <- e.uptime
	// e.scrapeErrors.Collect(ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	uptime := getUptime(e.hostname, e.username, e.password)

	e.up.Set(1)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, uptime, "collect.uptime")
}

// New : Creates a new instance of Exporter for scraping metrics
func New(hostname string, username string, password string, collect Collect) *Exporter {
	return &Exporter{
		hostname: hostname,
		username: username,
		password: password,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "up",
			Help:      "Indicates if the monitor is up",
		}),

		uptime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "uptime",
			Help:      "The uptime, in seconds, of the SmartRG 808AC device.",
		}),
	}
}

func main() {
	hostname := os.Getenv("SMARTRG_HOSTNAME")
	username := os.Getenv("SMARTRG_USERNAME")
	password := os.Getenv("SMARTRG_PASSWORD")

	collect := Collect{}

	smartrgCollector := New(hostname, username, password, collect)
	prometheus.MustRegister(smartrgCollector)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on :2112")
	http.ListenAndServe(":2112", nil)
}
