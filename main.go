package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	outPath = flag.String("out", "", "the output path for the Bazel BUILD files")
	exclude = flag.String("exclude", "", "exclude targets matching this regexp")
)

func init() {
	log.SetFlags(log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s: %[1]s <out_dir>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	dir := flag.String("dir", "", "the path to run gn commands in, must be a checkout")
	flag.Parse()

	outDir := flag.Arg(0)
	if outDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	exclude, err := regexp.Compile(*exclude)
	if err != nil {
		log.Fatalf("compiling -exclude failed: %v", err)
	}

	if !strings.Contains(outDir, string(filepath.Separator)) {
		outDir = filepath.Join("out.gn", outDir)
	}

	infoJSON, err := Run(*dir, "gn", "desc", "--format=json", outDir, "*")
	if err != nil {
		log.Fatal(err)
	}

	var targets map[string]targetProperties
	if err := json.Unmarshal(infoJSON, &targets); err != nil {
		log.Fatalf("json parsing failed: %v", err)
	}

	paths := make(map[string][]string)
	for name := range targets {
		if exclude.MatchString(name) {
			continue
		}

		idx := strings.LastIndex(name, ":")
		if idx < 0 {
			log.Fatalf("invalid target name %q", name)
		}

		dir := name[:idx]
		paths[dir] = append(paths[dir], name)
	}

	for dir, rules := range paths {
		dir = strings.TrimPrefix(dir, "//")
		sort.Strings(rules)
		convert(targets, dir, rules)
	}
}

func convert(targets map[string]targetProperties, dir string, sortedTargets []string) {
	out := filepath.Join(*outPath, dir, "BUILD")
	f, err := os.Create(out)
	if err != nil {
		log.Fatalf("creating output BUILD file failed: %v", err)
	}
	w := bufio.NewWriter(f)

	for i, name := range sortedTargets {
		target := targets[name]

		// TODO: configs ?

		// Unhandled fields:
		//  toolchain
		//  check_includes
		//  allow_circular_includes_from
		//  configs
		//  public_configs
		//  all_dependent_configs
		//  depfile
		//  arflags
		//  asmflags
		//  cflags
		//  cflags_c
		//  clfags_cc
		//  cflags_objc
		//  clfags_objcc
		//  defines
		//  precompiled_header
		//  precompiled_source
		//  lib_dirs
		//  runtime_deps
		//  source_outputs

		// Unsupported types:
		//  action_foreach
		//  component
		//  shared_library

		var tmplName string
		switch target.Type {
		case "action":
			tmplName = "genrule"
		case "copy":
			tmplName = "filegroup"
		case "executable":
			tmplName = "cc_binary"
		case "group":
			tmplName = "cc_library"
		case "source_set":
			// I don't think this is strictly correct. See
			// https://chromium.googlesource.com/chromium/src/+/eca97f87e275a7c9c5b7f13a65ff8635f0821d46/tools/gn/docs/reference.md#source_set_Declare-a-source-set-target
			tmplName = "cc_library"
		case "static_library":
			tmplName = "cc_library"
		case "test":
			tmplName = "cc_test"
		default:
			log.Fatalf("unknown target type %q", target.Type)
		}

		if i > 0 {
			w.WriteString("\n")
		}

		fmt.Fprintf(w, "# %s %s\n", target.Type, name)

		deps, data := filterDeps(&target, targets)
		if err := templates.ExecuteTemplate(w, tmplName, struct {
			Name string
			*targetProperties

			Deps []string
			Data []string
		}{name, &target, deps, data}); err != nil {
			log.Fatalf("executing cc_library template failed: %v", err)
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatalf("flush BUILD file failed: %v", err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("closing BUILD file failed: %v", err)
	}
}

func filterDeps(target *targetProperties, targets map[string]targetProperties) (deps, data []string) {
	deps = make([]string, 0, len(target.Deps))
	for _, dep := range target.Deps {
		if isDataTarget(targets[dep]) {
			data = append(data, dep)
		} else {
			deps = append(deps, dep)
		}
	}

	return deps, data
}

func isDataTarget(target targetProperties) bool {
	return target.Type == "copy"
}
