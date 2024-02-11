# Background
This project is based on [allure-docker-service](https://github.com/fescobar/allure-docker-service)

Ported and rewritten in Go, using [GoFiber](https://gofiber.io/) with some custom modification.



## Main Features
1. Send allure result files then generate the report
2. Save reports based on project
3. Serve allure report

> For more information, visit [allure-docker-service](https://github.com/fescobar/allure-docker-service)

> [Live Demo](https://apps.gungj.tech/allure-service/swagger/index.html) (Using my home server, it might be down temporarily)

## Before Running
1. Set max multipart to set max files that can be uploaded

    ```export GODEBUG=multipartmaxheaders=<values>,multipartmaxparts=<value>```

## Pre-requisites
1. Java must be installed and ```JAVA_HOME``` is already set
2. download and install Allure based on current OS from [here](https://github.com/allure-framework/allure) then modify ```LOCAL_ALLURE_EXECUTABLE``` in .env to point to Allure executable path 
    
    - run ```which allure``` to get current allure path if installing using homebrew(MacOS)/scoop(Windows)

    - for Linux, it is recommended to download the zip then extract it to current workspace then set ```LOCAL_ALLURE_EXECUTABLE```, example command

    > ```export ALLURE_VERSION=<version>```
    > 
    > ```wget https://github.com/allure-framework/allure2/releases/download/${ALLURE_VERSION}/allure-${ALLURE_VERSION}.zip```
    > 
    > ```unzip -q allure-<version>.zip```
    >
    > ```chmod +x allure-${ALLURE_VERSION}/bin/allure```

    - Verify allure can be executed
    > ```allure-${ALLURE_VERSION}/bin/allure --help```
    
    - set ```LOCAL_ALLURE_EXECUTABLE``` in env to allure executable path

## Running in local
1. go mod download
2. go get github.com/swaggo/swag/cmd/swag
3. run server and generate swagger using command
    
    ` swag init && go run . `

4. run command below if error `"swag: command not found"` occured

    `export PATH=$(go env GOPATH)/bin:$PATH`

5. Access swager in browser by opening this URL http://localhost:5050/allure-service/swagger/index.html

## Build Docker Image

### Build docker image only

    > ```docker compose build```

### Build and deploy/run container

    > ```docker compose up -d```