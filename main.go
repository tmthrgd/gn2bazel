package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
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
	out := flag.String("out", "BUILD", "the output path for the Bazel BUILD file")
	flag.Parse()

	outDir := flag.Arg(0)
	if outDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	infoJSON, err := Run(*dir, "gn", "desc", "--format=json", outDir, "*")
	if err != nil {
		log.Fatal(err)
	}

	var targets map[string]targetProperties
	if err := json.Unmarshal(infoJSON, &targets); err != nil {
		log.Fatalf("json parsing failed: %v", err)
	}

	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("creating output BUILD file failed: %v", err)
	}
	w := bufio.NewWriter(f)

	sortedTargets := make([]string, 0, len(targets))
	for name := range targets {
		sortedTargets = append(sortedTargets, name)
	}
	sort.Strings(sortedTargets)

	dataTargets := make(map[string]bool)
	for name, target := range targets {
		if target.Type == "copy" {
			dataTargets[name] = true
		}
	}

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

		deps, data := filterDeps(&target, dataTargets)
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

func filterDeps(target *targetProperties, dataTargets map[string]bool) (deps, data []string) {
	deps = make([]string, 0, len(target.Deps))
	for _, dep := range target.Deps {
		if dataTargets[dep] {
			data = append(data, dep)
		} else {
			deps = append(deps, dep)
		}
	}

	return deps, data
}
