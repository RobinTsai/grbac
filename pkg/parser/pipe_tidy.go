package parser

import (
	"fmt"
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
func addRoleKey(p *Permission) {
	sort.Sort(sort.StringSlice(p.AuthorizedRoles))
	ar := "+" + strings.Join(p.AuthorizedRoles, "+")

	sort.Sort(sort.StringSlice(p.ForbiddenRoles))
	fr := "-" + strings.Join(p.ForbiddenRoles, "-")

	aa := "=false"
	if p.AllowAnyone {
		aa = "=true"
	}

	key := strings.Join([]string{ar, fr, aa}, "")
	p.PermKey = key
}
func addFrags(p *Permission) {
	p.Frags = strings.Split(p.Path, "/")
}
func addSameFragCount(p *Permission, last []string, lastKey string) {
	if p.PermKey != lastKey {
		p.SameFragCountWithLast = -1
		return
	}
	max := len(last)
	if t := len(p.Frags); t < max {
		max = t
	}

	i := 0
	for ; i < max; i++ {
		if last[i] == p.Frags[i] {
			continue
		}
		if last[i] == "*" || last[i] == "**" ||
			p.Frags[i] == "*" || p.Frags[i] == "**" {
			continue
		}
		break
	}
	p.SameFragCountWithLast = i
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

type PSWithCount struct {
	PS
}

func (p PSWithCount) Less(i, j int) bool {
	return p.PS[i].SameFragCountWithLast < p.PS[j].SameFragCountWithLast
}
func (p PSWithCount) Len() int {
	return len(p.PS)
}
func (p PSWithCount) Swap(i, j int) {
	p.PS.Swap(i, j)
}

func AggregatePath(ps []*Permission) []*Permission {
	sort.Sort(PS(ps))

	for _, p := range ps {
		addRoleKey(p)
	}

	for _, p := range ps {
		addFrags(p)
	}

	for i, p := range ps {
		if i == 0 {
			addSameFragCount(p, nil, "")
			continue
		}
		last := ps[i-1]
		addSameFragCount(p, last.Frags, last.PermKey)
	}

	// 双指针
	firstP := 0
	secondP := 0
	for firstP < len(ps) {
		secondP++
		if secondP >= len(ps) || ps[secondP].SameFragCountWithLast == -1 {
			curPsGroup := ps[firstP:secondP]
			// todo: check this group

			other := append([]*Permission{}, ps[:firstP]...)
			other = append(other, ps[secondP:]...)
			optimizeByGroup(curPsGroup, other)
			_ = curPsGroup
			firstP = secondP
		}
	}

	fmt.Println("before aggregate ps len", len(ps))
	ps = aggregateSamePath(ps)
	reStampID(ps)

	return ps
}

func isFragEqual(a, b string) bool {
	if a == "*" && b == "*" ||
		a == "**" && b == "**" {
		return true
	}
	return false
}

func isFragsEqual(as, bs []string) bool {
	ai, bi := 0, 0
	for ai < len(as) && bi < len(bs) {
		if as[ai] == bs[bi] || as[ai] == "*" || bs[bi] == "*" {
			ai++
			bi++
			continue
		} else if bs[bi] == "**" || as[ai] == "**" {
			return true
		}
		return false
	}

	return true
}

// 优化 this, other 做对比用
func optimizeByGroup(this, other []*Permission) {
	sort.Sort(PSWithCount{PS(this)})
	if len(this) == 0 {
		return
	}
	okPrefixs := []string{}
First:
	for _, p0 := range this {
		for _, oked := range okPrefixs {
			if strings.HasPrefix(p0.Path, oked) {
				continue First
			}
		}
		// 是否可以
		okPrefix := ""
		sameCount := 0
	Second:
		for i := 0; i < len(p0.Frags); i++ {
			if p0.Frags[i] == "" {
				sameCount = 0
				continue Second
			}

			curFrags := p0.Frags[:i+1]
			for _, p := range this {
				if isFragsEqual(curFrags, p.Frags) {
					sameCount++
				}
			}

			if sameCount <= 2 {
				sameCount = 0
				continue Second
			}
			// 合适优化

			for _, otherP := range other {
				if isFragsEqual(curFrags, otherP.Frags) { // 冲突
					sameCount = 0
					continue Second
				}
			}

			okPrefix = strings.Join(curFrags, "/") + "/"
			break
		}

		// 可以，设置星号
		if okPrefix != "" {
			okPrefixs = append(okPrefixs, okPrefix)
			fmt.Println(" -------- ok", okPrefix)
			for _, p := range this {
				if strings.HasPrefix(p.Path, okPrefix) {
					p.Path = okPrefix + "**"
					p.Frags = strings.Split(p.Path, "/")
					p.SameFragCountWithLast = -1
				}
			}
		}
	}
}

func aggregateSamePath(ps []*Permission) []*Permission {
	group := make(map[string][]*Permission, 0)

	getKey := func(p *Permission) string {
		return p.Path + "+" + strings.Join(p.AuthorizedRoles, "+") +
			"-" + strings.Join(p.ForbiddenRoles, "-")
	}

	for _, p := range ps {
		group[getKey(p)] = append(group[getKey(p)], p)
	}

	for k, g := range group {
		fmt.Println("--- group info", k, len(g))
	}
	result := make([]*Permission, 0)
	for _, gps := range group {
		p0 := gps[0]
		p0.Methods = p0.GetMethodsFromMethodStr()
		for _, p := range gps[1:] {
			p0.Methods = append(p0.Methods, p.GetMethodsFromMethodStr()...)
		}

		p0.SetMethodFromMethods()
		result = append(result, p0)
	}
	return result
}

func reStampID(ps []*Permission) {
	for i, p := range ps {
		p.Id = i
	}
}
