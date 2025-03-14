package handlers

import (
	"net/http"
	"personal-notes-with-go/models"
	"personal-notes-with-go/repositories"
	"personal-notes-with-go/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo repositories.CategoryRepository
}

func NewCategoryHandler(repo repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.BindJSON(&cat); err != nil {
		utils.HandleBadRequestError(c, err)
		return
	}

	if cat.Name == "" {
		utils.HandleBadRequestError(c, utils.ErrCategoryNameEmpty)
		return
	}

	if err := h.repo.Create(&cat); err != nil {
		utils.HandleInternalServerError(c, err, "create category")
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		utils.HandleInternalServerError(c, err, "get categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := c.BindJSON(&cat); err != nil {
		utils.HandleBadRequestError(c, err)
		return
	}

	if cat.Name == "" {
		utils.HandleBadRequestError(c, utils.ErrCategoryNameEmpty)
		return
	}

	cat.ID = id
	if err := h.repo.Update(&cat); err != nil {
		utils.HandleInternalServerError(c, err, "update category")
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		utils.HandleInternalServerError(c, err, "delete category")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
