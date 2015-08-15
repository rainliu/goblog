package models

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"

	"appengine/user"
)

type Article struct {
	ArticleID   int
	Title       string
	Date        time.Time
	DisplayName string
	Author      user.User
	Category    string
	Tags        []string
	Content     template.HTML `datastore:"Content,noindex"`
	Views       int
	Comments    int
}

type Content struct {
	IsAdmin     bool
	Articles    []Article
	PrevEntries string
	NextEntries string
	IsSingle    bool
	Comments    []Comment
}

func Date2String(date time.Time) string {
	str := date.Local().String()

	return str[0:10]
}

func Author2String(author user.User) string {
	return author.String()
}

func Tags2String(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	var strBuffer bytes.Buffer
	for i := 0; i < len(tags); i++ {
		strBuffer.WriteString(strings.TrimSpace(tags[i]))
		if i != len(tags)-1 {
			strBuffer.WriteString(", ")
		}
	}

	return strBuffer.String()
}

func FullTextSelect(single bool) func(template.HTML) template.HTML {
	if single {
		return FullText2Full
	} else {
		return FullText2Excerpt
	}
}

func FullText2Full(text template.HTML) template.HTML {
	return text
}

func FullText2Excerpt(text template.HTML) template.HTML {
	if len(text) < MyBlogConfig.ExcerptsCharNumber {
		return text
	} else {
		return text[0:MyBlogConfig.ExcerptsCharNumber] + "[...]"
	}
}

func Comments2String(comments []string) string {
	return strconv.Itoa(len(comments))
}
