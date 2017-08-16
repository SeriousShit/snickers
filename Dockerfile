FROM serious/snickers-docker:latest

# Download snickers
RUN go get github.com/SeriousShit/snickers

# Run snickers!
RUN curl -O http://flv.io/snickers/config.json

RUN cd GOPATH/src/github.com/SeriousShit/snickers && make fetch && make run

#RUN go install github.com/SeriousShit/snickers
ENTRYPOINT snickers
EXPOSE 8000
