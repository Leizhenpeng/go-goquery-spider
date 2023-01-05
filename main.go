package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type BookInfo struct {
	BookName  string
	Author    string
	Publisher string
	Code      string
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://202.114.238.250/XSTB/xstb_right?code=all&Doctypecode=0&starttime=2020-12-5&endtime=2023-1-5&orderbyzd=DC_FABMC.DTFOUND&orderbypx=desc&SetTakeNum=1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cookie", "ASP.NET_SessionId=vfrnjtv2y4ir0a4kbncz1ufy")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://202.114.238.250/XSTB/xstb_left")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := parseContent(bodyText)

	writeContent(r)

}

func writeContent(r []BookInfo) {
	file := "result.csv"
	if _, err := os.Stat(file); err == nil {
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		for _, v := range r {
			_, err := f.WriteString(fmt.Sprintf("%s,%s,%s,%s \r ", v.BookName, v.Author, v.Publisher, v.Code))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		f, err := os.Create(file)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		f.WriteString("BookName, Author,Publisher,Code \r")
		for _, v := range r {
			f.WriteString(v.BookName + "," + v.Author + "," + v.Publisher + "," + v.Code + "\r")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func parseContent(text []byte) []BookInfo {

	// goQuery
	var result []BookInfo
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(text))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			var book BookInfo
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				str := s.Text()
				str = replace(str)
				switch i {
				case 1:
					book.BookName = str
				case 2:
					book.Author = str
				case 3:
					book.Publisher = str
				case 4:
					book.Code = str
				}
			})
			fmt.Printf("%+v", book)
			fmt.Println()
			if book.BookName != "" {
				result = append(result, book)
			}
		})
	})
	return result
}

func replace(str string) string {
	//replace ,  with _ in the string
	str = string(bytes.ReplaceAll([]byte(str), []byte(","), []byte("_")))
	return string(str)
}
