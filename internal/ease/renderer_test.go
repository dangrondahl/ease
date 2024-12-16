package ease_test

import (
	"bytes"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/dangrondahl/ease/internal/ease"
)

func TestRender(t *testing.T) {

	e := ease.Ease{
		Version: "v1",
		GitOps: ease.GitOps{
			URL: "https://github.com/owner/gitops-repo.git",
		},
		Template: ease.Template{
			FromGit: ease.FromGit{
				URL:  "https://github.com/owner/template-repo.git",
				Ref:  "ref",
				Path: "path",
			},
		},
	}

	approvals.UseFolder("testdata")

	t.Run("it renders a new ease.yaml file", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := ease.Render(&buf, &e); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}
