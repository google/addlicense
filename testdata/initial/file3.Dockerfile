# escape=`
FROM microsoft/nanoserver
COPY testfile.txt c:\
RUN dir c:\
