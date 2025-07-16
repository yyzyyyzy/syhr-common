package parse

import (
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
)

func ParseTags(lang string) []language.Tag {
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		logx.Errorw("parse accept-language failed", logx.Field("detail", err))
		return []language.Tag{language.Chinese}
	}

	return tags
}
