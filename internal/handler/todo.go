package handler

import (
	"net/http"
	"strconv"

	"golang-service-template/internal/dao/model"
	"golang-service-template/internal/service"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/do"

	"github.com/labstack/echo/v4"
)

type TodoController interface {
	GetTodos() echo.HandlerFunc
	GetTodo() echo.HandlerFunc
	CreateTodo() echo.HandlerFunc
	DeleteTodo() echo.HandlerFunc
}

type todoController struct {
	todoService service.TodoService
}

// GetTodo implements TodoController.
func (tc *todoController) GetTodo() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,number")
		if err != nil {
			return err
		}

		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}

		todo, err := tc.todoService.Get(c.Request().Context(), idUint)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"status":  http.StatusOK,
			},
			"data": todo,
		})
	}
}

// CreateTodo implements TodoController.
func (tc *todoController) CreateTodo() echo.HandlerFunc {
	// TODO: validator instance & translation creation can be moved to middleware
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // `en` should be from request header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	type todo struct {
		Text string `json:"text" validate:"required"`
	}

	return func(c echo.Context) error {
		t := todo{}

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

		createdTodo, err := tc.todoService.Create(c.Request().Context(), model.Todo{
			Text: t.Text,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"meta": map[string]any{
				"status": http.StatusCreated,
			},
			"data": createdTodo,
		})
	}
}

// DeleteTodo implements TodoController.
func (t *todoController) DeleteTodo() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,number")
		if err != nil {
			return err
		}

		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}

		err = t.todoService.Delete(c.Request().Context(), idUint)
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

// GetTodos implements TodoController.
func (t *todoController) GetTodos() echo.HandlerFunc {
	return func(c echo.Context) error {

		todos, err := t.todoService.GetAll(c.Request().Context())

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"total":  len(todos),
				"status": http.StatusOK,
			},
			"data": todos,
		})
	}
}

func NewTodoController(i *do.Injector) (TodoController, error) {
	return &todoController{
		todoService: do.MustInvoke[service.TodoService](i),
	}, nil
}
