#!/bin/sh
go build -o goose cmd/goose/*.go
mv goose ../../bin/
