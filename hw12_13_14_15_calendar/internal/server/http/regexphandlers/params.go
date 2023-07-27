package regexhandlers

import (
	"net/url"
	"regexp"
	"strings"
)

type QueryPathPattern struct {
	Pattern      string
	ParamsNaming []string
}

func NewQueryPathPattern(Pattern string, ParamsNaming []string) QueryPathPattern {
	return QueryPathPattern{
		Pattern, ParamsNaming,
	}
}

func (qpp *QueryPathPattern) normalize() string {
	if !strings.HasSuffix(qpp.Pattern, "$") {
		qpp.Pattern += `$`
	}
	qpp.Pattern = strings.ReplaceAll(qpp.Pattern, `{numeric}`, `(\d*)`)
	qpp.Pattern = strings.ReplaceAll(qpp.Pattern, `{string}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.Pattern = strings.ReplaceAll(qpp.Pattern, `{filename}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.Pattern = strings.ReplaceAll(qpp.Pattern, `{any}`, `([^/]*)`)
	return qpp.Pattern
}

func (qpp *QueryPathPattern) MustCompile() *regexp.Regexp {
	return regexp.MustCompile(qpp.normalize())
}

func (qpp *QueryPathPattern) Match(url string) bool {
	return qpp.MustCompile().MatchString(url)
}

func (qpp *QueryPathPattern) Fetch(url string) url.Values {
	parsedParamsInUrl := regexp.MustCompile(qpp.normalize()).FindAllStringSubmatch(url, -1)
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
