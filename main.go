package main

import (
	"telecom/spider"
	"log"
	"path/filepath"
	"os"
	"fmt"
	"flag"
)

var (
	h bool
	t string
	p string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&t, "t", "visit", "set `type`:spider,visit")
	flag.StringVar(&p, "p", "xici", "set `spider`:xici,kuai,cloud,en")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if t == "spider" {
		rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rootPath)
		var s spider.BaseSpiderInterface
		if p == "xici" {
			s = &spider.XiciSpider{}
		} else if p == "kuai" {
			s = &spider.KuaiSpider{}
		} else if p == "en" {
			s = &spider.EnSpider{}
		} else if p == "cloud" {
			s = &spider.CloudSpider{}
		} else if p == "ss" {
			s = &spider.SsSpider{}
		} else {
			flag.Usage()
			return
		}
		s.Init()
		s.Initialize()
		//s.SetPath(rootPath)
		s.Spider()
		s.CheckProxyPort()
		err = s.SaveProxy()
		if err != nil {
			log.Println(err)
		}
	} else if t == "visit" {
		rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		s := &spider.BaseSpider{}
		s.Init()
		s.Initialize()
		s.Path = rootPath
		s.Visit()
	} else {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `aaa`)
	flag.PrintDefaults()
}
