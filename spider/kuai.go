package spider

import (
	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"log"
)

type KuaiSpider struct {
	BaseSpider
}

func (s *KuaiSpider) Initialize() {
	s.Prefix = "kuai_"
}

func (s *KuaiSpider) Spider() {
	//urls := []string{"https://www.xicidaili.com/nn/1", "https://www.xicidaili.com/nn/2", "https://www.xicidaili.com/nn/3", "https://www.xicidaili.com/nn/4"}
	format := "https://www.kuaidaili.com/free/inha/%d/"
	for page := 1; page <= 10; page++ {
		u := fmt.Sprintf(format, page)
		log.Println("sleep 5 second")
		time.Sleep(5 * time.Second)
		s.GrabDom(u, "", "#list>table", func(e *colly.HTMLElement) {
			e.DOM.Find("tr").Each(func(i int, qs *goquery.Selection) {
				if i == 0 {
					return
				}
				row := make([]string, 7)
				qs.Find("td").Each(func(j int, td *goquery.Selection) {
					row[j] = td.Text()
				})
				if row[3] != "HTTP" {
					return
				}
				t := row[0] + ":" + row[1]
				s.SpiderProxyList[t] = t
				return
			})
		})
	}
}
