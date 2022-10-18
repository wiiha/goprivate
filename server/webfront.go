package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed webfront
var webfront embed.FS

/*
fsWebBaseDir will navigate into the webfront
folder and then return the sub dir.
The purpose is to remove the `webfront/`
part from all paths.
*/
func fsLandigPage() http.FileSystem {
	webBaseDir, err := fs.Sub(webfront, "webfront")
	if err != nil {
		log.Fatalf("subbing fs dir: %v", err)
	}
	return http.FS(webBaseDir)
}

func fsAssetsDir() http.FileSystem {
	assetsDir, err := fs.Sub(webfront, "webfront/assets")
	if err != nil {
		log.Fatalf("subbing fs dir: %v", err)
	}
	return http.FS(assetsDir)
}
