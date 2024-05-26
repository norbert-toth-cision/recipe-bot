package cloud

import (
	"github.com/NightMan-1/go-dropbox"
	"github.com/NightMan-1/go-dropy"
	"io"
	"recipebot/environment"
)

type Dropbox struct {
	client *dropy.Client
	config *environment.DropboxConfig
}

func NewDropbox(config *environment.DropboxConfig) *Dropbox {
	dbx := new(Dropbox)
	dbx.client = dropy.New(dropbox.New(dropbox.NewConfig(config.DropboxToken)))
	dbx.config = config
	return dbx
}

func (dbx *Dropbox) Upload(remoteName string, reader io.Reader) (*Metadata, error) {
	remotePath := dbx.config.DropboxFolder + remoteName
	err := dbx.client.Upload(remotePath, reader)
	if err != nil {
		return nil, err
	}
	shareIn := dropbox.CreateSharedLinkInput{
		Path:     remotePath,
		ShortURL: true,
	}
	shareOut, err := dbx.client.Sharing.CreateSharedLink(&shareIn)
	if err != nil {
		return nil, err
	}
	return &Metadata{CloudUrl: shareOut.URL}, nil
}
