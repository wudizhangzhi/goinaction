package matchers

import (
	"chapter2/search"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		Guid        string   `xml:"guid"`
		PubDate     string   `xml:"pubDate"`
		GeoRssPoint string   `xml:"georss:point"`
	}

	image struct {
		XMLName xml.Name `xml:"image"`
		Url     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}

	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

type rssMatcher struct{}

func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}

func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]search.Result, error) {
	var results []search.Result
	log.Printf("Search FeedType[%s] Site[%s] Uri[%s]\n", feed.Type, feed.Name, feed.URI)
	document, err := m.retrieve(feed)
	if err != nil {
		return nil, err
	}

	for _, channelItem := range document.Channel.Item {
		// 检查 title
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results,
				search.Result{
					Field:   "Title",
					Content: channelItem.Title,
				})
		}
		// 检查 description
		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
		if err != nil {
			return nil, err
		}
		if matched {
			results = append(results,
				search.Result{
					Field:   "Description",
					Content: channelItem.Description,
				})
		}
	}
	return results, nil
}

func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	resp, err := http.Get(feed.URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Http Response Error %d\n", resp.StatusCode)
	}

	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err
}
