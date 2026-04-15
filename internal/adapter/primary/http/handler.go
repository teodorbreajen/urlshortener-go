package http

import (
	"errors"
	"net/http"

	"urlshortener/internal/domain/model"
	"urlshortener/internal/domain/service"
	"urlshortener/internal/port"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service port.URLService
	baseURL string
}

func NewURLHandler(service port.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type URLResponse struct {
	Code        string `json:"code"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	Clicks      int64  `json:"clicks"`
	CreatedAt   string `json:"created_at"`
}

type ListResponse struct {
	URLs   []*URLResponse `json:"urls"`
	Total  int            `json:"total"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
}

func (h *URLHandler) CreateURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	url, err := h.service.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		if errors.Is(err, model.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, toResponse(url, h.baseURL))
}

func (h *URLHandler) GetURL(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, toResponse(url, h.baseURL))
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.service.IncrementClicks(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func (h *URLHandler) ListURLs(c *gin.Context) {
	offset := 0
	limit := 10

	if o := c.Query("offset"); o != "" {
		if _, err := c.GetQuery("offset"); err {
			offset = 0
		}
	}
	if l := c.Query("limit"); l != "" {
		if _, err := c.GetQuery("limit"); err {
			limit = 10
		}
	}

	urls, err := h.service.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	responses := make([]*URLResponse, len(urls))
	for i, url := range urls {
		responses[i] = toResponse(url, h.baseURL)
	}

	c.JSON(http.StatusOK, ListResponse{
		URLs:   responses,
		Total:  len(responses),
		Offset: offset,
		Limit:  limit,
	})
}

func toResponse(url *model.URL, baseURL string) *URLResponse {
	return &URLResponse{
		Code:        url.Code,
		OriginalURL: url.OriginalURL,
		ShortURL:    baseURL + "/r/" + url.Code,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
