FROM golang:1.8

WORKDIR /go/src/bitbucket-eng-sjc1.cisco.com/an/GoHospital/
COPY . .
RUN ["go", "get", "./..."]
RUN ["go", "build"]
EXPOSE 8080
CMD ["./GoHospital"]
