package parser

import (
	"sort"
	"strings"
)

func FilterNil(in []*Permission) []*Permission {
	swap := func(i, j int) {
		in[i], in[j] = in[j], in[i]
	}

	end := len(in) - 1
	for i := 0; i <= end; {
		if in[i] == nil {
			swap(i, end)
			end--
		} else {
			i++
		}
	}
	return in[:end+1]
}

func GroupPathByRoles(ps []*Permission) {
	gs := make(map[string][]*Permission)
	for _, p := range ps {
		roleKey := func(p *Permission) string {
			sort.Sort(sort.StringSlice(p.AuthorizedRoles))
			ar := "+" + strings.Join(p.AuthorizedRoles, "+")

			sort.Sort(sort.StringSlice(p.ForbiddenRoles))
			fr := "-" + strings.Join(p.ForbiddenRoles, "-")

			aa := "=false"
			if p.AllowAnyone {
				aa = "=true"
			}

			return strings.Join([]string{ar, fr, aa}, "")
		}(p)
		gs[roleKey] = append(gs[roleKey], p)
	}
	_ = gs
}

// PS 排序，按长的 URL 在前进行排序
// 一样长的按字母排序，字母序大的在前
type PS []*Permission

func (p PS) Less(i, j int) bool {
	as := strings.Split(p[i].PermissionDoc.Path, "/")
	bs := strings.Split(p[j].PermissionDoc.Path, "/")
	ai, bi := 0, 0
	for ai < len(as) && bi < len(bs) {
		if as[ai] == bs[bi] || as[ai] == "*" || bs[bi] == "*" {
			ai++
			bi++
			continue
		} else if bs[bi] == "**" {
			ai++
			return true
		} else if as[ai] == "**" {
			bi++
			return false
		}
		return as[ai] > bs[bi]
	}

	return len(as) > len(bs)
}
func (p PS) Len() int {
	return len(p)
}
func (p PS) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func AggregatePath(groupedPs []*Permission) {
	sort.Sort(PS(groupedPs))
}
