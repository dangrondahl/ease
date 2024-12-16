package ease

import (
	"fmt"
	"html/template"
	"os"

	"github.com/dangrondahl/ease/internal/static"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

type Ease struct {
	// Version is the version of the ease.yaml file.
	Version string `yaml:"version"`
	// GitOps contains the owner and repo of the git repository.
	GitOps struct {
		Owner string `yaml:"owner"`
		Repo  string `yaml:"repo"`
	} `yaml:"gitops"`

	// Template contains the URL, ref, and path of the template repository.
	Template struct {
		FromGit struct {
			URL  string `yaml:"url"`
			Ref  string `yaml:"ref"`
			Path string `yaml:"path"`
		} `yaml:"fromGit"`
	} `yaml:"template"`
}

func New() *Ease {
	return &Ease{
		Version: "v1",
	}
}

// New config from a file
func FromFile(path string) (*Ease, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	cfg := &Ease{}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	return cfg, nil
}

func (c *Ease) WithVersion(version string) *Ease {
	c.Version = version
	return c
}

func (c *Ease) WithGitOps(owner, repo string) *Ease {
	c.SetGitOps(owner, repo)
	return c
}

func (c *Ease) WithTemplate(url, ref, path string) *Ease {
	c.SetTemplate(url, ref, path)
	return c
}

func (c *Ease) SetGitOps(owner, repo string) {
	c.GitOps.Owner = owner
	c.GitOps.Repo = repo
}

func (c *Ease) SetTemplate(url, ref, path string) {
	c.Template.FromGit.URL = url
	c.Template.FromGit.Ref = ref
	c.Template.FromGit.Path = path
}

func (c *Ease) CreateEaseFile(path string) error {

	// Check if the file already exists
	if _, err := os.Stat(path); err == nil {
		return errors.New("ease.yaml already exists")
	}

	// Parse the embedded YAML template
	tmpl, err := template.New("ease.yaml").Parse(string(static.EaseExampleConfig))
	if err != nil {
		return fmt.Errorf("failed to parse embedded template: %w", err)
	}

	// Create the output file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create ease.yaml: %w", err)
	}
	defer file.Close()

	// Execute the template with the provided data
	if err := tmpl.Execute(file, c); err != nil {
		return fmt.Errorf("failed to write ease.yaml: %w", err)
	}

	return nil
}
