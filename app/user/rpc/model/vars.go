package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
var UserTypeLocal = "local"
var UserTypeGoogle = "google"
