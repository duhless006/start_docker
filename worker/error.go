package worker

import "errors"

var ErrWorkerNotFound = errors.New("user not found")
var ErrWorkerAlreadyExists = errors.New("user already exists")
