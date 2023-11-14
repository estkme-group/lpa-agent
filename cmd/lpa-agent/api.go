package main

import (
	"errors"
	"github.com/esimclub/lpa-agent/lpac"
	"net/http"
	"strconv"
	"strings"
)
import "github.com/labstack/echo/v4"

func NewAPIHTTPHandler(cmdline *lpac.CommandLine) http.Handler {
	mux := echo.New()
	mux.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var h *echo.HTTPError
		if !errors.As(err, &h) {
			err = &echo.HTTPError{Message: strings.TrimSpace(err.Error())}
		}
		_ = c.JSON(code, err)
	}
	mux.GET("/", func(c echo.Context) error {
		info, err := cmdline.Info(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, info)
	})
	mux.PATCH("/", func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		var actual, expected *lpac.Information
		if actual, err = cmdline.Info(ctx); err != nil {
			return err
		}
		expected = new(lpac.Information)
		if err = c.Bind(expected); err != nil {
			return err
		}
		if expected.DefaultSMDP != actual.DefaultSMDP {
			if err = cmdline.SetDefaultSMDP(ctx, expected.DefaultSMDP); err != nil {
				return err
			}
		}
		return
	})
	mux.DELETE("/", func(c echo.Context) error {
		return cmdline.Purge(c.Request().Context())
	})
	profileGroup := mux.Group("/profile")
	profileGroup.GET("/", func(c echo.Context) error {
		profiles, err := cmdline.ListProfile(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, profiles)
	})
	profileGroup.POST("/download", func(c echo.Context) error {
		var profile lpac.DownloadProfile
		if err := c.Bind(&profile); err != nil {
			return err
		}
		return cmdline.DownloadProfile(c.Request().Context(), &profile)
	})
	profileGroup.GET("/{iccid}", func(c echo.Context) error {
		profile, err := cmdline.SpecificProfile(c.Request().Context(), c.Param("iccid"))
		if err != nil {
			return err
		}
		if profile == nil {
			return echo.ErrNotFound
		}
		return c.JSON(http.StatusOK, profile)
	})
	profileGroup.PUT("/{iccid}", func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		iccid := c.Param("iccid")
		var actual, expected *lpac.Profile
		if actual, err = cmdline.SpecificProfile(ctx, iccid); err != nil {
			return err
		}
		expected = new(lpac.Profile)
		if err = c.Bind(expected); err != nil {
			return err
		}
		if expected.ProfileState != actual.ProfileState {
			switch expected.ProfileState {
			case lpac.ProfileStateEnabled:
				err = cmdline.EnableProfile(ctx, iccid)
			case lpac.ProfileStateDisabled:
				err = cmdline.DisableProfile(ctx, iccid)
			}
			if err != nil {
				return err
			}
		}
		if expected.ProfileName != actual.ProfileName {
			if err = cmdline.SetProfileName(ctx, iccid, expected.ProfileName); err != nil {
				return err
			}
		}
		return
	})
	profileGroup.DELETE("/{iccid}", func(c echo.Context) error {
		return cmdline.DeleteProfile(c.Request().Context(), c.Param("iccid"))
	})
	notificationGroup := mux.Group("/notification")
	notificationGroup.GET("/", func(c echo.Context) error {
		notifications, err := cmdline.ListNotification(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, notifications)
	})
	notificationGroup.GET("/{index}", func(c echo.Context) error {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			return err
		}
		notification, err := cmdline.SpecificNotification(c.Request().Context(), index)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &notification)
	})
	notificationGroup.POST("/{index}", func(c echo.Context) error {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			return err
		}
		return cmdline.ProcessNotification(c.Request().Context(), index)
	})
	notificationGroup.DELETE("/{index}", func(c echo.Context) error {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			return err
		}
		return cmdline.RemoveNotification(c.Request().Context(), index)
	})
	return mux
}
