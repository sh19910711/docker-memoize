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

config_env() {
{{ range $env := .Env }}
  if [[ "${{$env}}" =~ "=" ]]; then
    echo "-e {{$env}}"
  else
    echo "-e {{$env}}=\"${{$env}}\""
  fi
{{ end }}
}

config_exec_env() {
{{ range $env := .ExecEnv }}
  echo '{{$env}}'
{{ end }}
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

docker_run() {
  local name=$1
  local image=$2
  local cmd=()

  cmd+=(docker run -d)
  cmd+=($(config_env))
  cmd+=(-u "$UID:$GROUPS")
  cmd+=(-e "TERM=${TERM:-xterm}")
  cmd+=(-e "HOME=$HOME")
  cmd+=(-v "$HOME:$HOME")
  cmd+=(-w "$(workspace)")
  cmd+=(-v "$(workroot):/workspace")
  cmd+=(--name "$name")
  cmd+=("$image")

  echo ${cmd[@]} /bin/sh -c 'while true; do sleep 1000000; done'
  ${cmd[@]} /bin/sh -c 'while true; do sleep 1000000; done'
}

command() {
  local image='{{.Image}}'
  local command='{{.Command}}'
  local name=docker_path_${command}_$(container_id)
  if container_exist $name; then
    : nothing to do
  else
    msg create $name
    docker_run $name $image
  fi

  if test -t 0; then
    # tty
    local args=$@
    docker exec -ti $name bash -c "$(config_exec_env) $command $args"
  else
    # input
    local args=$@
    docker exec -i $name bash -c "$(config_exec_env) $command $args"
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
