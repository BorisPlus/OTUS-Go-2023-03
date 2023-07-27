package regexhandlers

import (
	"net/url"
	"regexp"
	"strings"
)

type QueryPathPattern struct {
	pattern      string
	ParamsNaming []string
	compiled     *regexp.Regexp
}

func NewQueryPathPattern(pattern string, ParamsNaming []string) *QueryPathPattern {
	qpp := new(QueryPathPattern)
	qpp.pattern = pattern
	qpp.mustCompile()
	qpp.ParamsNaming = ParamsNaming
	return qpp
}

func (qpp *QueryPathPattern) normalize() string {
	if !strings.HasSuffix(qpp.pattern, "$") {
		qpp.pattern += `$`
	}
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{numeric}`, `(\d*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{string}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{filename}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{any}`, `([^/]*)`)
	return qpp.pattern
}

func (qpp *QueryPathPattern) mustCompile() {
	qpp.compiled = regexp.MustCompile(qpp.normalize())
}

func (qpp *QueryPathPattern) match(url string) bool {
	return qpp.compiled.MatchString(url)
}

func (qpp *QueryPathPattern) fetch(url string) url.Values {
	parsedParamsInUrl := qpp.compiled.FindAllStringSubmatch(url, -1)
	paramsValues := make(map[string][]string)
	for _, submatch := range parsedParamsInUrl {
		for orderIndex, value := range submatch {
			if orderIndex == 0 {
				continue
			}
			paramsValues[qpp.ParamsNaming[orderIndex-1]] = append(paramsValues[qpp.ParamsNaming[orderIndex-1]], value)
		}
	}
	return paramsValues
}
