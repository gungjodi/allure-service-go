ARG APP_VERSION="0.0.1-$(date)"
ARG GODEBUG=multipartmaxheaders=100000,multipartmaxparts=100000

FROM golang:alpine as builder
RUN apk add git --no-cache tzdata openjdk11 unzip curl

WORKDIR /app
ARG GODEBUG
ARG APP_VERSION
COPY . ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV GODEBUG ${GODEBUG}
ENV APP_VERSION ${APP_VERSION}

RUN swag init && GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o allure_server

FROM alpine:latest as runner
WORKDIR /app
ARG APP_VERSION
ARG GODEBUG
ARG KEEP_RESULTS_HISTORY=true
ARG KEEP_HISTORY_LATEST=1
ARG ALLURE_RELEASE=2.24.1
ARG ALLURE_REPO=https://repo.maven.apache.org/maven2/io/qameta/allure/allure-commandline
ARG UID=1002
ARG GID=1002
ARG HOST=0.0.0.0
ARG PORT=5050

ENV APP_VERSION ${APP_VERSION}
ENV APP_MODE release
ENV BASE_PATH /allure-service
ENV ALLURE_HOME=/allure
ENV APP_DATA_DIR=/app/AppData
ENV BACKUP_DATA_DIR=/app/BackupData
ENV PATH=$PATH:$ALLURE_HOME/bin
ENV ALLURE_RESOURCES=/app/resources
ENV ALLURE_VERSION=${APP_DATA_DIR}/version

ENV KEEP_RESULTS_HISTORY ${KEEP_RESULTS_HISTORY}
ENV KEEP_HISTORY_LATEST ${KEEP_HISTORY_LATEST}
ENV DOWNLOAD_REPORT_CSV_DESTINATION_PATH ${APP_DATA_DIR}/reports/csv

COPY --from=builder /app/allure_server /app/allure_server
RUN chmod +x ./allure_server

RUN apk add git --no-cache tzdata openjdk11 unzip curl doas
RUN curl ${ALLURE_REPO}/${ALLURE_RELEASE}/allure-commandline-${ALLURE_RELEASE}.zip -L -o /tmp/allure-commandline.zip && \
    unzip -q /tmp/allure-commandline.zip -d / && \
    apk del unzip curl --purge && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    mv /allure-${ALLURE_RELEASE} ${ALLURE_HOME} && \
    chmod -R +x ${ALLURE_HOME}/bin

ENV USER_ID ${UID}
ENV GROUP_ID ${GID}
ENV USER_NAME=allure
ENV GROUP_NAME=allure

RUN adduser $USER_NAME -G wheel --disabled-password --no-create-home \
    --home ${ALLURE_HOME} --uid ${USER_ID} && \
    echo 'permit nopass :wheel as root' >> /etc/doas.d/doas.conf

RUN mkdir -p ${APP_DATA_DIR} && \
    echo -n $(allure --version) > ${ALLURE_VERSION} && \
    echo "ALLURE_VERSION: "$(cat ${ALLURE_VERSION}) && \
    mkdir allure-results && allure generate -c -o /tmp/resources && \
    mkdir ${ALLURE_RESOURCES} && \
    mkdir -p $DOWNLOAD_REPORT_CSV_DESTINATION_PATH && \
    chmod -R 777 $DOWNLOAD_REPORT_CSV_DESTINATION_PATH && \
    cp /tmp/resources/app.js ${ALLURE_RESOURCES} && \
    cp /tmp/resources/styles.css ${ALLURE_RESOURCES} && \
    rm -rf /tmp/resources

RUN chown -R allure:wheel ${APP_DATA_DIR}

ENV HOST ${HOST}
ENV PORT ${PORT}
ENV GODEBUG ${GODEBUG}

EXPOSE ${PORT}

ENTRYPOINT [ "/app/allure_server" ]