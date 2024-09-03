package github

import (
	"fmt"
)

func setProxyWeb(urlsStr string) string {

	return fmt.Sprintf("https://gh.xmly.dev/%s", urlsStr)
}
