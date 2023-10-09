package response

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/gin-gonic/gin"
)

//go:embed html/*.html
var HTMLtemplates embed.FS

// Sets HTML templates from FS
func LoadHTMLTemplates(router *gin.Engine) {
	subFS, _ := fs.Sub(HTMLtemplates, "html")
	html := template.Must(template.ParseFS(subFS, "*.html"))
	router.SetHTMLTemplate(html)
}
