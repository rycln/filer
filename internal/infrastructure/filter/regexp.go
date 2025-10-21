package filter

import (
	"regexp"
	"slices"
)

type RegexpFilter struct {
	pattern string
}

func NewRegexpFilter(pattern string) *RegexpFilter {
	return &RegexpFilter{pattern: pattern}
}

func (f *RegexpFilter) Filter(filenames []string) ([]string, error) {
	if f.pattern == "" {
		return filenames, nil
	}

	var filtered []string

	for _, filename := range filenames {
		matched, err := regexp.MatchString(f.pattern, filename)
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, filename)
		}
	}

	slices.Sort(filtered)

	return filtered, nil
}
