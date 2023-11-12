FROM golang:1.21.3-alpine3.18 AS BUILD_IMAGE

WORKDIR /opt/fizzy-build/

COPY . .

RUN go build

FROM alpine:3.18

WORKDIR /opt/fizzy/

COPY --from=BUILD_IMAGE /opt/fizzy-build/fizzy .
COPY --from=BUILD_IMAGE /opt/fizzy-build/views/ ./views/
COPY --from=BUILD_IMAGE /opt/fizzy-build/static/ ./static/

ENV LISTEN_ADDRESS=0.0.0.0:3000

EXPOSE 3000

CMD [ "./fizzy" ]