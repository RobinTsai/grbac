package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type PermissionDoc struct {
	Id              int      `json:"id"`
	Host            string   `json:"host"`             // default *
	Path            string   `json:"path"`             // default *
	Method          string   `json:"method"`           // "{GET}"
	AuthorizedRoles []string `json:"authorized_roles"` //
	ForbiddenRoles  []string `json:"forbidden_roles"`
	AllowAnyone     bool     `json:"allow_anyone"`
}

type Permission struct {
	*PermissionDoc
	Pkg                   string   `json:"pkg"`
	Filepath              string   `json:"filepath"`
	RawRouterLine         string   `json:"rawRouterLine"`
	RawAuthRolesLine      string   `json:"rawAuthRolesLine"`
	RawForbiddenRolesLine string   `json:"rawForbiddenRolesLine"`
	Tags                  []string `json:"-"`
}

func (p *Permission) Parse() error {
	if p.RawRouterLine == "" {
		return fmt.Errorf("empty permission raw router line")
	}

	if err := p.parseRouterLine(); err != nil {
		return err
	}

	p.AuthorizedRoles = p.parseRolesLine(p.RawAuthRolesLine, "*")
	p.ForbiddenRoles = p.parseRolesLine(p.RawForbiddenRolesLine, "")

	return nil
}

// @Router       /admin/users [get]
var routerPattern = regexp.MustCompile(`^(/[\w./\-{}+:$]*)[[:blank:]]+\[(\w+)]`)

var routerRegex = regexp.MustCompile(`\{.*}`)

func (p *Permission) parseRouterLine() error {
	matches := routerPattern.FindStringSubmatch(p.RawRouterLine)
	if len(matches) != 3 {
		return fmt.Errorf("can not parse router comment \"%s\"", p.RawRouterLine)
	}

	p.Path = string(routerRegex.ReplaceAll([]byte(matches[1]), []byte("*")))
	p.Method = strings.ToUpper(matches[2])

	return nil
}

func (p *Permission) parseRolesLine(line string, def string) []string {
	res := make([]string, 0, 3)
	arr := strings.Split(line, ",")
	for _, s := range arr {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		res = append(res, s)
	}
	if def != "" && len(res) == 0 {
		res = []string{def}
	}
	return res
}

func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.PermissionDoc)
}
