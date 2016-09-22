@echo off

setlocal

set /P release=<version.txt

echo ���ઠ �ਫ������ ��� Windows

set GOOS=windows
set GOARCH=386

go build -ldflags "-X main.Version=%release%" -o verctl-%release%.exe verctl.go 

if errorlevel 1 (
    echo �訡�� �� ᡮથ �ਫ������
    exit /b 1
)

echo ���ઠ �ਫ������ ��� GNU/Linux

set GOOS=linux
set GOARCH=amd64

go build -ldflags "-X main.Version=%release%" -o verctl-%release% verctl.go

if errorlevel 1 (
    echo �訡�� �� ᡮથ �ਫ������
    exit /b 1
)

echo ���ઠ �ਫ������ �ᯥ譮 �����襭