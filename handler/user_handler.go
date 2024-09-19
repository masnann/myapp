package handler

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
    "github.com/username/myapp/service"
)

type UserHandler struct {
    Service service.UserService
}

// NewUserHandler adalah konstruktor untuk UserHandler
func NewUserHandler(service service.UserService) *UserHandler {
    return &UserHandler{
        Service: service,
    }
}

// DeleteUser menangani HTTP DELETE /users/:id
func (h *UserHandler) DeleteUser(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid user id",
        })
    }

    err = h.Service.DeleteUser(id)
    if err != nil {
        if err.Error() == "user not found" {
            return c.JSON(http.StatusNotFound, map[string]string{
                "error": err.Error(),
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "internal server error",
        })
    }

    return c.NoContent(http.StatusNoContent)
}
