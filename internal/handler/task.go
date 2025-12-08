package handler

import (
	"net/http"

	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/middleware"
	"golang-service-template/internal/service"

	"github.com/samber/do"

	"github.com/labstack/echo/v4"
)

type TaskController interface {
	Create() echo.HandlerFunc
	Find() echo.HandlerFunc
	FindByUserId() echo.HandlerFunc
	GetById() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type taskController struct {
	taskService service.TaskService
}

// Create implements TaskController.
func (tc *taskController) Create() echo.HandlerFunc {
	type task struct {
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		t := task{}

		// Use middleware's validator helper function for validation
		if err := middleware.ValidateRequest(c, &t); err != nil {
			return err
		}

		createdTask, err := tc.taskService.Create(c.Request().Context(), model.Task{
			Description: t.Description,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"meta": map[string]any{
				"status": http.StatusCreated,
			},
			"data": createdTask,
		})
	}
}

// Get implements TaskController.
func (t *taskController) Find() echo.HandlerFunc {
	return func(c echo.Context) error {

		tasks, err := t.taskService.Find(c.Request().Context())

		if err != nil {
			return err
		}

		return c.JSON(
			http.StatusOK,
			NewResponse().
				AddMeta("total", len(tasks)).
				AddMeta("status", http.StatusOK).
				SetData(tasks),
		)
	}
}

// Get implements TaskController.
func (t *taskController) FindByUserId() echo.HandlerFunc {
	return func(c echo.Context) error {

		userId := c.Get(middleware.ContextKeyUserId)

		userIdStr, _ := userId.(string)

		tasks, err := t.taskService.FindByUserId(c.Request().Context(), userIdStr)

		if err != nil {
			return err
		}

		return c.JSON(
			http.StatusOK,
			NewResponse().
				AddMeta("total", len(tasks)).
				AddMeta("status", http.StatusOK).
				SetData(tasks),
		)
	}
}

// GetTask implements TaskController.
func (tc *taskController) GetById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		// Get validator from middleware
		validate := middleware.GetValidator(c)
		if validate == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Validator not configured")
		}

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		task, err := tc.taskService.Get(c.Request().Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(
			http.StatusOK,
			NewResponse().
				AddMeta("status", http.StatusOK).
				SetData(task),
		)
	}
}

// Update implements TaskController.
func (tc *taskController) Update() echo.HandlerFunc {
	type task struct {
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		id := c.Param("id")

		// Get validator from middleware
		validate := middleware.GetValidator(c)
		if validate == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Validator not configured")
		}

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		t := task{}

		// Use middleware's validator helper function for validation
		if err := middleware.ValidateRequest(c, &t); err != nil {
			return err
		}

		createdTask, err := tc.taskService.Update(c.Request().Context(), id, map[string]any{
			"description": t.Description,
		})
		if err != nil {
			return err
		}

		return c.JSON(
			http.StatusOK,
			NewResponse().
				AddMeta("status", http.StatusCreated).
				SetData(createdTask),
		)
	}
}

// Delete implements TaskController.
func (t *taskController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		// Get validator from middleware
		validate := middleware.GetValidator(c)
		if validate == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Validator not configured")
		}

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		err = t.taskService.Delete(c.Request().Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(
			http.StatusOK,
			NewResponse().
				AddMeta("status", http.StatusOK).
				SetMessage("deleted"),
		)
	}
}

func NewTaskController(i *do.Injector) (TaskController, error) {
	return &taskController{
		taskService: do.MustInvoke[service.TaskService](i),
	}, nil
}
