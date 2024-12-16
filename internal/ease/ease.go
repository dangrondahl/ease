package ease

import (
	"bytes"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"sigs.k8s.io/yaml"
)

type GitOps struct {
	URL string `yaml:"url"`
}

type FromGit struct {
	URL  string `yaml:"url"`
	Ref  string `yaml:"ref"`
	Path string `yaml:"path"`
}

type Template struct {
	FromGit FromGit `yaml:"fromGit"`
}

type Ease struct {
	// Version is the version of the ease.yaml file.
	Version string `yaml:"version"`
	// GitOps contains the owner and repo of the git repository.
	GitOps GitOps `yaml:"gitOps"`
	// Template contains the URL, ref, and path of the template repository.
	Template Template `yaml:"template"`
}

func New() *Ease {
	return &Ease{
		Version: "v1",
	}
}

func NewFromPrompt() *Ease {
	e := New()

	form := huh.NewForm(

		huh.NewGroup(
			huh.NewNote().Title("Create new ease config").Description("Happy to see you use _ease_.\n\nLet's configure the tool?\n\n").Next(true).NextLabel("Next"),
		),

		huh.NewGroup(
			huh.NewNote().Title("GitOps repository").Description("First, declare which GitOps repository\nThis service should release to:\n"),
			huh.NewInput().Title("URL of GitOps repository").Placeholder("https://github.com/owner/gitops-repo").Value(&e.GitOps.URL),
		),

		huh.NewGroup(
			huh.NewNote().Title("HelmRelease template repository").Description("Now, declare where your HelmRelease template is:\n"),
			huh.NewInput().Title("URL").Placeholder("Enter the URL of the repository").Value(&e.Template.FromGit.URL),
			huh.NewInput().Title("Ref").Placeholder("Enter the ref of the repository\n ref can be a branch, tag, or commit SHA\n").Value(&e.Template.FromGit.Ref),
			huh.NewInput().Title("Path").Placeholder("Enter the path to the HelmRelease template\n").Value(&e.Template.FromGit.Path),
		),
	)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	return e
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

func (e *Ease) WithGitOps(url string) *Ease {
	e.GitOps.URL = url
	return e
}

func (e *Ease) WithTemplate(url, ref, path string) *Ease {
	e.Template.FromGit.URL = url
	e.Template.FromGit.Ref = ref
	e.Template.FromGit.Path = path
	return e
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
		return fmt.Errorf("ease.yaml already exists")
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
