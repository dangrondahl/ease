package ease

import (
	"embed"
	"io"
	"text/template"
)

var (
	//go:embed templates/ease-example.yaml
	EaseTemplate embed.FS
)

func Render(w io.Writer, e *Ease) error {

	templ, err := template.ParseFS(EaseTemplate, "templates/ease-example.yaml")

	if err != nil {
		return err
	}

	if err := templ.Execute(w, e); err != nil {
		return err
	}

	return nil
}
