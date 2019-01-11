package spider

import (
	"sync"
	"net/http"
	"github.com/gocolly/colly"
	"log"
	"time"
	"strings"
	"strconv"
	"net"
	"path/filepath"
	"io/ioutil"
	"os"
	"github.com/gocolly/colly/proxy"
	"math/rand"
	"fmt"
)

type BaseSpider struct {
	Wg                sync.WaitGroup
	SpiderProxyList   map[string]string
	CheckOkProxy      []string
	Path              string
	FileName          string
	UserAgents        []string
	UserAgentFileName string
	MaxCoroutineNum   int
	Rate              int
	VisitWg           sync.WaitGroup
	Prefix            string
}

type BaseSpiderInterface interface {
	GetHeader(newHeader map[string]string) http.Header
	GrabDom(rUrl string, cookie string, selector string, cb func(e *colly.HTMLElement)) error
	NowTimeStr() string
	CheckProxyPort()
	Verify(proxyPort string, c *chan int)
	SaveProxy() error
	Read(name string) []string
	Init()
	Initialize()
	Spider()
	Visit()
	SetPath(path string)
}

func (s *BaseSpider) GetHeader(newHeader map[string]string) http.Header {
	hdr := http.Header{}
	header := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36",
	}
	for i, v := range newHeader {
		header[i] = v
	}
	for i, v := range header {
		hdr.Set(i, v)
	}
	return hdr
}

func (s *BaseSpider) GrabDom(rUrl string, cookie string, selector string, cb func(e *colly.HTMLElement)) error {
	c := colly.NewCollector()
	header := map[string]string{}
	header["Cookie"] = cookie
	hdr := s.GetHeader(header)
	c.OnResponse(func(r *colly.Response) {

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println(s.NowTimeStr(), " -- request -- ", r.URL.String())
	})

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		cb(e)
	})
	err := c.Request("GET", rUrl, nil, nil, hdr)
	if err != nil {
		log.Fatal(s.NowTimeStr(), " -- request -- err -- ", err)
	}
	return nil
}

func (s *BaseSpider) NowTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (s *BaseSpider) CheckProxyPort() {
	preCheck := []string{}
	for key, _ := range s.SpiderProxyList {
		preCheck = append(preCheck, key)
	}
	length := len(preCheck)
	s.Wg = sync.WaitGroup{}
	s.Wg.Add(length)
	coroutineNum := s.MaxCoroutineNum

	if length <= coroutineNum {
		coroutineNum = length / s.Rate
	}
	coroutineNum = 100
	preChan := make(chan int, coroutineNum)
	for _, val := range preCheck {
		preChan <- 1
		go s.Verify(val, &preChan)
	}
	s.Wg.Wait()
}

func (s *BaseSpider) Verify(proxyPort string, c *chan int) {
	defer s.Wg.Done()
	t := strings.Split(proxyPort, ":")
	fmt.Println(t)
	port, _ := strconv.Atoi(t[1])
	tcpAddr := net.TCPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if err == nil {
		s.CheckOkProxy = append(s.CheckOkProxy, proxyPort)
		conn.Close()
	}
	<-*c
}

func (s *BaseSpider) SaveProxy() error {
	data := strings.Join(s.CheckOkProxy, "\n")
	fileName := filepath.Join(s.Path, s.Prefix+s.FileName)
	err := ioutil.WriteFile(fileName, []byte(data), 0644)
	return err
}

func (s *BaseSpider) Read(name string) []string {
	fileName := filepath.Join(s.Path, name)
	fp, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(data), "\n")
}

func (s *BaseSpider) Init() {
	s.FileName = "proxy.txt"
	s.UserAgentFileName = "useragent.txt"
	s.UserAgents = s.Read(s.UserAgentFileName)
	s.SpiderProxyList = map[string]string{}
	s.CheckOkProxy = []string{}
	s.MaxCoroutineNum = 50
	s.Rate = 3
}

func (s *BaseSpider) Initialize() {
	s.Prefix = ""
}

func (s *BaseSpider) Spider() {
}

func (s *BaseSpider) Visit() {
	lenUserAgent := len(s.UserAgents)
	if lenUserAgent <= 1 {
		log.Println("user agent is empty")
		return
	}
	l, err := filepath.Glob(filepath.Join(s.Path, "*proxy.txt"))
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range l {
		t := s.Read(filepath.Base(f))
		s.CheckOkProxy = append(s.CheckOkProxy, t...)
	}
	lenProxy := len(s.CheckOkProxy)
	if lenProxy < 1 {
		log.Println("proxy is empty")
		return
	}
	for start := 0; start < lenProxy; start += 50 {
		high := start + 50
		if high >= lenProxy {
			high = lenProxy - 1
		}
		s.VisitWg = sync.WaitGroup{}
		dd := s.CheckOkProxy[start:high]
		us := []string{}
		vchan := make(chan int, 50)
		for _, v := range dd {
			if lenUserAgent < 4 {
				us = s.UserAgents
			} else {
				rand.Seed(time.Now().UnixNano())
				max := lenUserAgent - 3
				ts := rand.Intn(max) + 0
				//fmt.Println(ts)
				//fmt.Println(lenUserAgent)
				us = s.UserAgents[ts:ts+2]
			}
			for _, ua := range us {
				s.VisitWg.Add(1)
				vchan <- 1
				go s.DoVisit(ua, &vchan, "http://"+v)
			}
		}
		s.VisitWg.Wait()
	}
	s.VisitWg = sync.WaitGroup{}
	s.VisitWg.Wait()
}

func (s *BaseSpider) DoVisit(userAgent string, ch *chan int, proxyUrl ...string) {
	defer s.VisitWg.Done()
	c := colly.NewCollector()
	c.UserAgent = userAgent
	rp, err := proxy.RoundRobinProxySwitcher(proxyUrl...)
	if err != nil {
		log.Println(err)
		return
	}
	c.SetProxyFunc(rp)
	c.OnResponse(func(r *colly.Response) {
		log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
		log.Printf("User Agent %s\n", r.Request.Headers.Get("User-Agent"))
	})
	c.Visit("http://deal.189store.com/getvip/?channel=b1HMcQG0fceIigI3")
	<-*ch
}

func (s *BaseSpider) SetPath(path string) {
	s.Path = path
}
