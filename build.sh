#!/bin/sh

# Copyright 2018 NTT Group

# Permission is hereby granted, free of charge, to any person obtaining a copy of this 
# software and associated documentation files (the "Software"), to deal in the Software 
# without restriction, including without limitation the rights to use, copy, modify, 
# merge, publish, distribute, sublicense, and/or sell copies of the Software, and to 
# permit persons to whom the Software is furnished to do so, subject to the following 
# conditions:

# The above copyright notice and this permission notice shall be included in all copies 
# or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR 
# PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
# FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR 
# OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
# DEALINGS IN THE SOFTWARE.

if [ x = x$DOCKER_USER ]; then
	echo "The variable $DOCKER_USER is not set. Try: export DOCKER_USER=<your-docker-id>" 
	exit 1
fi

# Build Frontend
if [ xfrontend = x$1 ]; then
    cd html
    grunt
    if [ $? -ne 0 ]; then
        echo "ERROR building fronted"
        exit 1
    fi
    cd .. 
fi

# Build Backend
docker build -t $DOCKER_USER/edge-anki-adas .
if [ $? -ne 0 ]; then
    echo "ERROR building backend"
    exit 1
fi

# Push image
docker push $DOCKER_USER/edge-anki-adas
if [ $? -ne 0 ]; then
    echo "ERROR pushing image"
    exit 1
fi
