package httphandler

import (
	"github.com/hlf-mipt/saac-v2-client/internal/service/hlf"
	"github.com/hlf-mipt/saac-v2-core/pkg/model"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type HttpHandler struct {
	log *logrus.Entry
	hlf *hlf.HLFService
}

func NewHttpHandler(log *logrus.Entry, hlf *hlf.HLFService) *HttpHandler {
	return &HttpHandler{
		log: log,
		hlf: hlf,
	}
}

func (ctrl *HttpHandler) CreateAsset(e echo.Context) error {
	user := e.Request().Header.Get("user")
	req := new(model.Asset)
	err := e.Bind(req)
	if err != nil {
		ctrl.log.Errorf("failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input data")
	}

	err = ctrl.hlf.CreateAsset(user, req)
	if err != nil {
		ctrl.log.Errorf("failed to create asset: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, req)
}

func (ctrl *HttpHandler) UpdateAsset(e echo.Context) error {
	user := e.Request().Header.Get("user")
	req := new(model.Asset)
	err := e.Bind(req)
	if err != nil {
		ctrl.log.Errorf("failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input data")
	}

	err = ctrl.hlf.UpdateAsset(user, req)
	if err != nil {
		ctrl.log.Errorf("failed to update asset: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, req)
}

func (ctrl *HttpHandler) ReadAsset(e echo.Context) error {
	user := e.Request().Header.Get("user")
	id, _ := strconv.Atoi(e.Request().Header.Get("id"))

	item, err := ctrl.hlf.ReadAsset(user, id)
	if err != nil {
		ctrl.log.Errorf("failed to read asset: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, item)
}
