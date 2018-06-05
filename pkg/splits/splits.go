package splits

import "regexp"

// SplitParameters splits shell command parameters, taking quoting in account.
func SplitParameters(s string) []string {
	r := regexp.MustCompile(`"[^']*"|'[^']*'|[^ ]+`)
	params := r.FindAllString(s, -1)
	return unquote(params)
}

// SplitAttrs returns
func SplitAttrs(s string) []string {
	r := regexp.MustCompile(`"[^']*"|'[^']*'|\[[^]]*\]|[^.]+`)
	params := r.FindAllString(s, -1)
	return unquote(params)
}

func unquote(params []string) []string {
	for i, p := range params {
		if p[0] == '"' || p[0] == '\'' {
			params[i] = p[1 : len(p)-1]
		}
	}
	return params
}
