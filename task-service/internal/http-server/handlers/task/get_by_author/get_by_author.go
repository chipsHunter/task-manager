package get_by_author

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
	task_model "task-service/internal/model/task"
)

type TaskGetter interface {
	GetUserTasks(author string) ([]task_model.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get_by_author.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		author := chi.URLParam(r, "author")
		if author == "" {
			log.Error("missing author", sl.Err(errors.New("author is missing")))
			render.JSON(w, r, resp.Error("author is required"))
			return
		}

		tasks, err := taskGetter.GetUserTasks(author)
		if err != nil {
			log.Error("failed to get tasks for author", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get tasks for author"))
			return
		}

		if len(tasks) == 0 {
			log.Info("no tasks found for author", slog.String("author", author))
		} else {
			log.Info("tasks retrieved for author", slog.String("author", author), slog.Int("task_count", len(tasks)))
		}

		render.JSON(w, r, resp.OKWithData(tasks))
	}
}
