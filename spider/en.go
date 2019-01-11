package spider

import (
	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"log"
	"strings"
)

type EnSpider struct {
	BaseSpider
}

func (s *EnSpider) Initialize() {
	s.Prefix = "en_"
}

func (s *EnSpider) Spider() {
	//urls := []string{"https://www.xicidaili.com/nn/1", "https://www.xicidaili.com/nn/2", "https://www.xicidaili.com/nn/3", "https://www.xicidaili.com/nn/4"}
	format := "http://www.89ip.cn/index_%d.html"
	for page := 1; page <= 10; page++ {
		u := fmt.Sprintf(format, page)
		log.Println("sleep 5 second")
		time.Sleep(5 * time.Second)
		s.GrabDom(u, "", ".layui-table", func(e *colly.HTMLElement) {
			e.DOM.Find("tr").Each(func(i int, qs *goquery.Selection) {
				if i == 0 {
					return
				}
				row := make([]string, 5)
				qs.Find("td").Each(func(j int, td *goquery.Selection) {
					row[j] = td.Text()
				})
				t := strings.TrimSpace(row[0]) + ":" + strings.TrimSpace(row[1])
				s.SpiderProxyList[t] = t
				return
			})
		})
	}
}
