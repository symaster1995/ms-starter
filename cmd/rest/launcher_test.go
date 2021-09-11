package main_test

import (
	"context"
	"github.com/symaster1995/ms-starter/cmd/rest"
	"testing"
)

var ctx = context.Background()

func TestLauncher_Setup(t *testing.T) {
	l := main.NewTestLauncher()
	l.RunOrFail(t,ctx)

	defer l.ShutdownOrFail(t)
}