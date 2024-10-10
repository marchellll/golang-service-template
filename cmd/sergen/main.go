package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)



func main() {

	ModuleName := flag.String("ModuleName", "github.com/marchellll/something", "this is the main module of the project")
	EntityName := flag.String("EntityName", "User", "the entity name to be generated")
	EntityNamePlural := flag.String("EntityNamePlural", "users", "OTIONAL. the entity name to be generated, but in plural. when not supplied, it will be the entity + 's'")
	flag.Parse()

	data := struct {
		ModuleName  string
		EntityName  string
		EntityNameLow	string
		EntityNamePlural	string
		EntityNameLowPlural string
	}{
		ModuleName:  *ModuleName,
		EntityName:  *EntityName,
		EntityNamePlural:  *EntityNamePlural,
	}


	if data.EntityNamePlural == "" {
		data.EntityNamePlural = data.EntityName + "s"
	}
	data.EntityNameLow = firstToLower(data.EntityName)
	data.EntityNameLowPlural = firstToLower(data.EntityNamePlural)


	serviceFile, err := os.Create("internal/service/"+ strings.ToLower(data.EntityName) +".go")
	die(err)
	defer serviceFile.Close()

	err = serviceTemplate.Execute(serviceFile, data)
	die(err)

	handlerFile, err := os.Create("internal/handler/"+ strings.ToLower(data.EntityName) +".go")
	die(err)
	defer handlerFile.Close()

	err = controllerTemplate.Execute(handlerFile, data)
	die(err)


	routesFile, err := os.OpenFile("internal/app/routes.go", os.O_APPEND|os.O_WRONLY, 0644)
	die(err)
	defer routesFile.Close()


	err = routeTemplate.Execute(routesFile, data)
	die(err)


	diFile, err := os.OpenFile("internal/app/di.go", os.O_APPEND|os.O_WRONLY, 0644)
	die(err)
	defer routesFile.Close()


	err = diTemplate.Execute(diFile, data)
	die(err)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func firstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
			return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
			return s
	}
	return string(lc) + s[size:]
}

var controllerTemplate = template.Must(template.New("").Parse(`package handler

import (
	"net/http"

	"{{ .ModuleName }}/internal/dao/model"
	"{{ .ModuleName }}/internal/service"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/do"

	"github.com/labstack/echo/v4"
)

type {{ .EntityName }}Controller interface {
	Create() echo.HandlerFunc
	Find() echo.HandlerFunc
	GetById() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type {{ .EntityNameLow }}Controller struct {
	{{ .EntityNameLow }}Service service.{{ .EntityName }}Service
}

// Create implements {{ .EntityName }}Controller.
func (tc *{{ .EntityNameLow }}Controller) Create() echo.HandlerFunc {
	// TODO: validator instance & translation creation can be moved to middleware
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // 'en' should be from request header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	type {{ .EntityNameLow }} struct {
		Description string `+"`"+`json:"description" validate:"required"`+"`"+`
	}

	return func(c echo.Context) error {
		t := {{ .EntityNameLow }}{}

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

		created, err := tc.{{ .EntityNameLow }}Service.Create(c.Request().Context(), model.{{ .EntityName }}{
			// TODO: fill
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"meta": map[string]any{
				"status": http.StatusCreated,
			},
			"data": created,
		})
	}
}

// Get implements {{ .EntityName }}Controller.
func (t *{{ .EntityNameLow }}Controller) Find() echo.HandlerFunc {
	return func(c echo.Context) error {

		list, err := t.{{ .EntityNameLow }}Service.Find(c.Request().Context())

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"total":  len(list),
				"status": http.StatusOK,
			},
			"data": list,
		})
	}
}

// Get{{ .EntityName }} implements {{ .EntityName }}Controller.
func (tc *{{ .EntityNameLow }}Controller) GetById() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		entity, err := tc.{{ .EntityNameLow }}Service.Get(c.Request().Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"meta": map[string]any{
				"status": http.StatusOK,
			},
			"data": entity,
		})
	}
}

// Update implements {{ .EntityName }}Controller.
func (tc *{{ .EntityNameLow }}Controller) Update() echo.HandlerFunc {

	// TODO: validator instance & translation creation can be moved to middleware
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // 'en' should be from request header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	type {{ .EntityNameLow }} struct {
		Description string `+"`"+`json:"description" validate:"required"`+"`"+`
	}

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		t := {{ .EntityNameLow }}{}

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

		updated, err := tc.{{ .EntityNameLow }}Service.Update(c.Request().Context(), id, map[string]any{
			"description": t.Description,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"meta": map[string]any{
				"status": http.StatusCreated,
			},
			"data": updated,
		})
	}
}

// Delete implements {{ .EntityName }}Controller.
func (t *{{ .EntityNameLow }}Controller) Delete() echo.HandlerFunc {
	validate := validator.New()

	return func(c echo.Context) error {
		id := c.Param("id")

		err := validate.Var(id, "required,uuid")
		if err != nil {
			return err
		}

		err = t.{{ .EntityNameLow }}Service.Delete(c.Request().Context(), id)
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

func New{{ .EntityName }}Controller(i *do.Injector) ({{ .EntityName }}Controller, error) {
	return &{{ .EntityNameLow }}Controller{
		{{ .EntityNameLow }}Service: do.MustInvoke[service.{{ .EntityName }}Service](i),
	}, nil
}

`));


