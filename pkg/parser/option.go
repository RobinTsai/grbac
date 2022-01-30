package parser

func SetExcludeFiles(files []string) func(parser *Parser) {
	return func(parser *Parser) {
		parser.excludeFiles = files
	}
}
