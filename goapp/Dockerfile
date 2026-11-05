# npm install
FROM node:18.6.0 as npm-build
RUN mkdir /ng-app
WORKDIR /ng-app
RUN mkdir /node_modules
COPY ./artemis/package.json .
COPY ./artemis/package-lock.json .
RUN npm install

# angular app
FROM node:18.6.0 as ng-build
ARG ngBuildCommand=build-local
RUN mkdir /ng-app
COPY ./artemis /ng-app
WORKDIR /ng-app
COPY --from=npm-build /ng-app/node_modules node_modules
ENV NODE_OPTIONS=--max-old-space-size=4096

RUN npm run ${ngBuildCommand}

# angular ares app
FROM node:18.6.0 as ares-build
ARG aresBuildCommand=ares:build
RUN mkdir /app
WORKDIR /app
RUN mkdir /app/ng-app
RUN mkdir /app/bin
RUN mkdir /app/bin/static
RUN mkdir /app/bin/static/css
COPY ./artemis /ng-app
COPY ./bin/static/css /bin/static/css
WORKDIR /ng-app
COPY --from=npm-build /ng-app/node_modules node_modules

RUN npm run ${aresBuildCommand}

# angular login app
FROM node:18.6.0 as login-build
ARG loginBuildCommand=login:build
COPY ./artemis /ng-app
COPY ./bin/static/css /bin/static/css
WORKDIR /ng-app
COPY --from=npm-build /ng-app/node_modules node_modules

RUN npm run ${loginBuildCommand}

# elements build
FROM node:18.6.0 as elements-build
ARG elementsBuildCommand=elements:docker:build-local
RUN mkdir /ng-app
WORKDIR /ng-app
COPY ./artemis/angular.json .
COPY ./artemis/tsconfig.app.json .
COPY ./artemis/tsconfig.spec.json .
COPY ./artemis/tsconfig.json .
COPY ./artemis/scripts scripts
COPY ./artemis/package-lock.json .
COPY ./artemis/package.json .
COPY --from=npm-build /ng-app/node_modules node_modules
COPY ./artemis/projects projects
COPY ./artemis/src src
RUN npm run ${elementsBuildCommand}

# Second-stage
FROM golang:alpine as builder

RUN apk update
RUN apk add --no-cache bash tzdata
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app
ENV GOBIN /build/bin

RUN mkdir /build
COPY . /build
WORKDIR /build
RUN rm -rf artemis

RUN mkdir log
RUN mkdir uploads
RUN mkdir /build/bin/exports

WORKDIR hermes
RUN go install hermes.go

# Final image
FROM alpine
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=ng-build /bin/static/artemis /app/bin/static/artemis
COPY --from=ares-build /bin/static/ares /app/bin/static/ares
COPY --from=login-build /bin/static/login /app/bin/static/login
COPY --from=elements-build /ng-app/dist/js /app/bin/static/js
COPY --from=elements-build /ng-app/dist/css /app/bin/static/css

RUN chown -R app:app /app
WORKDIR /app/bin
USER app
CMD ["./hermes"]
