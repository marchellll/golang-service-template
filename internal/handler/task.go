package handler

import (
	"net/http"

	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/service"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/do"

	"github.com/labstack/echo/v4"
)

type TaskController interface {
	Create() echo.HandlerFunc
	Find() echo.HandlerFunc
	GetById() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type taskController struct {
	taskService service.TaskService
}

// Create implements TaskController.
func (tc *taskController) Create() echo.HandlerFunc {
	// TODO: validator instance & translation creation can be moved to middleware
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // `en` should be from request header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	type task struct {
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		t := task{}

		if err := c.Bind(&t); err != nil {
			// TODO: this specific type in middleware
			// to show 400 for bind error
			return err
		}

		err := validate.Struct(t)
		if err != nil {
			// TODO: this specific type in middleware
			// to show 400 for validation errors
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

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"total":  len(tasks),
				"status": http.StatusOK,
			},
			"data": tasks,
		})
	}
}

// GetTask implements TaskController.
func (tc *taskController) GetById() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		task, err := tc.taskService.Get(c.Request().Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"status": http.StatusOK,
			},
			"data": task,
		})
	}
}

// Update implements TaskController.
func (tc *taskController) Update() echo.HandlerFunc {

	// TODO: validator instance & translation creation can be moved to middleware
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // `en` should be from request header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	type task struct {
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		t := task{}

		if err := c.Bind(&t); err != nil {
			// TODO: this specific type in middleware
			// to show 400 for bind error
			return err
		}

		err = validate.Struct(t)
		if err != nil {
			// TODO: this specific type in middleware
			// to show 400 for validation errors
			return err
		}

		createdTask, err := tc.taskService.Update(c.Request().Context(), id, map[string]any{
			"description": t.Description,
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

// Delete implements TaskController.
func (t *taskController) Delete() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		err = t.taskService.Delete(c.Request().Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"status":  http.StatusOK,
				"message": "deleted",
			},
		})
	}
}

func NewTaskController(i *do.Injector) (TaskController, error) {
	return &taskController{
		taskService: do.MustInvoke[service.TaskService](i),
	}, nil
}
