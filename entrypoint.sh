#!/bin/sh -l

/opt/guac/guacone collect files $1 --gql-addr $2 >> $GITHUB_OUTPUT