var serviceTemplate = template.Must(template.New("").Parse(`package service

import (
	"context"
	"{{ .ModuleName }}/internal/dao/model"
	"{{ .ModuleName }}/internal/dao/query"
	"{{ .ModuleName }}/internal/errz"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type {{ .EntityName }}Service interface {
	Create(ctx context.Context, {{ .EntityNameLow }} model.{{ .EntityName }}) (*model.{{ .EntityName }}, error)
	Get(ctx context.Context, id string) (*model.{{ .EntityName }}, error)
	Find(ctx context.Context) ([]*model.{{ .EntityName }}, error)
	Update(ctx context.Context, id string, entity map[string]any) (*model.{{ .EntityName }}, error)
	Delete(ctx context.Context, id string) error
}

type {{ .EntityNameLow }}Service struct {
	db    *gorm.DB
	q     *query.Query
	redis *redis.Client
}

func New{{ .EntityName }}Service(i *do.Injector) ({{ .EntityName }}Service, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &{{ .EntityNameLow }}Service{
		db:    db,
		q:     query.Use(db),
		redis: do.MustInvoke[*redis.Client](i),
	}, nil
}

// Create implements {{ .EntityName }}Service.
func (s *{{ .EntityNameLow }}Service) Create(ctx context.Context, entity model.{{ .EntityName }}) (*model.{{ .EntityName }}, error) {
	newID, err := uuid.NewV7()

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to generate new id", err)
	}

	entityp := &entity
	entityp.ID = newID.String()
	if err := query.Use(s.db).WithContext(ctx).{{ .EntityName }}.Create(entityp); err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to create {{ .EntityNameLow }}", err)
	}
	return entityp, nil
}

// Get implements {{ .EntityName }}Service.
func (s *{{ .EntityNameLow }}Service) Get(ctx context.Context, id string) (*model.{{ .EntityName }}, error) {
	entity, err := s.q.WithContext(ctx).{{ .EntityName }}.Where(s.q.{{ .EntityName }}.ID.Eq(id)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entity", err)
	}

	return entity, nil
}

// GetAll implements {{ .EntityName }}Service.
func (s *{{ .EntityNameLow }}Service) Find(ctx context.Context) ([]*model.{{ .EntityName }}, error) {
	// var entities []model.{{ .EntityName }}

	entities, err := s.q.WithContext(ctx).{{ .EntityName }}.Find()

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to get entities", err)
	}

	return entities, nil
}

// Update implements {{ .EntityName }}Service.
// using map here to avoid headache of handling Go's zero value
// we pass whatever passed validation in handler
func (s *{{ .EntityNameLow }}Service) Update(ctx context.Context, id string, entity map[string]any) (*model.{{ .EntityName }}, error) {
	_, err := s.q.WithContext(ctx).{{ .EntityName }}.Where(s.q.{{ .EntityName }}.ID.Eq(id)).Updates(entity)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errz.NewPrettyError(http.StatusNotFound, "not_found", "entity not found", err)
	}

	if err != nil {
		return nil, errz.NewPrettyError(http.StatusInternalServerError, "internal_server_error", "failed to update entities", err)
	}

	return s.Get(ctx, id)
}

// Delete implements {{ .EntityName }}Service.
func (s *{{ .EntityNameLow }}Service) Delete(ctx context.Context, id string) error {
	_, err := s.q.WithContext(ctx).{{ .EntityName }}.Where(s.q.{{ .EntityName }}.ID.Eq(id)).Delete()
	if err != nil {
		return errors.Wrap(err, "failed to delete {{ .EntityNameLow }}")
	}
	return nil
}

`))


var routeTemplate = template.Must(template.New("").Parse(`
add{{ .EntityName }}Routes(injector, e) // FIXME: move me
func add{{ .EntityName }}Routes(injector *do.Injector, e *echo.Echo) {
	group := e.Group("/{{ .EntityNameLowPlural }}")

	group.POST("", do.MustInvoke[handler.{{ .EntityName }}Controller](injector).Create())
	group.GET("", do.MustInvoke[handler.{{ .EntityName }}Controller](injector).Find())
	group.GET("/:id", do.MustInvoke[handler.{{ .EntityName }}Controller](injector).GetById())
	group.PATCH("/:id", do.MustInvoke[handler.{{ .EntityName }}Controller](injector).Update())
	group.DELETE("/:id", do.MustInvoke[handler.{{ .EntityName }}Controller](injector).Delete())
}

`));


var diTemplate = template.Must(template.New("").Parse(`
// FIXME: move me
// services
do.Provide(injector, service.New{{ .EntityName }}Service)

// handler
do.Provide(injector, handler.New{{ .EntityName }}Controller)

`));
