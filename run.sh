#!/usr/bin/env bash
set -xe
cd src;
go build;
cd ..;
./src/playma;
