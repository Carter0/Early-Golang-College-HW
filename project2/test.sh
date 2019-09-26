#!/bin/bash

go run testhanger.go $@ &
sleep 5
go run router/router.go $@
