package service

import (
	"context"
	"sync"
)

var (
	ctx  = context.Background()
	once sync.Once
)
