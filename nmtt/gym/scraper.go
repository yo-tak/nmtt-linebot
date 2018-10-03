package gym

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type columnType int

func GetNakanoCurrentNotification() (string, error) {
	resp, err := http.Get("http://www.nakano-taiikukan.com/nakano/news/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("failed to read response:", err)
	}

	var notificationURL string
	// find a tag for news of gym schedule
	doc.Find("#event h4").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if idx := strings.LastIndex(s.Text(), "一般開放　【卓球】　週間予定"); idx > 0 {
			log.Println(s.Html())
			href, ok := s.Find("a").Attr("href")
			if !ok {
				return true
			}
			notificationURL = href
			log.Println("extracted url:", notificationURL)
			return false
		}
		return true
	})
	return notificationURL, nil
}
