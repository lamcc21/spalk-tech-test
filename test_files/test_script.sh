#!/bin/bash

# Test success case
cat ./test_success.ts | ../mpegts-parser
if [ $? -eq 0 ]; then
    echo "Success test passed"
else
    echo "Success test failed"
fi

# Random video success case
cat ./sample_1280x720.ts | ../mpegts-parser
if [ $? -eq 0 ]; then
    echo "Random test passed"
else
    echo "Random test failed"
fi

# Test failure case
cat ./test_failure.ts | ../mpegts-parser
if [ $? -eq 0 ]; then
    echo "Failure test failed"
else
    echo "Failure test passed"
fi