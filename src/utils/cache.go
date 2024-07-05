package utils

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = cache.New(30*time.Minute, 10*time.Minute)