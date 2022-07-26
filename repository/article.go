package repository

import (
	"linblog/model"
	"linblog/repository/articleSource"
)

type ArticleRepository struct {
}

func (article *ArticleRepository) GetAllArticles(page, size int) ([]*model.Article, int, error) {
	res := make([]*model.Article, 0)
	min_index, max_index := (page-1)*size, page*size
	totalLength := 0
	authors, err := articleSource.Article.GetAuthors()
	if err != nil {
		return nil, 0, err
	}
	for _, author := range authors {
		//先拿到所有的分类
		cates, err := articleSource.Article.GetCategories(author)
		if err != nil {
			return nil, 0, err
		}
		//按照所有分类拿到所有的文章名
		for _, cate := range cates {
			articles, err := articleSource.Article.GetArticleNames(author, cate)
			if err != nil {
				return nil, 0, err
			}
			articleLength := len(articles)
			//小于所需起始条目的直接continue
			if articleLength+totalLength < min_index {
				totalLength += articleLength
				continue
			}

			for _, article := range articles {
				if totalLength >= max_index {
					return res, totalLength, nil
				}
				//每个分类对应的文章
				if totalLength >= min_index && totalLength < max_index {
					newArticle, _ := articleSource.Article.GetArticleInfo(author, cate, article)
					res = append(res, newArticle)
				}
				totalLength++
			}
		}
	}
	return res, totalLength, nil
}

func (article *ArticleRepository) GetArticlesByCategory(cate string, page, size int) ([]*model.Article, int, error) {
	//按照所有分类拿到所有的文章名
	res := make([]*model.Article, 0)
	min_index, max_index := (page-1)*size, page*size
	totalLength := 0
	authors, err := articleSource.Article.GetAuthors()
	if err != nil {
		return nil, 0, err
	}
	for _, author := range authors {
		articles, err := articleSource.Article.GetArticleNames(author, cate)
		if err != nil {
			return nil, 0, err
		}
		for _, article := range articles {
			if totalLength >= max_index {
				return res, totalLength, nil
			}
			//每个分类对应的文章
			if totalLength >= min_index && totalLength < max_index {
				newArticle, _ := articleSource.Article.GetArticleInfo(author, cate, article)
				res = append(res, newArticle)
			}
			totalLength++
		}
	}
	return res, totalLength, nil
}

func (article *ArticleRepository) GetArticleContent(category, articleName string) (string, error) {
	authors, err := articleSource.Article.GetAuthors()
	if err != nil {
		return "", err
	}
	for _, author := range authors {
		content, err := articleSource.Article.GetArticleHtml(author, category, articleName)
		if err == nil && content != "" {
			return content, nil
		}
	}
	return "", err
}

func (article *ArticleRepository) GetArticleInfo(category string, articleName string) (*model.Article, error) {
	authors, err := articleSource.Article.GetAuthors()
	if err != nil {
		return nil, err
	}
	for _, author := range authors {
		info, err := articleSource.Article.GetArticleInfo(author, category, articleName)
		if err == nil && info != nil {
			return info, nil
		}
	}
	return nil, err
}
