package main

import (
	"fmt"
	"regexp"
)

type Codec interface {
	Encode(input string) string
	Decode(input string) []string
}

/*
* RESP message parser.
/*   - decode()
/*  - encode()
*
*/
type RESPParser struct{}

func (parser *RESPParser) Encode(input string) string {
	/* Server 返回 RESP bulk string 格式消息 */
	return fmt.Sprintf("$%d\r\n%s\r\n", len(input), input)
}

func (parser *RESPParser) Decode(input string) []string {
	/* 解析 client command 消息提取命令+参数 */

	// 正则提取参数
	regElement := regexp.MustCompile(`\$\d*\r\n([^\r\n]+)\r\n`)
	elementMatch := regElement.FindAllStringSubmatch(input, -1)

	elements := make([]string, 0, len(elementMatch))
	for _, match := range elementMatch {
		elements = append(elements, match[1])
	}

	return elements
}
