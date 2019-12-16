#!/bin/bash


if [ $(pwd) != "/home/daniel/go/src/github.com/dosarudaniel/CS438_Project/" ]
then
    cd ../../../..
    exit
fi




cd ../../../..
mkdir -p tmp_Submission/src/github.com/dosarudaniel/
cp -r src/github.com/dosarudaniel/CS438_Project tmp_Submission/src/github.com/dosarudaniel/
# delete .git folder
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/.git

# Remove binary files
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/CS438_Project
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/IntegrationTests/*
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/client/client
# Remove pdf files
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/*.pdf

# Remove pdf files
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/make_submission.sh
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/_Downloads
rm -rf tmp_Submission/src/github.com/dosarudaniel/CS438_Project/_SharedFiles

cd tmp_Submission
tar -czvf ~/Videos/latestDSESubmission.tar.gz src/

echo ""
echo "Complete, upload the ~/Videos/latestDSESubmission.tar.gz on https://cs438.epfl.ch/"
