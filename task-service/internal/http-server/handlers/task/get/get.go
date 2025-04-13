package get

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
	task_model "task-service/internal/model/task"
)

type TaskGetter interface {
	GetTask(id int) (task_model.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get.New"

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

		task, err := taskGetter.GetTask(id)
		if err != nil {
			log.Error("failed to get task", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get task"))
			return
		}

		log.Info("task retrieved", slog.Any("task", task))
		render.JSON(w, r, resp.OKWithData(task))
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
