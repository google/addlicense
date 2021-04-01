#!/bin/sh -l
# Copyright 2021 Nordcloud Oy or its affiliates. All Rights Reserved.

cd /github/workspace
echo $TEMPLATE > LICENSE.tmpl
/go/src/app/addlicense -c "$HOLDER" -f LICENSE.tmpl -u $PATTERN
rm LICENSE.tmpl