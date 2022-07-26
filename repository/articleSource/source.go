package articleSource

import (
	"io/ioutil"
	"linblog/model"
	"net/http"
)

var Article ArticleSource = NewGithub()

type ArticleSource interface {
	GetAuthors() ([]string, error)
	GetCategories(author string) ([]string, error)
	GetArticleNames(author, category string) ([]string, error)
	GetArticleHtml(author, category string, articleName string) (string, error)
	GetArticleInfo(author, category string, articleName string) (*model.Article, error)
	GetImageUrl(category string, articleName string, imageName string) (string, error)
}

func DoRequest(url string, token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "token "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()                // 这步是必要的，防止以后的内存泄漏，切记
	body, err := ioutil.ReadAll(response.Body) // 读取响应 body, 返回为 []byte
	if err != nil {
		return nil, err
	}
	return body, nil
}

type article struct {
	IsTop     bool   `json:"is_top"`
	PubTime   string `json:"publish_time"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
}
