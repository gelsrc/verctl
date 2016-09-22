@echo off

setlocal

set /P release=<version.txt

echo Сборка приложения для Windows

set GOOS=windows
set GOARCH=386

go build -ldflags "-X main.Version=%release%" -o verctl-%release%.exe verctl.go 

if errorlevel 1 (
    echo Ошибка при сборке приложения
    exit /b 1
)

echo Сборка приложения для GNU/Linux

set GOOS=linux
set GOARCH=amd64

go build -ldflags "-X main.Version=%release%" -o verctl-%release% verctl.go

if errorlevel 1 (
    echo Ошибка при сборке приложения
    exit /b 1
)

echo Сборка приложения успешно завершена