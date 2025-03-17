package update

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
	Id          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Author      string `json:"author,omitempty" validate:"required"`
	Type        string `json:"type,omitempty"`
}

type TaskUpdater interface {
	UpdateTask(taskToUpdate task_model.Task) error
}

func New(log *slog.Logger, taskUpdater TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.update.New"

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

		task := task_model.Task{
			Id:          req.Id,
			Name:        req.Name,
			Description: req.Description,
			Status:      req.Status,
			Author:      req.Author,
			Type:        req.Type,
		}

		err = taskUpdater.UpdateTask(task)
		if err != nil {
			log.Error("failed to update task", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update task"))
			return
		}

		log.Info("task updated", slog.String("name", task.Name))
		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
