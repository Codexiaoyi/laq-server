package repository

import (
	"linblog/model"
	"linblog/repository/articleSource"
)

type CategoryRepository struct {
}

func (cate *CategoryRepository) GetCategories() []*model.Category {
	authors, err := articleSource.Article.GetAuthors()
	if err != nil {
		return nil
	}
	res := make([]*model.Category, 0)
	for _, author := range authors {
		cates, err := articleSource.Article.GetCategories(author)
		if err != nil {
			return nil
		}
		for index, cate := range cates {
			res = append(res, &model.Category{Id: index, Title: cate, Href: "category/" + cate})
		}
	}
	return res
}
