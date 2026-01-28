//go:build embed_dicts

package lojban_password_gen

import (
	_ "embed"
)

//go:embed gismu.txt
var embedGismu string

//go:embed cmavo.txt
var embedCmavo string

func GetEmbeddedDicts() (string, string, bool) {
	return embedGismu, embedCmavo, true
}
