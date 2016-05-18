package main_test

import (
	"."
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderImage(t *testing.T) {
	locals := main.Command{Command: "my-command", Image: "my-image"}
	script := main.Render(&locals)
	assert.Regexp(t, "my-image", script, "should contain image name")
}

func TestRenderCommand(t *testing.T) {
	locals := main.Command{Command: "my-command", Image: "my-image"}
	script := main.Render(&locals)
	assert.Regexp(t, "my-command", script, "should contain command")
}
