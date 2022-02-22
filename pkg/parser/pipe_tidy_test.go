package parser

import "testing"

func TestFilterNil(t *testing.T) {
	match := func(in, out []*Permission) int {
		if len(in) != len(out) {
			return 1
		}

		for i := 0; i < len(in); i++ {
			if in[i] == nil || out[i] == nil {
				return 2
			}
			if in[i].Pkg != out[i].Pkg {
				return 3
			}
		}
		return 0
	}

	cases := []struct {
		Input []*Permission
		Want  []*Permission
	}{
		{
			[]*Permission{{Pkg: "1"}, {Pkg: "2"}, nil},
			[]*Permission{{Pkg: "1"}, {Pkg: "2"}},
		},
		{
			[]*Permission{{Pkg: "1"}, nil, {Pkg: "2"}},
			[]*Permission{{Pkg: "1"}, {Pkg: "2"}},
		},
		{
			[]*Permission{{Pkg: "1"}, nil, {Pkg: "2"}, nil},
			[]*Permission{{Pkg: "1"}, {Pkg: "2"}},
		},
		{
			[]*Permission{nil},
			[]*Permission{},
		},
	}

	for _, c := range cases {
		output := FilterNil(c.Input)
		if res := match(c.Want, output); res != 0 {
			t.Error("-------- err code", res, c.Want, output)
		}
	}
}
