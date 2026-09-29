package mask

func Apply(s string) string {
	switch {
	case s == "":
		return ""
	case len(s) < 4:
		return string(s[0]) + "..."
	default:
		return string(s[0:4]) + "..."
	}
}
