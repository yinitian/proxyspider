#Golang Proxy Spider
目前只支持http协议代理，因为记录到文件，虽然发现也可以支持https和socks5，但是懒得修改了，有兴趣的可以添加<br/>
使用者请不要频繁抓取，一天一次即可，否则可能被封IP或者添加防止爬虫规则

##采集代理
* 66IP
* 西刺代理
* 快代理
* 89ip
* ip3366

##使用方式
```
./main -t spider -p ss
./main -t spider -p xici
./main -t spider -p kuaifa
./main -t spider -p en
./main -t spider -p cloud
```

##后记
如果需要拓展其他网站抓取代理，请参考spider目录下xici.go
同时在main.go里面添加新抓取类型
