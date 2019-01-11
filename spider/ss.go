package spider

import (
	"github.com/gocolly/colly"
	"log"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"regexp"
	"strings"
	"strconv"
	"github.com/robertkrimen/otto"
)

type SsSpider struct {
	BaseSpider
}

func (s *SsSpider) Initialize() {
	s.Prefix = "ss_"
}

func (s *SsSpider) Spider() {
	//urls := []string{"https://www.xicidaili.com/nn/1", "https://www.xicidaili.com/nn/2", "https://www.xicidaili.com/nn/3", "https://www.xicidaili.com/nn/4"}
	c := colly.NewCollector(
		colly.Async(false),
	)
	cookie := ""
	c.OnError(func(r *colly.Response, err error) {
		//log.Println("response received", r.Headers.Get("Cookie"))
		log.Println(r.Headers.Get("Set-Cookie"))
		log.Println("haha")
		c1 := r.Headers.Get("Set-Cookie")
		body := string(r.Body)
		fmt.Println(body)
		reg, _ := regexp.Compile(`window.onload=setTimeout\("[a-zA-Z]+\((\d+)\)",\s*\d+\);`)
		matchs := reg.FindAllString(body, -1)
		fmt.Println(matchs)
		//[window.onload=setTimeout("ky(157)", 200);]
		if len(matchs) == 0 {
			log.Println("regrex fail")
			return
		}
		pl, _ := strconv.Atoi(strings.Split(strings.Split(matchs[0], "(")[2], ")")[0])
		funcName := strings.TrimLeft(strings.TrimSpace(strings.Split(matchs[0], "(")[1]), `"`)
		fmt.Println(funcName)
		fmt.Println(pl)
		reg2, _ := regexp.Compile(`; function.*</script> </body>`)
		matchs2 := reg2.FindAllString(body, -1)
		if len(matchs2) == 0 {
			log.Println("regrex fail func")
			return
		}
		funcBody := strings.TrimRight(strings.TrimLeft(matchs2[0], "; "), "</script> </body>")
		funcBody = strings.Replace(funcBody, `eval("qo=eval;qo(po);");`, "return po;", -1)
		fmt.Println(funcBody)
		vm := otto.New()
		vm.Run(funcBody)
		cc, err := vm.Call(funcName, nil, pl)
		if err != nil {
			log.Fatal("func exec fail")
		}
		fmt.Println(cc)
		cookie += strings.Split(c1, ";")[0] + "; " + strings.TrimLeft(strings.Split(cc.String(), ";")[0], "document.cookie='")
		fmt.Println(cookie)

	})

	header := map[string]string{}
	hdr := s.GetHeader(header)
	c.Request("GET", "http://www.66ip.cn", nil, nil, hdr)
	log.Println("haah")
	if cookie == "" {
		return
	}
	format := "http://www.66ip.cn/areaindex_%d/1.html"
	for page := 1; page <= 24; page++ {
		u := fmt.Sprintf(format, page)
		log.Println("sleep 10 second")
		time.Sleep(10 * time.Second)
		s.GrabDom(u, cookie, "#footer>div>table", func(e *colly.HTMLElement) {
			e.DOM.Find("tr").Each(func(i int, qs *goquery.Selection) {
				if i == 0 {
					return
				}
				row := make([]string, 5)
				qs.Find("td").Each(func(j int, td *goquery.Selection) {
					row[j] = td.Text()
				})
				t := strings.TrimSpace(row[0]) + ":" + strings.TrimSpace(row[1])
				//fmt.Println(t)
				s.SpiderProxyList[t] = t
				return
			})
		})
	}
}
