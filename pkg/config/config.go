package config

type RootConfig struct {
	Project ProjectConfig `yaml:"project"`
	Build   BuildConfig   `yaml:"build"`
}

type ProjectConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Module  string `yaml:"module"`
}

type BuildConfig struct {
	Binaries  []BuildBinaryConfig `yaml:"binaries"`
	Target    string              `yaml:"target"`
	Env       []string            `yaml:"env"`
	Flags     string              `yaml:"flags"`
	Ldflags   string              `yaml:"ldflags"`
	Platforms []string            `yaml:"platforms"`
}

type BuildBinaryConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
