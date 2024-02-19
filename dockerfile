FROM golang:1.21.3-alpine3.18 AS build-stage

WORKDIR /home/abhishek/Clean/

COPY ./ /home/abhishek/Clean/

RUN mkdir -p /home/abhishek/build
RUN go mod download
RUN go build -v -o /home/abhishek/build/api ./

FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /home/abhishek/build/api /api
COPY --from=build-stage /home/abhishek/Clean/template /template
COPY --from=build-stage /home/abhishek/Clean/.env /

EXPOSE 8080
CMD ["/api", "-host", "0.0.0.0", "-port", "8080"]


