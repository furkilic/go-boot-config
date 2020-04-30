@echo off
title %0
set BASE_DIR=%~dp0..

docker-compose -f "%BASE_DIR%\deployments\docker-compose.yml" up