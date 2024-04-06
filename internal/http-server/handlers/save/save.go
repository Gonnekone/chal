package save

import (
	"context"
	"log/slog"
	"net/http"

	resp "github.com/Gonnekone/challenge/internal/lib/api/response"
	"github.com/Gonnekone/challenge/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CarSaver interface {
	SaveCar(ctx context.Context, log *slog.Logger, regNums []string) error
}

type Request struct {
	RegNums []string `json:"regNums" example:"AAA111, BBB222, CCC333"`
}

type Response struct {
	resp.Response
}

// @Summary Save cars
// @Description Saves a list of cars with the provided registration numbers.
// @Tags cars
// @Accept json
// @Produce json
// @Param body body Request true "Request body containing registration numbers"
// @Success 200 "ok"
// @Failure 400 "client error"
// @Router / [post]
func New(log *slog.Logger, carSaver CarSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

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

		if req.RegNums == nil {
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err = carSaver.SaveCar(r.Context(), log, req.RegNums)

		if err != nil {
			log.Error("failed to save regNum", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save regNum"))

			return
		}

		log.Info("cars added")

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
