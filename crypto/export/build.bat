@echo off
REM set PATH=D:\User\dev\bin\msys64\mingw32\bin;%PATH%
REM set GOARCH=386
REM set CGO_ENABLED=1
go build -buildmode=c-shared -o LDSCrypto.dll main.go
