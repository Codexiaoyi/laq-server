package controllers

import (
	"fmt"
	"linblog/model"
	"linblog/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/Codexiaoyi/linweb"
	"github.com/Codexiaoyi/linweb/interfaces"
)

type ArticleController struct {
	ArticleRepo *repository.ArticleRepository
}

//[GET("/articles")]
func (a *ArticleController) GetHomeArticles(c interfaces.IContext) {
	category := c.Request().Query("category")
	page, _ := strconv.Atoi(c.Request().Query("page"))
	size, _ := strconv.Atoi(c.Request().Query("size"))
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 5
	}

	//缓存
	cacheKey := fmt.Sprintf("articles_%s_%v_%v", category, page, size)
	if cache, ok := linweb.Cache.Get(cacheKey); ok {
		Response(c, http.StatusOK, cache)
		return
	}

	total, err := 0, (error)(nil)
	articles := make([]*model.Article, 0)
	if category != "" {
		articles, total, err = a.ArticleRepo.GetArticlesByCategory(category, page, size)
	} else {
		articles, total, err = a.ArticleRepo.GetAllArticles(page, size)
	}
	if err != nil {
		fmt.Println("Error by GetHomeArticles, error message is ", err)
		Response(c, http.StatusInternalServerError, nil)
		return
	}
	response := &articleListResponseDto{
		Total:       total,
		Items:       make([]*articleResponseDto, 0, len(articles)),
		HasNextPage: total >= page*size,
		Page:        page,
		Size:        size,
	}
	for _, article := range articles {
		dto := &articleResponseDto{}
		err := linweb.NewModel(article).MapToByFieldName(dto).ModelError()
		if err != nil {
			fmt.Println("Error by GetHomeArticles, error message is ", err)
			Response(c, http.StatusInternalServerError, nil)
			return
		}
		response.Items = append(response.Items, dto)
	}

	//添加缓存，1小时过期
	linweb.Cache.AddWithExpire(cacheKey, response, time.Hour*1)

	Response(c, http.StatusOK, response)
}

//[GET("/article/info/:cate/:title")]
func (a *ArticleController) GetArticleInfo(c interfaces.IContext) {
	cate := c.Request().Param("cate")
	title := c.Request().Param("title")

	//缓存
	cacheKey := fmt.Sprintf("article_info_%s_%s", cate, title)
	if cache, ok := linweb.Cache.Get(cacheKey); ok {
		Response(c, http.StatusOK, cache)
		return
	}

	info, err := a.ArticleRepo.GetArticleInfo(cate, title)
	if err != nil {
		Response(c, http.StatusInternalServerError, nil)
		return
	}
	dto := &articleResponseDto{}
	err = linweb.NewModel(info).MapToByFieldName(dto).ModelError()
	if err != nil {
		Response(c, http.StatusInternalServerError, nil)
		return
	}

	linweb.Cache.AddWithExpire(cacheKey, dto, time.Hour*1)
	Response(c, http.StatusOK, dto)
}

//[GET("/article/:cate/:title")]
func (a *ArticleController) GetArticleContent(c interfaces.IContext) {
	cate := c.Request().Param("cate")
	title := c.Request().Param("title")

	//缓存
	cacheKey := fmt.Sprintf("article_content_%s_%s", cate, title)
	if cache, ok := linweb.Cache.Get(cacheKey); ok {
		Response(c, http.StatusOK, cache)
		return
	}

	article, err := a.ArticleRepo.GetArticleContent(cate, title)
	if err != nil {
		Response(c, http.StatusInternalServerError, nil)
		return
	}

	linweb.Cache.AddWithExpire(cacheKey, article, time.Hour*1)
	Response(c, http.StatusOK, article)
}

//[POST("/articles/change")]
func (a *ArticleController) ArticleChangedEvent(c interfaces.IContext) {

}

type articleListResponseDto struct {
	Total       int                   `json:"total"`
	Items       []*articleResponseDto `json:"items"`
	HasNextPage bool                  `json:"hasNextPage"`
	Page        int                   `json:"page"`
	Size        int                   `json:"size"`
}

type articleResponseDto struct {
	Id            int    `json:"id"`
	IsTop         bool   `json:"isTop"`
	Banner        string `json:"banner"`
	IsHot         bool   `json:"isHot"`
	PubTime       string `json:"pubTime"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Category      string `json:"category"`
	Publisher     string `json:"publisher"`
	ViewsCount    int    `json:"viewsCount"`
	CommentsCount int    `json:"commentsCount"`
}
