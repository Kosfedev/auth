package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PlaceholderDollar = "$"
)

func Pretty(query string, placeholderBase string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		placeholder := fmt.Sprintf("%s%s", placeholderBase, strconv.Itoa(i+1))
		query = strings.Replace(query, placeholder, value, -1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
