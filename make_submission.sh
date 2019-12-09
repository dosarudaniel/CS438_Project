#!/bin/bash


if [ $(pwd) != "/home/daniel/go/src/github.com/2_alt_hw2/Peerster" ]
then
    cd ../../../..
    exit
fi




cd ../../../..
mkdir -p tmp_Submission/src/github.com/2_alt_hw2/
cp -r src/github.com/2_alt_hw2/Peerster tmp_Submission/src/github.com/2_alt_hw2/
# delete .git folder
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/.git

# Remove binary files
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/Peerster
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/IntegrationTests/*
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/client/client
# Remove pdf files
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/*.pdf

# Remove pdf files
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/make_submission.sh
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/_Downloads
rm -rf tmp_Submission/src/github.com/2_alt_hw2/Peerster/_SharedFiles

cd tmp_Submission
tar -czvf ~/Videos/latestDSESubmission.tar.gz src/

echo ""
echo "Complete, upload the ~/Videos/latestDSESubmission.tar.gz on https://cs438.epfl.ch/"
