package delete

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"

	resp "task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
)

type TaskDeleter interface {
	DeleteTask(id int) error
}

func New(log *slog.Logger, taskDeleter TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := parseIDFromRequest(r)
		if err != nil {
			log.Error("invalid task id", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid task id"))
			return
		}

		err = taskDeleter.DeleteTask(id)
		if err != nil {
			log.Error("failed to delete task", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete task"))
			return
		}

		log.Info("task deleted", slog.Int("task_id", id))
		responseOK(w, r)
	}
}

func parseIDFromRequest(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("task id is missing")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid task id format")
	}

	return id, nil
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
