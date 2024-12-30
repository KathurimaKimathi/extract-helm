package helm

type Values struct {
	Image Image `yaml:"image"`
}

type Image struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}
