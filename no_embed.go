//go:build !embed_dicts

package lojban_password_gen

func GetEmbeddedDicts() (string, string, bool) {
	return "", "", false
}
