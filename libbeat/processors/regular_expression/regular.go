package regular_expression

import (
	"fmt"
	"regexp"
)

// 编译正则表达式
func newRegular(regular string) (*regexp.Regexp, error) {
	regExp, err := regexp.Compile(regular)
	if err != nil {
		fmt.Println("正则表达式不合规")
		return nil, err
	}
	fmt.Println("正则表达式合规", regExp)
	return regExp, nil
}
