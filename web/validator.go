package web

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

const errMsgValidationFailed = "Request body validation failed"

type Validator interface {
	Validate(request any) error
}

type DefaultValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewDefaultValidator() *DefaultValidator {
	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	err := en_translations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		slog.Error("Failed to register default translations for validator", slog.Any("err", err))
	}

	return &DefaultValidator{
		validate:   validate,
		translator: translator,
	}
}

func (v *DefaultValidator) Validate(request any) error {
	err := v.validate.Struct(request)
	if err == nil {
		return nil
	}

	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		details := make(map[string]string, len(validationErrors))

		for _, fieldError := range validationErrors {
			_, path, _ := strings.Cut(fieldError.Namespace(), ".")
			details[path] = fieldError.Translate(v.translator)
		}

		return NewHTTPErrorWithDetails(
			http.StatusUnprocessableEntity,
			errMsgValidationFailed,
			err,
			details,
		)
	}

	return err
}

func (v *DefaultValidator) Engine() *validator.Validate {
	return v.validate
}
