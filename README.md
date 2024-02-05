# Background
This project is based on [allure-docker-service](https://github.com/fescobar/allure-docker-service)

Ported and rewritten in Go, using [GoFiber](https://gofiber.io/) with some custom modification.

### Before Running
Set max multipart to set max files that can be uploaded
> ```export GODEBUG=multipartmaxheaders=<values>,multipartmaxparts=<value>```

