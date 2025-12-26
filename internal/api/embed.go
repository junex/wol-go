package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/static
var staticFiles embed.FS

//go:embed web/templates
var templateFiles embed.FS

// StaticFileHandler 静态文件处理器
func StaticFileHandler() http.FileSystem {
	fsys, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

// IndexHTML 索引页面内容
func IndexHTML() []byte {
	content, err := templateFiles.ReadFile("web/templates/index.html")
	if err != nil {
		panic(err)
	}
	return content
}
