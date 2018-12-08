package main

import "text/template"

var templates = template.New("gn2bazel").Funcs(template.FuncMap{
	"print_slice":         printStringSlice,
	"filter_sources":      filterHeadersFromSources,
	"to_string_slice":     toStringSlice,
	"to_bazel_visibility": toBazelVisibility,
	"print_bool":          printBool,
	"format_cmd":          formatCmd,
	"merge_slices":        mergeSlices,
	"rule_name":           ruleName,
	"resolve_locations":   resolveLocations,
	"unique":              uniqueSlice,
})

// action -> genrule
func init() {
	template.Must(templates.New("genrule").Parse(`{{/**/ -}}
genrule(
	name = {{rule_name .Name | printf "%q"}},
	srcs = [{{merge_slices .Inputs .Sources .Deps .Data | resolve_locations | print_slice}}],
	outs = [{{resolve_locations .Outputs | print_slice}}],
	cmd = {{format_cmd .Script .Args | printf "%q"}},
	visibility = [{{to_bazel_visibility .Visibility | print_slice}}],
	testonly = {{print_bool .TestOnly}},
)
`))
}

// copy -> filegroup
func init() {
	// TODO: handle outputs, $target_out_dir, $target_gen_dir and
	// {{source*}} expansion.

	template.Must(templates.New("filegroup").Parse(`{{/**/ -}}
filegroup(
	name = {{rule_name .Name | printf "%q"}},
	srcs = [{{resolve_locations .Sources | print_slice}}],
	data = [{{print_slice .Data}}],
	visibility = [{{to_bazel_visibility .Visibility | print_slice}}],
	testonly = {{print_bool .TestOnly}},
)
`))
}

// executable -> cc_binary
func init() {
	template.Must(templates.New("cc_binary").Parse(`{{/**/ -}}
cc_binary(
	name = {{rule_name .Name | printf "%q"}},
	deps = [{{print_slice .Deps}}],
	srcs = [{{resolve_locations .Sources | unique | print_slice}}],
	data = [{{print_slice .Data}}],
	copts = [{{merge_slices .Cflags .Asmflags | print_slice}}],
	defines = [{{print_slice .Defines}}],
	includes = [{{print_slice .IncludeDirs}}],
	linkopts = [{{print_slice .Arflags}}],
	visibility = [{{to_bazel_visibility .Visibility | print_slice}}],
	testonly = {{print_bool .TestOnly}},
)
`))
}

// group -> cc_library?
func init() {
	// TODO: handle data_deps, deps, public_deps.
	// TODO: implement
}

// source_set -> cc_library
func init() {
	// TODO: implement
}

// static_library -> cc_library
func init() {
	template.Must(templates.New("cc_library").Parse(`{{/**/ -}}
cc_library(
	name = {{rule_name .Name | printf "%q"}},
	deps = [{{print_slice .Deps}}],
	srcs = [{{unique .Sources | resolve_locations | print_slice}}],
	data = [{{print_slice .Data}}],
	hdrs = [
{{- if eq (print .Public) "*" -}}
	{{filter_sources .Sources | resolve_locations | unique | print_slice}}
{{- else -}}
	{{to_string_slice .Public | resolve_locations | unique | print_slice}}
{{- end -}}
	],
	copts = [{{merge_slices .Cflags .Asmflags | print_slice}}],
	defines = [{{print_slice .Defines}}],
	includes = [{{print_slice .IncludeDirs}}],
	linkopts = [{{print_slice .Arflags}}],
	visibility = [{{to_bazel_visibility .Visibility | print_slice}}],
	testonly = {{print_bool .TestOnly}},
)
`))
}

// cc_test
func init() {
	template.Must(templates.New("cc_test").Parse(`{{/**/ -}}
cc_test(
	name = {{rule_name .Name | printf "%q"}},
	deps = [{{print_slice .Deps}}],
	srcs = [{{resolve_locations .Sources | print_slice}}],
	data = [{{print_slice .Data}}],
	hdrs = [
{{- if eq (print .Public) "*" -}}
	{{filter_sources .Sources | resolve_locations | unique | print_slice}}
{{- else -}}
	{{to_string_slice .Public | resolve_locations | unique | print_slice}}
{{- end -}}
	],
	includes = [{{print_slice .IncludeDirs}}],
	visibility = [{{to_bazel_visibility .Visibility | print_slice}}],
	testonly = {{print_bool .TestOnly}},
)
`))
}
