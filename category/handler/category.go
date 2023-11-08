package handler

import (
	"context"

	"github.com/prometheus/common/log"
	"github.com/wagfog/gomdemo/common"
	"github.com/wagfog/gomdemo/domain/model"
	"github.com/wagfog/gomdemo/domain/service"
	"github.com/wagfog/gomdemo/proto/category"
)

type Category struct {
	CategoryDataService service.ICategoryDataService
}

func (c *Category) CreateCategory(ctx context.Context, request *category.CategoryRequest, resp *category.CreateCategoryResponse) error {
	category := &model.Category{}
	//
	err := common.SwapTo(request, category)
	if err != nil {
		return err
	}
	categoryId, err := c.CategoryDataService.AddCategory(category)
	if err != nil {
		return err
	}
	resp.Message = "add category success"
	resp.CategoryId = categoryId
	return nil
}

func (c *Category) UpdateCategory(ctx context.Context, request *category.CategoryRequest, resp *category.UpdateCategoryResponse) error {
	category := &model.Category{}
	err := common.SwapTo(resp, category)
	if err != nil {
		return err
	}
	err = c.CategoryDataService.UpdateCategory(category)
	if err != nil {
		return nil
	}
	resp.Message = "update category success"
	return nil
}

func (c *Category) DeleteCategory(ctx context.Context, request *category.DeleteCategoryRequest, resp *category.DeleteCategoryResponse) error {
	err := c.CategoryDataService.DeleteCategory(request.CategoryId)
	if err != nil {
		return nil
	}
	resp.Message = "delete success"
	return nil
}

func (c *Category) FindCategory(ctx context.Context, request *category.FindByNameRequest, resp *category.CategoryResponse) error {
	category, err := c.CategoryDataService.FindCategoryByName(request.CategoryName)
	if err != nil {
		return err
	}
	return common.SwapTo(category, resp)
}

// 根据分类ID查找分类
func (c *Category) FindCategoryByID(ctx context.Context, request *category.FindByIdRequest, response *category.CategoryResponse) error {
	category, err := c.CategoryDataService.FindCategoryByID(request.CategoryId)
	if err != nil {
		return err
	}
	return common.SwapTo(category, response)
}

func (c *Category) FindCategoryByLevel(ctx context.Context, request *category.FindByLevelRequest, response *category.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByLevel(request.Level)
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func (c *Category) FindCategoryByParent(ctx context.Context, request *category.FindByParentRequest, response *category.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByParent(request.ParentId)
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func (c *Category) FindAllCategory(ctx context.Context, request *category.FindAllRequest, response *category.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindAllCategory()
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func categoryToResponse(categorySlice []model.Category, response *category.FindAllResponse) {
	for _, cg := range categorySlice {
		cr := &category.CategoryResponse{}
		err := common.SwapTo(cg, cr)
		if err != nil {
			log.Error(err)
			break
		}
		response.Category = append(response.Category, cr)
	}
}
