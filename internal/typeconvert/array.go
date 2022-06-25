package typeconvert

import "github.com/hashicorp/terraform-plugin-framework/types"

func ToTypesStrings(s []string) []types.String {
	if s == nil {
		return nil
	}
	t := make([]types.String, len(s))
	for i, se := range s {
		t[i] = types.String{
			Value: se,
		}
	}
	return t
}

func ToStringSlice(t []types.String) []string {
	if t == nil {
		return nil
	}
	s := make([]string, len(t))
	for i, te := range t {
		s[i] = te.Value
	}
	return s
}
