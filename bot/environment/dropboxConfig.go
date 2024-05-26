package environment

import "strings"

const (
	dropboxToken  = "DROPBOX_TOKEN"
	dropboxFolder = "DROPBOX_FOLDER"
)

type DropboxConfig struct {
	DropboxToken  string
	DropboxFolder string
}

func (dConf *DropboxConfig) Load(env map[string]any) error {
	var err error
	if dConf.DropboxToken, err = GetRequiredString(env, dropboxToken); err != nil {
		return err
	}
	if folder, err := GetRequiredString(env, dropboxFolder); err != nil {
		return err
	} else {
		dConf.DropboxFolder = startTrailSep(folder)
	}
	return nil
}

func startTrailSep(path string) string {
	separator := "/"
	if !strings.HasSuffix(path, separator) {
		path += separator
	}
	if !strings.HasPrefix(path, separator) {
		path = separator + path
	}
	return path
}
