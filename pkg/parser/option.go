package parser

import "strings"

func SetExcludeFiles(files []string) func(parser *Parser) {
	return func(parser *Parser) {
		parser.excludeFiles = files
	}
}

func SetTag(tag string) func(parser *Parser) {
	return func(parser *Parser) {
		parser.tag = strings.ToLower(tag)
	}
}

func SetSsrole(role string) func(parser *Parser) {
	return func(parser *Parser) {
		parser.ssRole = role
	}
}
