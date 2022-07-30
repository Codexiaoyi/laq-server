package articleSource

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"linblog/config"
	"linblog/model"
	"linblog/utils"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

var _ ArticleSource = &Github{}

type Github struct {
	access_token string
	owner        string
	repo         string
}

type GithubResponse struct {
	Content     string `json:"content"`
	DownloadUrl string `json:"download_url"`
	Encoding    string `json:"encoding"`
	HtmlUrl     string `json:"html_url"`
	GitUrl      string `json:"git_url"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}

type Links struct {
	Git  string `json:"git"`
	Self string `json:"self"`
	Html string `json:"html"`
}

func NewGithub() ArticleSource {
	return &Github{
		access_token: config.AccessToken,
		owner:        config.Owner,
		repo:         config.Repo,
	}
}

func (g *Github) GetAuthors() ([]string, error) {
	response, err := DoRequest(fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/", g.owner, g.repo), g.access_token)
	if err != nil {
		return nil, err
	}
	authors := []*GithubResponse{}
	fmt.Println(string(response))
	err = json.Unmarshal([]byte(string(response)), &authors)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(authors))
	for _, author := range authors {
		res = append(res, author.Name)
	}
	return res, nil
}

func (g *Github) GetCategories(author string) ([]string, error) {
	response, err := DoRequest(fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", g.owner, g.repo, author), g.access_token)
	if err != nil {
		return nil, err
	}
	categories := []*GithubResponse{}
	err = json.Unmarshal([]byte(string(response)), &categories)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(categories))
	for _, category := range categories {
		if category.Name == "info.json" {
			continue
		}
		res = append(res, category.Name)
	}
	return res, nil
}

func (g *Github) GetArticleNames(author, category string) ([]string, error) {
	response, err := DoRequest(fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s/%s", g.owner, g.repo, author, category), g.access_token)
	if err != nil {
		return nil, err
	}
	articles := []*GithubResponse{}
	err = json.Unmarshal([]byte(string(response)), &articles)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(articles))
	for _, article := range articles {
		res = append(res, article.Name)
	}
	return res, nil
}

func (g *Github) GetArticleHtml(author, category string, articleName string) (string, error) {
	//直接获取相应的文章markdown字符串
	response, err := DoRequest(fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s/%s/%s/%s.md", g.owner, g.repo, author, category, articleName, articleName), g.access_token)
	if err != nil {
		return "", err
	}
	article := &GithubResponse{}
	err = json.Unmarshal([]byte(string(response)), &article)
	if err != nil {
		return "", err
	}
	mdStr, err := base64.StdEncoding.DecodeString(article.Content)
	if err != nil {
		return "", err
	}

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderFunc := func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		switch node.(type) {
		// case *ast.Image:
		// 	//这里将md中的图片在转换为html的时候替换
		// 	image := node.(*ast.Image)
		// 	name := string(image.Destination)
		// 	url, err := g.GetImageUrl(category, articleName, name)
		// 	if err != nil {
		// 		return ast.GoToNext, false
		// 	}
		// 	image.Destination = []byte(url)
		case *ast.CodeBlock:
			code := node.(*ast.CodeBlock)
			newCode, err := utils.ToGoSyntaxHighlight(code.Literal)
			if err == nil {
				io.WriteString(w, newCode)
				//跳过这个node
				return ast.GoToNext, true
			}
		}
		return ast.GoToNext, false
	}

	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: renderFunc}
	renderer := html.NewRenderer(opts)
	html := markdown.ToHTML(mdStr, nil, renderer)
	return string(html), nil
}

func (g *Github) GetImageUrl(category string, articleName string, imageName string) (string, error) {
	return fmt.Sprintf("https://raw.githubusercontent.com/Codexiaoyi/laq-articles/main/%s/%s/%s", category, articleName, imageName), nil
}

func (g *Github) GetArticleInfo(author, category string, articleName string) (*model.Article, error) {
	response, err := DoRequest(fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s/%s/%s/info.json", g.owner, g.repo, author, category, articleName), g.access_token)
	if err != nil {
		return nil, err
	}
	info := &GithubResponse{}
	err = json.Unmarshal([]byte(string(response)), &info)
	if err != nil {
		return nil, err
	}
	mdStr, err := base64.StdEncoding.DecodeString(info.Content)
	if err != nil {
		return nil, err
	}

	article := &article{}
	err = json.Unmarshal(mdStr, &article)
	if err != nil {
		return nil, err
	}
	newArticle := &model.Article{
		IsTop:         article.IsTop,
		Banner:        utils.GetRandomImageUrl(),
		IsHot:         true,
		PubTime:       article.PubTime,
		Title:         articleName,
		Summary:       article.Summary,
		Publisher:     article.Publisher,
		Category:      category,
		ViewsCount:    1000,
		CommentsCount: 100,
	}
	return newArticle, nil
}
