package controllers

import (
	"net/http"

	"github.com/byuoitav/hateoas"
	"github.com/labstack/echo"
)

// Root offers HATEOAS links from the root endpoint of the API
func Root(c echo.Context) error {
	hateoasObject := hateoas.GetInfo()

	links, err := hateoas.AddLinks(c.Path(), []string{})
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Could not add HATEOAS")
	}

	hateoasObject.Links = links

	return c.JSON(http.StatusOK, hateoasObject)
}
