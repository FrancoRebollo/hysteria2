@ECHO OFF

SET PATH=C:\oracle\64bits\instantclient_19_22;%PATH%

call cls

rmdir /S /Q logs

REM Obteniendo la fecha y hora actual en formato JSON
for /f "tokens=2,3,4,5,6 usebackq delims=:/ " %%a in ('%date% %time%') do set mydate={"FECHA":"%%a/%%b/%%c", "HORA":"%%d:%%e"}


REM Instalando dependencias usando Go Modules
go mod tidy

REM Compilando y ejecutando la aplicación Go
go run cmd/apis/main.go

pause