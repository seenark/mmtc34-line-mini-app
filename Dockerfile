# syntax=docker/dockerfile:1
# Stage 1: TS types from Go structs (tygo needs Go; sqlc output is committed, so no sqlc here).
FROM golang:1.27 AS types
WORKDIR /src/apps/server
RUN go install github.com/gzuidhof/tygo@v0.2.21
COPY apps/server/go.mod apps/server/go.sum ./
RUN go mod download
COPY apps/server/ ./
RUN mkdir -p ../web/src/api && /go/bin/tygo generate

# Stage 2: SPA. BASE_PATH is baked into asset URLs; the server refuses to start on a mismatch.
FROM oven/bun:1.4.4 AS web
ARG BASE_PATH=/mmtc34
ENV BASE_PATH=$BASE_PATH
WORKDIR /src
COPY package.json bun.lock ./
COPY apps/web apps/web
COPY packages packages
RUN bun install --frozen-lockfile
COPY --from=types /src/apps/web/src/api/types.gen.ts apps/web/src/api/types.gen.ts
WORKDIR /src/apps/web
RUN bun run build # writes ../server/internal/spa/dist

# Stage 3: static binary with the SPA embedded.
FROM types AS build
COPY --from=web /src/apps/server/internal/spa/dist internal/spa/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mmtc34 . && mkdir /out/data

# Stage 4: runtime.
FROM gcr.io/distroless/static:nonroot
ARG BASE_PATH=/mmtc34
ENV BASE_PATH=$BASE_PATH PORT=8080 DATA_DIR=/data
COPY --from=build /out/mmtc34 /mmtc34
COPY --from=build --chown=65532:65532 /out/data /data
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["/mmtc34"]
