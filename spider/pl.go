package spider

import (
	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"log"
	"strings"
	"net/url"
)

type PlSpider struct {
	BaseSpider
}

func (s *PlSpider) Initialize() {
	s.Prefix = "pl_"
}

func (s *PlSpider) Spider() {
	//urls := []string{"https://www.xicidaili.com/nn/1", "https://www.xicidaili.com/nn/2", "https://www.xicidaili.com/nn/3", "https://www.xicidaili.com/nn/4"}
	domain := "http://www.proxylists.net/"
	u := "cn_0.html"
	c := colly.NewCollector(
		colly.Async(false),
	)
	urls := []string{}
	c.OnHTML("table>tbody>tr>td>table>tbody", func(qs *colly.HTMLElement) {
		qs.DOM.Find("tr:last-child").Find("td>b>a").Each(func(i int, qs *goquery.Selection) {

			href, _ := qs.Attr("href")
			fmt.Println(qs.Attr("href"))
			urls = append(urls, domain+href)
		})
	})
	c.Visit(domain + u)
	log.Println(urls)
	if len(urls) == 0 {
		log.Println("get no url")
		return
	}
	time.Sleep(5 * time.Second)
	for _, l := range urls {
		log.Println("sleep 5 second")
		time.Sleep(5 * time.Second)
		s.GrabDom(l, "", "table>tbody>tr>td>table>tbody", func(e *colly.HTMLElement) {
			e.DOM.Find("tr:not(:last-child)").Each(func(i int, qs *goquery.Selection) {
				if i < 2 {
					return
				}
				row := make([]string, 2)
				qs.Find("td").Each(func(j int, td *goquery.Selection) {
					row[j] = td.Text()
				})
				//eval(unescape('%73%65%6c%66%2e%64%6f%63%75%6d%65%6e%74%2e%77%72%69%74%65%6c%6e%28%22%31%32%34%2e%31%37%32%2e%32%33%32%2e%34%39%22%29%3b'));Please enable javascript
				//self.document.writeln("221.226.11.228");
				tt := strings.Split(strings.TrimSpace(row[0]), "'")
				ho,_ := url.QueryUnescape(tt[1])
				host:= strings.Split(ho, `"`)
				t := host[1] + ":" + strings.TrimSpace(row[1])
				s.SpiderProxyList[t] = t
				fmt.Println(t)
				return
			})
		})
	}
}
