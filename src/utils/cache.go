package utils

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = cache.New(80*time.Minute, 10*time.Minute)
