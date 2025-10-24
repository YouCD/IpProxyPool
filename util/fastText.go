package util

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var textBufPool = sync.Pool{
	New: func() interface{} { return new(strings.Builder) },
}

func FastText(sel *goquery.Selection) string {
	buf := textBufPool.Get().(*strings.Builder)
	buf.Reset()
	defer textBufPool.Put(buf)

	// 手动遍历 TextNode，避免 goquery.Text() 内部再 new Builder
	sel.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			buf.WriteString(strings.TrimSpace(s.Text()))
		}
	})
	return buf.String()
}
