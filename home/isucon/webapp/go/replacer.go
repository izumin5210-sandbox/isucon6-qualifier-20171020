package main

import (
	"fmt"
	"html"
	"net/url"
	"strings"
)

var (
	replacer = strings.NewReplacer()
)

func replace(s string) string {
	return replacer.Replace(s)
}

func updateReplacer() {
	keywords := getKeywordsOrderByLength()
	oldnew := make([]string, 0, 2*len(keywords))
	for _, k := range keywords {
		u, err := url.Parse(fmt.Sprintf("%s/keyword/%s", baseUrl.String(), pathURIEscape(k)))
		panicIf(err)
		a := fmt.Sprintf("<a href=\"%s\">%s</a>", u, html.EscapeString(k))
		oldnew = append(oldnew, k, a)
	}
	replacer = strings.NewReplacer(oldnew...)
}
