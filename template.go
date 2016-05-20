package main

import (
	"bytes"
	"text/template"
)

const TEMPLATE_SCRIPT string = `
#!/bin/bash
# image = {{.Image}}

DOCKER_BIN=$(which docker)

enable_git() {
  test '{{.Git}}' == 'true'
}

msg() {
  echo $@ >&2
}

__sudo() {
  if sudo -v >/dev/null 2>&1; then
    sudo "$@"
  else
    "$@"
  fi
}

docker() {
  __sudo ${DOCKER_BIN} "$@"
}

container_names() {
  docker ps --format {{"{{.Names}}"}}
}

container_exist() {
  container_names | grep $1 >/dev/null 2>&1
}

is_git_repo() {
  git rev-parse --git-dir >/dev/null 2>&1
}

workroot() {
  if enable_git && is_git_repo; then
    git rev-parse --show-toplevel
  else
    echo $PWD
  fi
}

workspace() {
  if enable_git; then
    local root=$(git rev-parse --show-toplevel)
    local rel=$(realpath --relative-base "${root}" "$PWD")
    echo "/workspace/$rel"
  else
    echo "/workspace"
  fi
}

container_id() {
  env | sha1sum | awk '{print $1}'
}

command() {
  local image='{{.Image}}'
  local command='{{.Command}}'
  local name=docker_path_${command}_$(container_id)
  if container_exist $name; then
    : nothing to do
  else
    msg create $name
    docker run -d \
      -u "$UID:$GROUPS" \
      -e "TERM=${TERM:-xterm}" \
      -e "HOME=/opt/user" \
      -v "$HOME:/opt/user" \
      -w "$(workspace)" \
      -v "$(workroot):/workspace" \
      --name $name \
      "$image" \
      /bin/sh -c 'while true; do sleep 1000000; done'
  fi

  if test -t 0; then
    # tty
    docker exec -ti $name "$command" "$@"
  else
    # input
    docker exec -i $name "$command" "$@"
  fi
}

command $@
`

func Render(locals *Command) string {
	tmpl, err := template.New("script").Parse(TEMPLATE_SCRIPT)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, locals); err != nil {
		panic(err)
	}
	return doc.String()
}
