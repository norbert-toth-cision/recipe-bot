package environment

const (
	downloaderApi   = "TIKTOK_DOWNLOADER_API"
	downloadWorkDir = "TIKTOK_DOWNLOADER_WORKDIR"
)

type TiktokProcConfig struct {
	DownloaderApi string
	WorkingDir    string
}

func (tConf *TiktokProcConfig) Load(env map[string]any) error {
	var err error
	tConf.DownloaderApi, err = GetRequiredString(env, downloaderApi)
	if err != nil {
		return err
	}
	tConf.WorkingDir, err = GetRequiredString(env, downloadWorkDir)
	if err != nil {
		return err
	}
	return nil
}

type ApiResponse struct {
	Data *Data
}

type Data struct {
	Id    string
	Title string
	Play  string
}
