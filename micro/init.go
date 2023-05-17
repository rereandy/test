package main

import (
	"github.com/rereandy/micro"
)

var Provider = mirco.NewProvider(
	&MySvcImpl{},
	&Other{},
)
