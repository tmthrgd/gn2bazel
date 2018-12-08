package main

// targetProperties is the schema output from gn desc --format=json.
//
// It's described here: https://gn.googlesource.com/gn/+/master/tools/gn/desc_builder.cc
type targetProperties struct {
	Type                      string // matching Target::GetStringForOutputType
	Toolchain                 string
	Visibility                []string
	TestOnly                  bool
	CheckIncludes             bool
	AllowCircularIncludesFrom []string
	Sources                   []string
	Public                    interface{} // either string("*") or []interface{} -> []string
	Inputs                    []string
	Configs                   []string
	PublicConfigs             []string
	AllDependentConfigs       []string
	Script                    string
	Args                      []string
	Depfile                   string
	Outputs                   []string
	Arflags                   []string
	Asmflags                  []string
	Cflags                    []string
	CflagsC                   []string
	CflagsCC                  []string `json:"clfags_cc"`
	CflagsObjC                []string `json:"clfags_objc"`
	CflagsObjCC               []string `json:"clfags_objcc"`
	Defines                   []string
	IncludeDirs               []string
	PrecompiledHeader         string
	PrecompiledSource         string
	Deps                      []string
	Libs                      []string
	LibDirs                   []string

	// Optionally, if "what" is specified while generating
	// description, two other properties can be requested that are not
	// included by default.
	RuntimeDeps   []string
	SourceOutputs []map[string][]string
}
