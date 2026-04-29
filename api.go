package main

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

const base = "https://news.hada.io"

var (
	reBR     = regexp.MustCompile(`(?i)<br\s*/?>`)
	reTag    = regexp.MustCompile(`<[^>]+>`)
	reSpaces = regexp.MustCompile(`[ \t]+`)
	reNL     = regexp.MustCompile(`\n{3,}`)

	reTopicID  = regexp.MustCompile(`data-topic-state-id='(\d+)'`)
	reTitle    = regexp.MustCompile(`(?s)<h1>(.*?)</h1>`)
	reNofollow = regexp.MustCompile(`href='([^']+)' rel='nofollow'`)
	reDomain   = regexp.MustCompile(`<span class=topicurl>\(([^)]+)\)</span>`)
	reDesc     = regexp.MustCompile(`(?s)class='c99 breakall'>(.*?)</a>`)
	rePoints   = regexp.MustCompile(`<span id='tp\d+'>(\d+)</span>`)
	reAuthor   = regexp.MustCompile(`href='/@([^']+)'>[^<]+</a>`)
	reTime     = regexp.MustCompile(`href='/@[^']*'>[^<]+</a>\s*([^<]+?)\s*<span id='unvote`)
	reCmtCount = regexp.MustCompile(`go=comments[^>]*>([^<]+)</a>`)

	reBoldTitle = regexp.MustCompile(`(?s)class='bold ud'><h1>(.*?)</h1>`)
	reBoldURL   = regexp.MustCompile(`href='([^']+)' class='bold ud'`)
	reSpanTitle = regexp.MustCompile(`<span title='([^']+)'>`)

	reCID     = regexp.MustCompile(`id=cid(\d+)`)
	reDepth   = regexp.MustCompile(`style=--depth:(\d+)`)
	reCTime   = regexp.MustCompile(`href='comment\?id=\d+'>([^<]+)</a>`)
	reContent = regexp.MustCompile(`(?s)class='comment_contents'>(.*?)</span>`)
)

type Topic struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Domain   string `json:"domain"`
	Desc     string `json:"desc"`
	Points   int    `json:"points"`
	Author   string `json:"author"`
	Time     string `json:"time"`
	Comments string `json:"comments"`
	TopicURL string `json:"topic_url"`
}

type Comment struct {
	ID     string `json:"id"`
	Depth  int    `json:"depth"`
	Author string `json:"author"`
	Time   string `json:"time"`
	Text   string `json:"text"`
}

type TopicDetail struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Points   int       `json:"points"`
	Author   string    `json:"author"`
	Time     string    `json:"time"`
	Body     string    `json:"body"`
	Comments []Comment `json:"comments"`
	TopicURL string    `json:"topic_url"`
}

func cleanHTML(text string) string {
	text = reBR.ReplaceAllString(text, "\n")
	text = reTag.ReplaceAllString(text, "")
	text = reSpaces.ReplaceAllString(text, " ")
	lines := strings.Split(text, "\n")
	for i, ln := range lines {
		lines[i] = strings.TrimSpace(ln)
	}
	text = strings.Join(lines, "\n")
	text = reNL.ReplaceAllString(text, "\n\n")
	return html.UnescapeString(strings.TrimSpace(text))
}

func sub1(m []string) string {
	if m == nil {
		return ""
	}
	return m[1]
}

func resolveURL(u, tid string) string {
	if strings.HasPrefix(u, "topic?") || strings.HasPrefix(u, "/topic?") {
		return base + "/" + strings.TrimPrefix(u, "/")
	}
	return u
}

func parseItem(chunk string) *Topic {
	idM := reTopicID.FindStringSubmatch(chunk)
	titleM := reTitle.FindStringSubmatch(chunk)
	if idM == nil || titleM == nil {
		return nil
	}

	tid := idM[1]
	title := html.UnescapeString(strings.TrimSpace(titleM[1]))

	u := base + "/topic?id=" + tid
	if m := reNofollow.FindStringSubmatch(chunk); m != nil {
		u = resolveURL(m[1], tid)
	}

	domain := sub1(reDomain.FindStringSubmatch(chunk))

	var desc string
	if m := reDesc.FindStringSubmatch(chunk); m != nil {
		desc = cleanHTML(m[1])
	}

	points := 0
	if m := rePoints.FindStringSubmatch(chunk); m != nil {
		points, _ = strconv.Atoi(m[1])
	}

	author := sub1(reAuthor.FindStringSubmatch(chunk))

	t := ""
	if m := reTime.FindStringSubmatch(chunk); m != nil {
		t = strings.TrimSpace(m[1])
	}

	cmt := ""
	if m := reCmtCount.FindStringSubmatch(chunk); m != nil {
		cmt = strings.TrimSpace(m[1])
	}

	return &Topic{
		ID:       tid,
		Title:    title,
		URL:      u,
		Domain:   domain,
		Desc:     desc,
		Points:   points,
		Author:   author,
		Time:     t,
		Comments: cmt,
		TopicURL: base + "/topic?id=" + tid,
	}
}

func topics(section string, page int) ([]Topic, error) {
	u := base + section
	if page > 1 {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u += sep + "page=" + strconv.Itoa(page)
	}
	htmlStr, err := getHTML(u)
	if err != nil {
		return nil, err
	}
	result := []Topic{}
	for _, chunk := range strings.Split(htmlStr, "<div class='topic_row'")[1:] {
		if item := parseItem(chunk); item != nil {
			result = append(result, *item)
		}
	}
	return result, nil
}

func topic(tid string) (*TopicDetail, error) {
	htmlStr, err := getHTML(base + "/topic?id=" + tid)
	if err != nil {
		return nil, err
	}

	title := ""
	if m := reBoldTitle.FindStringSubmatch(htmlStr); m != nil {
		title = html.UnescapeString(strings.TrimSpace(m[1]))
	}

	u := base + "/topic?id=" + tid
	if m := reBoldURL.FindStringSubmatch(htmlStr); m != nil {
		u = resolveURL(m[1], tid)
	}

	points := 0
	if m := rePoints.FindStringSubmatch(htmlStr); m != nil {
		points, _ = strconv.Atoi(m[1])
	}

	author := sub1(reAuthor.FindStringSubmatch(htmlStr))
	t := sub1(reSpanTitle.FindStringSubmatch(htmlStr))

	var body string
	const marker = "id='topic_contents'>"
	if bodyStart := strings.Index(htmlStr, marker); bodyStart >= 0 {
		raw := htmlStr[bodyStart+len(marker):]
		if end := strings.Index(raw, "</div></div></div>"); end > 0 {
			raw = raw[:end]
		} else if len(raw) > 10000 {
			raw = raw[:10000]
		}
		body = cleanHTML(raw)
	}

	comments := []Comment{}
	for _, chunk := range strings.Split(htmlStr, "<div class=comment_row")[1:] {
		cidM := reCID.FindStringSubmatch(chunk)
		if cidM == nil {
			continue
		}
		depth := 0
		if m := reDepth.FindStringSubmatch(chunk); m != nil {
			depth, _ = strconv.Atoi(m[1])
		}
		cauthor := sub1(reAuthor.FindStringSubmatch(chunk))
		ctime := ""
		if m := reCTime.FindStringSubmatch(chunk); m != nil {
			ctime = strings.TrimSpace(m[1])
		}
		text := ""
		if m := reContent.FindStringSubmatch(chunk); m != nil {
			text = cleanHTML(m[1])
		}
		comments = append(comments, Comment{
			ID:     cidM[1],
			Depth:  depth,
			Author: cauthor,
			Time:   ctime,
			Text:   text,
		})
	}

	return &TopicDetail{
		ID:       tid,
		Title:    title,
		URL:      u,
		Points:   points,
		Author:   author,
		Time:     t,
		Body:     body,
		Comments: comments,
		TopicURL: base + "/topic?id=" + tid,
	}, nil
}
