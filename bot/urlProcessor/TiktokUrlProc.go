package urlProcessor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"recipebot/cloud"
	"recipebot/environment"
	"recipebot/urlextract"
	"strconv"
	"time"
)

const (
	extenstion = ".mp4"
)

type TiktokUrlProc struct {
	config *environment.TiktokProcConfig
	cloud  cloud.Cloud
}

func NewTiktokProc(cnf *environment.TiktokProcConfig, cloud cloud.Cloud) (*TiktokUrlProc, error) {
	p := new(TiktokUrlProc)
	p.cloud = cloud
	p.config = cnf
	if err := os.MkdirAll(cnf.WorkingDir, 0755); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *TiktokUrlProc) CanHandle(uType urlextract.UrlType) bool {
	return uType == urlextract.VIDEO_TIKTOK
}

func (p *TiktokUrlProc) Process(request *Request) (*Result, error) {

	apiResp, err := p.getVideoDownloadUrl(request)
	if err != nil {
		return nil, err
	}
	localFile, err := p.downloadVideo(apiResp)
	if err != nil {
		return nil, err
	}
	defer func() {
		rem := localFile.Name()
		err := os.Remove(rem)
		if err != nil {
			fmt.Println("Could not delete file", rem, err)
		}
	}()
	meta, err := p.cloud.Upload(getRemoteFileName(apiResp.Data.Id), localFile)
	if err != nil {
		return nil, err
	}

	return &Result{StoredUrl: meta.CloudUrl}, nil
}

func (p *TiktokUrlProc) getVideoDownloadUrl(request *Request) (*environment.ApiResponse, error) {
	hc := http.Client{}
	req, err := http.NewRequest("POST", p.config.DownloaderApi, nil)

	q := url.Values{}
	q.Add("url", request.Details.MatchedUrl.String())
	req.URL.RawQuery = q.Encode()

	response, err := hc.Do(req)

	if err != nil || response.StatusCode != 200 {
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	resp := new(environment.ApiResponse)
	if err := json.Unmarshal(bytes, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *TiktokUrlProc) downloadVideo(apiResp *environment.ApiResponse) (*os.File, error) {
	response, err := http.Get(apiResp.Data.Play)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	fileName := apiResp.Data.Id + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".mp4"
	file, err := os.Create(getFilePath(p.config.WorkingDir, fileName))
	writer := bufio.NewWriter(file)
	buff := make([]byte, 1024*256)
	for {
		n, err := response.Body.Read(buff)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		if n == 0 || errors.Is(err, io.EOF) {
			break
		}
		if _, err := writer.Write(buff[:n]); err != nil {
			return nil, err
		}
		if err := writer.Flush(); err != nil {
			return nil, err
		}
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getFilePath(dir string, fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, dir, fileName)
}

func getRemoteFileName(id string) string {
	return id + extenstion
}
