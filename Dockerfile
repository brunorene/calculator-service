ARG UID=10001

FROM golang:1.18 AS builder
WORKDIR /app/
COPY --chown=${UID} go.mod go.sum Makefile ./
RUN make setup
COPY --chown=${UID} . ./
RUN make build

FROM scratch AS output
COPY --from=builder /app/build/calculator-service /
ENTRYPOINT ["/calculator-service"]
