package controllers

import (
	"net/http"

	"github.com/aasoft24/golara/wpkg/gola"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (c *HomeController) Index(ctx *gola.Context) {
	ctx.String(http.StatusOK, "Home Controller Index")
}
