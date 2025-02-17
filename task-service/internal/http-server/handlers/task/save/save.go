package save

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	task_model "task-service/internal/model/task"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"

	resp "task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
)

type Request struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.3	 --name=URLSaver
type TaskSaver interface {
	SaveTask(taskToSave task_model.Task) error
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		var task task_model.Task
		task.Name = req.Name
		task.Description = req.Description
		task.Status = req.Status

		err = taskSaver.SaveTask(task)

		if err != nil {
			log.Error("failed to add task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("task added", slog.String("name", task.Name))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
