FROM golang:1.18-buster AS build

WORKDIR /build

ENV PORT=8080
ENV MODE=production

COPY . ./

RUN make build

FROM scratch

WORKDIR /app

COPY --from=build /build/build/cdn /app/cdn

CMD ["/app/cdn"]