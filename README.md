# Background
This project is based on [allure-docker-service](https://github.com/fescobar/allure-docker-service)

Ported and rewritten in Go, using [GoFiber](https://gofiber.io/) with some custom modification.

# Main Features
1. Send allure result files then generate the report
2. Save reports based on project
3. Serve allure report
> For more information, visit [allure-docker-service](https://github.com/fescobar/allure-docker-service)

### Before Running
Set max multipart to set max files that can be uploaded
> ```export GODEBUG=multipartmaxheaders=<values>,multipartmaxparts=<value>```

