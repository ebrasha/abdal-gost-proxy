@echo off
chcp 65001 >nul
setlocal
cd /d "%~dp0"

echo === Updating dependencies (go mod tidy) ===
go mod tidy
if errorlevel 1 (
    echo If you get 403/network errors, try: set GOPROXY=https://proxy.golang.org,direct
    goto :err
)

echo Creating dist folders...
if not exist "dist\windows" mkdir "dist\windows"
if not exist "dist\linux"   mkdir "dist\linux"

echo.
echo === Building Windows (server, client, reality-keygen) ===
go build -o dist\windows\abdal-gost-proxy-server.exe main.go
if errorlevel 1 goto :err
go build -o dist\windows\abdal-gost-proxy-client.exe client_main.go
if errorlevel 1 goto :err
go build -o dist\windows\reality-keygen.exe .\tools\reality-keygen
if errorlevel 1 goto :err

echo.
echo === Building Linux (server, client, reality-keygen) ===
set GOOS=linux
set GOARCH=amd64
go build -o dist\linux\abdal-gost-proxy-server main.go
if errorlevel 1 goto :err
go build -o dist\linux\abdal-gost-proxy-client client_main.go
if errorlevel 1 goto :err
go build -o dist\linux\reality-keygen .\tools\reality-keygen
if errorlevel 1 goto :err
set GOOS=
set GOARCH=

echo.
echo === Copying config files ===
copy /Y "abdal-gost-proxy-server.json" "dist\windows\" >nul
copy /Y "abdal-gost-proxy-client.json" "dist\windows\" >nul
copy /Y "abdal-gost-proxy-server.json" "dist\linux\"   >nul
copy /Y "abdal-gost-proxy-client.json" "dist\linux\"   >nul

echo.
echo Done. Output:
echo   dist\windows\  - Server, Client, reality-keygen (.exe) + configs
echo   dist\linux\    - Server, Client, reality-keygen + configs
goto :eof

:err
echo Build failed.
exit /b 1
