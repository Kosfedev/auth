package model

import "github.com/pkg/errors"

// TODO: вопрос - почему на сервис слое? ошибка же для репозитория, не?
var ErrorNoteNotFound = errors.New("user not found")
