FROM golang:1.8

WORKDIR /go/src/bitbucket-eng-sjc1.cisco.com/an/GoHospital/
COPY . .
RUN ["go", "get", "-d", "-v", "-u", "github.com/mongodb/mongo-go-driver/..."]
#RUN ["dep", "ensure", "-add", "github.com/mongodb/mongo-go-driver/mongo"]
RUN ["go", "build"]
EXPOSE 8080
CMD ["./GoHospital"]
