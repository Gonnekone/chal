package update

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	resp "github.com/Gonnekone/challenge/internal/lib/api/response"
	"github.com/Gonnekone/challenge/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CarUpdater interface {
	UpdateCar(ctx context.Context, log *slog.Logger, id int64, updates map[string]interface{}) error
}

type Response struct {
	resp.Response
}

type Request struct {
	Id     int64  `json:"id"`
	RegNum string `json:"regNum,omitempty"`
    Mark   string `json:"mark,omitempty"`
    Model  string `json:"model,omitempty"`
    Year   int64  `json:"year,omitempty"`
}

// @Summary Update car
// @Description Updates the details of a car with the provided ID.
// @Tags cars
// @Accept json
// @Produce json
// @Param body body Request true "Request body containing updated car details"
// @Success 200 "ok"
// @Failure 400 "client error"
// @Router / [patch]
func New(log *slog.Logger, carUpdater CarUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.update.New"

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

		err = carUpdater.UpdateCar(r.Context(), log, req.Id, params)
		if err != nil {
			log.Error("failed to update car with id: " + strconv.FormatInt(req.Id, 10), sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update car with id: " + strconv.FormatInt(req.Id, 10)))
			return
		}

		log.Info("car updated")

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
