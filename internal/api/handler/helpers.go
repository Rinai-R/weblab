package handler

import (
	"errors"
	"net/http"
	"strconv"
	"weblab/internal/dao"
	"weblab/internal/middleware"
	"weblab/internal/service"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

func userIDFromContext(c *gin.Context) (int64, bool) {
	value, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

func pathInt64(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}

func queryInt(c *gin.Context, name string, defaultValue int) (int, error) {
	raw := c.Query(name)
	if raw == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(raw)
}

func queryInt64(c *gin.Context, name string, defaultValue int64) (int64, error) {
	raw := c.Query(name)
	if raw == "" {
		return defaultValue, nil
	}
	return strconv.ParseInt(raw, 10, 64)
}

func handleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		utils.Fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrUnauthorized):
		utils.Fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrForbidden):
		utils.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, dao.ErrNotFound):
		utils.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, dao.ErrDuplicated):
		utils.Fail(c, http.StatusConflict, err.Error())
	default:
		utils.Fail(c, http.StatusInternalServerError, "internal error")
	}
}
