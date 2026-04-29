package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func toJSON(v any) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
	return strings.TrimRight(buf.String(), "\n")
}

func tg(text string) string {
	return "```\n" + text + "\n```"
}

func fmtTopics(data []Topic, format string) string {
	if format == "json" {
		return toJSON(data)
	}
	if len(data) == 0 {
		return tg("토픽이 없습니다.")
	}
	var lines []string
	for _, item := range data {
		meta := fmt.Sprintf("%dP  %s  %s", item.Points, item.Author, item.Time)
		if item.Comments != "" {
			meta += "  " + item.Comments
		}
		lines = append(lines, fmt.Sprintf("[%s] %s", item.ID, item.Title))
		lines = append(lines, "  "+item.URL)
		lines = append(lines, "  "+meta)
		if item.Desc != "" {
			runes := []rune(item.Desc)
			if len(runes) > 120 {
				lines = append(lines, "  "+string(runes[:120])+"...")
			} else {
				lines = append(lines, "  "+item.Desc)
			}
		}
		lines = append(lines, "")
	}
	return tg(strings.TrimRight(strings.Join(lines, "\n"), "\n"))
}

func fmtTopic(data *TopicDetail, format string) string {
	if format == "json" {
		return toJSON(data)
	}
	lines := []string{
		data.Title,
		data.URL,
		fmt.Sprintf("%dP  %s  %s", data.Points, data.Author, data.Time),
		"토픽: " + data.TopicURL,
		strings.Repeat("─", 50),
	}
	if data.Body != "" {
		lines = append(lines, "", data.Body)
	}
	if len(data.Comments) > 0 {
		lines = append(lines, "", fmt.Sprintf("── 댓글 %d개 %s", len(data.Comments), strings.Repeat("─", 30)))
		for _, c := range data.Comments {
			indent := strings.Repeat("  ", c.Depth)
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("%s[%s] %s", indent, c.Author, c.Time))
			for _, ln := range strings.Split(c.Text, "\n") {
				lines = append(lines, indent+"  "+ln)
			}
		}
	}
	return tg(strings.Join(lines, "\n"))
}
