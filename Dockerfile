FROM golang:alpine as builder
RUN apk add git --no-cache tzdata openjdk11 unzip curl

WORKDIR /app
COPY . ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init && GODEBUG=multipartmaxheaders=100000,multipartmaxparts=100000 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o allure_server

FROM alpine:latest as runner
ARG KEEP_RESULTS_HISTORY=true
ARG KEEP_HISTORY_LATEST=1
ARG ARCH=amd64
ARG SOURCE_DIR=/app
ARG ALLURE_RELEASE=2.24.1
ARG ALLURE_REPO=https://repo.maven.apache.org/maven2/io/qameta/allure/allure-commandline
ARG UID=1002
ARG GID=1002
ARG HOST=0.0.0.0
ARG PORT=5050
ARG THREADS=4

ENV GIN_MODE release
ENV BASE_PATH /allure-docker-service
ENV KEEP_RESULTS_HISTORY ${KEEP_RESULTS_HISTORY}
ENV KEEP_HISTORY_LATEST ${KEEP_HISTORY_LATEST}
ENV APP_DATA_DIR=/AppData
ENV SOURCE_DIR $SOURCE_DIR
ENV ALLURE_HOME=/allure-$ALLURE_RELEASE
ENV ALLURE_HOME_SL=/allure
ENV PATH=$PATH:$ALLURE_HOME/bin
ENV ALLURE_RESOURCES=$APP_DATA_DIR/resources
ENV RESULTS_DIRECTORY=$APP_DATA_DIR/allure-results
ENV REPORT_DIRECTORY=$APP_DATA_DIR/allure-report
ENV RESULTS_HISTORY=$RESULTS_DIRECTORY/history
ENV REPORT_HISTORY=$REPORT_DIRECTORY/history
ENV ALLURE_VERSION=$APP_DATA_DIR/version
ENV EMAILABLE_REPORT_FILE_NAME='emailable-report-allure-docker-service.html'
ENV STATIC_CONTENT_PROJECTS=$APP_DATA_DIR/projects
ENV DEFAULT_PROJECT=default
ENV DEFAULT_PROJECT_ROOT=$STATIC_CONTENT_PROJECTS/$DEFAULT_PROJECT
ENV DEFAULT_PROJECT_RESULTS=$DEFAULT_PROJECT_ROOT/results
ENV DEFAULT_PROJECT_REPORTS=$DEFAULT_PROJECT_ROOT/reports
ENV EXECUTOR_FILENAME=executor.json
ENV DOWNLOAD_REPORT_CSV_DESTINATION_PATH=${ROOT}/reports/csv
ENV CHECK_RESULTS_EVERY_SECONDS=NONE

COPY --from=builder /app/allure_server ./allure_server
RUN chmod u+x ./${name}

RUN apk add git --no-cache tzdata openjdk11 unzip curl
RUN curl ${ALLURE_REPO}/${ALLURE_RELEASE}/allure-commandline-${ALLURE_RELEASE}.zip -L -o /tmp/allure-commandline.zip && \
    unzip -q /tmp/allure-commandline.zip -d / && \
    apk del unzip curl --purge && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    chmod -R +x /allure-$ALLURE_RELEASE/bin && \
    mkdir -p $APP_DATA_DIR

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/allure" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "allure"

RUN echo -n $(allure --version) > ${ALLURE_VERSION} && \
    echo "ALLURE_VERSION: "$(cat ${ALLURE_VERSION}) && \
    mkdir $ALLURE_HOME_SL && ln -s $ALLURE_HOME/* $ALLURE_HOME_SL
    # ln -s $STATIC_CONTENT_PROJECTS $APP_DATA_DIR/projects && \
    # ln -s $DEFAULT_PROJECT_REPORTS $APP_DATA_DIR/default-reports

RUN chown -R allure:allure $APP_DATA_DIR

VOLUME [ "$RESULTS_DIRECTORY" ]

EXPOSE $PORT

USER allure

CMD ["/bin/sh", "-c", "GODEBUG=multipartmaxheaders=100000,multipartmaxparts=100000 ./allure_server"]