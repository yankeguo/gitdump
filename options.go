package gitdump

type Options struct {
	Dir      string `yaml:"dir" default:"." validate:"required"`
	Accounts []struct {
		Vendor   string `yaml:"vendor" validate:"required"`
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"accounts"`
}
