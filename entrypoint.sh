#!/bin/sh -l
# Copyright 2021 Nordcloud Oy or its affiliates. All Rights Reserved.

pwd=$PWD
echo $TEMPLATE > LICENSE.tmpl
cd /github/workspace
$pwd/addlicense -c "$HOLDER" -f $pwd/LICENSE.tmpl -u $PATTERN
