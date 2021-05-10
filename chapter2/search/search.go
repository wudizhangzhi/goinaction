package search

import (
	"log"
	"sync"
)

// 注册过的Matcher
var matchers = make(map[string]Matcher)


func Run(searchTerm string) {
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	var results = make(chan *Result)

	var waitGroup sync.WaitGroup

	for _, feed := range feeds {
		waitGroup.Add(1)

		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, &feed)
	}

	go func() {
		waitGroup.Wait()

		close(results)
	}()

	Display(results)
}

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "已经存在")
	}

	matchers[feedType] = matcher
}