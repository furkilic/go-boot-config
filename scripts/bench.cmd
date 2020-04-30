@echo off
title %0
set BASE_DIR=%~dp0..

"%BASE_DIR%\gow.cmd" test -bench=. "%BASE_DIR%\..."