package spider

import (
	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

type XiciSpider struct {
	BaseSpider
}

func (s *XiciSpider) Initialize() {
	s.Prefix = "xici_"
}

func (s *XiciSpider) Spider() {
	//urls := []string{"https://www.xicidaili.com/nn/1", "https://www.xicidaili.com/nn/2", "https://www.xicidaili.com/nn/3", "https://www.xicidaili.com/nn/4"}
	format := "https://www.xicidaili.com/nn/%d"
	for page := 1; page <= 10; page++ {
		u := fmt.Sprintf(format, page)
		s.GrabDom(u, "", "#ip_list", func(e *colly.HTMLElement) {
			e.DOM.Find("tr").Each(func(i int, qs *goquery.Selection) {
				if i == 0 {
					return
				}
				row := make([]string, 10)
				qs.Find("td").Each(func(j int, td *goquery.Selection) {
					row[j] = td.Text()
				})
				if row[5] != "HTTP" {
					return
				}
				t := row[1] + ":" + row[2]
				s.SpiderProxyList[t] = t
				return
			})
		})
	}
}
