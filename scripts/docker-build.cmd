@echo off
title %0
set BASE_DIR=%~dp0..

set GOOS=linux
call "%BASE_DIR%\scripts\build.cmd"

docker build  -t go-boot-config:0.0.1 "%BASE_DIR%" -f "%BASE_DIR%\build\package\Dockerfile"
