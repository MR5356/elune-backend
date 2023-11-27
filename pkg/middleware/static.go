package middleware

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strings"
)

func Static(prefix string, fs *StaticFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if prefix != "" {
		fileServer = http.StripPrefix(prefix, fileServer)
	}
	return func(ctx *gin.Context) {
		if fs.Exists(prefix, ctx.Request.URL.Path) {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
			ctx.Abort()
		}
	}
}

type StaticFileSystem struct {
	fs   http.FileSystem
	root string
}

func (s *StaticFileSystem) Open(name string) (http.File, error) {
	openPath := path.Join(s.root, name)
	logrus.Debugf("openPath: %s", openPath)
	return s.fs.Open(openPath)
}

func (s *StaticFileSystem) Exists(prefix string, filepath string) bool {
	logrus.Debugf("filepath: %s, prefix: %s", filepath, prefix)
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		var name string
		if p == "" {
			name = path.Join(s.root, p, "index.html")
		} else {
			name = path.Join(s.root, p)
		}
		if _, err := s.fs.Open(name); err != nil {
			return false
		}
		return true
	}
	return false
}

func NewStaticFileSystem(data embed.FS, root string) *StaticFileSystem {
	return &StaticFileSystem{
		fs:   http.FS(data),
		root: root,
	}
}
