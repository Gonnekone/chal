package del

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

type CarDeleter interface {
	DeleteCar(ctx context.Context, log *slog.Logger, id int64) error
}

type Response struct {
	resp.Response
}

// @Summary Delete a car
// @Description Deletes a car by its identifier.
// @Tags cars
// @Param id query int true "Car identifier"
// @Success 200 "ok"
// @Failure 400 "client error"
// @Router / [delete]
func New(log *slog.Logger, carDeleter CarDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			log.Error("id is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to parse id", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to parse id"))

			return
		}

		err = carDeleter.DeleteCar(r.Context(), log, int64(id))
		if err != nil {
			log.Error("failed to delete car", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("car deleted")

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
