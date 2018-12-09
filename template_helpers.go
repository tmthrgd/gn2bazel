package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func printStringSlice(v []string) string {
	var b strings.Builder
	for i, s := range v {
		if i > 0 {
			b.WriteString(", ")
		}

		fmt.Fprintf(&b, "%q", s)
	}

	return b.String()
}

func isHeader(name string) bool {
	switch filepath.Ext(name) {
	case ".h", ".hpp":
		return true
	default:
		return false
	}
}

func filterHeadersFromSources(v []string) []string {
	var hdrs []string
	for _, s := range v {
		if isHeader(s) {
			hdrs = append(hdrs, s)
		}
	}

	return hdrs
}

func toStringSlice(v []interface{}) []string {
	var sv []string
	for _, s := range v {
		sv = append(sv, s.(string))
	}

	return sv
}

var (
	publicVisibility  = []string{"//visibility:public"}
	privateVisibility = []string{"//visibility:private"}
)

func toBazelVisibility(v []string) []string {
	if len(v) == 1 {
		switch v[0] {
		case "*":
			return publicVisibility
		case "//.:*":
			return privateVisibility
		}
	}

	var vis []string
	for _, s := range v {
		if strings.HasPrefix(s, "//.:") {
			s = "//:" + strings.TrimPrefix(s, "//.:")
		}

		switch {
		case s == "*":
			return publicVisibility
		case strings.HasSuffix(s, ":*"):
			pkg := strings.TrimSuffix(s, ":*") + ":__pkg__"
			vis = append(vis, pkg)
		default:
			vis = append(vis, s)
		}
	}

	return vis
}

func printBool(v bool) string {
	if v {
		return "True"
	}

	return "False"
}

func formatCmd(script string, args []string) string {
	if strings.HasPrefix(script, "//") {
		script = resolveLocation(script)
		script = fmt.Sprintf("$(locations %s)", script)
	}

	var argsStr strings.Builder
	for _, arg := range args {
		// TODO: handle $target_out_dir and $target_gen_dir.

		argsStr.WriteString(" ")
		argsStr.WriteString(arg)
	}

	return fmt.Sprintf("./%s%s", script, argsStr.String())
}

func mergeSlices(args ...[]string) []string {
	var size int
	for _, arg := range args {
		size += len(arg)
	}

	merged := make([]string, 0, size)
	for _, arg := range args {
		merged = append(merged, arg...)
	}

	return merged
}

func ruleName(name string) string {
	return name[strings.LastIndex(name, ":")+1:]
}

func resolveLocations(v []string) []string {
	out := make([]string, 0, len(v))
	for _, s := range v {
		if strings.HasPrefix(s, "//out.gn") {
			s = strings.TrimPrefix(s, "//")
		} else {
			_, file := filepath.Split(s)
			if file == "preparse-data-format.h" {
				continue
			}

			s = resolveLocation(s)
		}

		out = append(out, s)
	}

	return out
}

func resolveLocation(s string) string {
	if strings.Contains(s, ":") {
		return s
	}

	dir, file := filepath.Split(s)
	return strings.TrimSuffix(dir, "/") + ":" + file
}

func uniqueSlice(v []string) []string {
	out := make([]string, 0, len(v))
outer:
	for _, s := range v {
		for _, os := range out {
			if os == s {
				continue outer
			}
		}

		out = append(out, s)
	}

	return out
}

func sliceOf(v ...string) []string {
	return v
}
