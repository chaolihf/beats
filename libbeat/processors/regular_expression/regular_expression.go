package regular_expression

import (
	"fmt"
	"regexp"
)

// Map  represents the keys and their values extracted with the defined tokenizer.
type Map = map[string]string
type MapConverted = map[string]interface{}

// Regular is a tokenizer based on the Dissect syntax as defined at:
// https://www.elastic.co/guide/en/logstash/current/plugins-filters-dissect.html
type Regular struct {
	raw    string
	regexp *regexp.Regexp
}

func (d *Regular) Regular(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, errEmpty
	}
	regexp := d.regexp

	regex := regexp.MatchString(s)

	if !regex {
		fmt.Println("正则匹配不成功", regex, regexp)
		return nil, errParsingFailure
	}

	fmt.Println("正则匹配成功", regex)
	matches := regexp.FindStringSubmatch(s)
	if matches != nil {
		fmt.Println("匹配结果:", matches)
		result := make([]string, len(matches))
		// 遍历 matches 切片
		for i, match := range matches {
			if i > 0 {
				result[i] = match
				fmt.Printf("匹配 %d: %s\n", i, match)
			}
		}
		fmt.Printf("%q\n", result)

		return result, nil
	}
	return nil, nil
}

// Raw returns the raw tokenizer used to generate the actual parser.
func (d *Regular) Raw() string {
	return d.raw
}

// New creates a new Regular from a tokenized string.
func New(regular string) (*Regular, error) {
	//编译正则表达式
	r, err := newRegular(regular)
	if err != nil {
		return nil, err
	}
	return &Regular{regexp: r, raw: regular}, nil
}
