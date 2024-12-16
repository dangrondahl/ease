package ease

import (
	"bytes"
	"fmt"
	"os"

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

func (e *Ease) WithVersion(version string) *Ease {
	e.Version = version
	return e
}

func (e *Ease) WithGitOps(owner, repo string) *Ease {
	e.SetGitOps(owner, repo)
	return e
}

func (e *Ease) WithTemplate(url, ref, path string) *Ease {
	e.SetTemplate(url, ref, path)
	return e
}

func (e *Ease) SetGitOps(owner, repo string) {
	e.GitOps.Owner = owner
	e.GitOps.Repo = repo
}

func (e *Ease) SetTemplate(url, ref, path string) {
	e.Template.FromGit.URL = url
	e.Template.FromGit.Ref = ref
	e.Template.FromGit.Path = path
}

// CreateEaseFile creates a new ease.yaml file at the specified path using an embedded YAML template.
// If the file already exists, it returns an error.
//
// Parameters:
//   - path: The file path where the ease.yaml file will be created.
//
// Returns:
//   - error: An error if the file already exists, if the template parsing fails, if the file creation fails,
//     or if writing to the file fails.
func (e *Ease) CreateEaseFile(path string) error {

	// Check if the file already exists
	if _, err := os.Stat(path); err == nil {
		return errors.New("ease.yaml already exists")
	}

	buf := bytes.Buffer{}

	if err := Render(&buf, e); err != nil {
		return fmt.Errorf("failed to render ease.yaml: %w", err)
	}

	// write the buffer to the file
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write ease.yaml: %w", err)
	}

	return nil
}
