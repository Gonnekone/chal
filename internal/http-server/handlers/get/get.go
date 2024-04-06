package get

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Gonnekone/challenge/internal/domain/models"
	resp "github.com/Gonnekone/challenge/internal/lib/api/response"
	"github.com/Gonnekone/challenge/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const limitDefault = "20"
const offsetDefault = "0"

type CarGetter interface {
	GetCar(ctx context.Context, log *slog.Logger, filters map[string]interface{}, limit, offset int) ([]models.Car, error)
}

type Request struct {
	RegNum string `json:"regNum,omitempty" example:"AAA111"`
    Mark   string `json:"mark,omitempty" example:"BMW"`
    Model  string `json:"model,omitempty" example:"X5"`
    Year   int64  `json:"year,omitempty" example:"2020"`
	OwnerName string `json:"ownerName,omitempty" example:"Ivan"`
	OwnerSurname string `json:"ownerSurname,omitempty" example:"Ivanov"`
	OwnerPatronymic string `json:"ownerPatronymic,omitempty" example:"Ivanovich"`
}

type Response struct {
	resp.Response
}

// @Summary Get cars
// @Description Retrieves a list of cars based on specified filters.
// @Tags cars
// @Param body body Request true "Request body containing filters"
// @Param limit path int true "limit"
// @Param offset path int true "offset"
// @Success 200 "ok"
// @Failure 400 "client error"
// @Router / [get]
func New(log *slog.Logger, carGetter CarGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request decoded", slog.Any("request", req))

		params := make(map[string]interface{})
		if req.RegNum != "" {
			params["regNum"] = req.RegNum
		}
		if req.Mark != "" {
			params["mark"] = req.Mark
		}
		if req.Model != "" {
			params["model"] = req.Model
		}
		if req.Year != 0 {
			params["year"] = req.Year
		}
		if req.OwnerName != "" {
			params["name"] = req.OwnerName
		}
		if req.OwnerSurname != "" {
			params["surname"] = req.OwnerSurname
		}
		if req.OwnerPatronymic != "" {
			params["patronymic"] = req.OwnerPatronymic
		}

		limitStr := r.URL.Query().Get("limit")
		if limitStr == "" {
			limitStr = limitDefault
		}
		offsetStr := r.URL.Query().Get("offset")
		if offsetStr == "" {
			offsetStr = offsetDefault
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Error("failed to parse limit", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse limit"))
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Error("failed to parse offset", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse offset"))
			return
		}

		cars, err := carGetter.GetCar(r.Context(), log, params, limit, offset)
		if err != nil {
			log.Error("failed to get cars", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get cars"))
			return
		}

		render.JSON(w, r, cars)
	}
}

// @Param regNum query string false "Car registration number"
// @Param mark query string false "Car mark"
// @Param model query string false "Car model"
// @Param year query int false "Car manufacturing year"
// @Param ownerName query string false "Owner's name"
// @Param ownerSurname query string false "Owner's surname"
// @Param ownerPatronymic query string false "Owner's patronymic"
// @Param limit query int false "Number of results to return (default: 20)"
// @Param offset query int false "Offset for paginated results (default: 0)"
