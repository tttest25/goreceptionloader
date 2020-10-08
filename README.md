# #goreceptionloader
![Much Win!](https://www.pauladamsmith.com/images/gopher.png)

Scrapper for paralell load data to mysql/mariadb  (similar for ETL), prepare data for reports 

## Describe
 Project for start in docker stack for accumulate data from SMD service to reception DB 
auto rotate logging (log.txt), null dependencies + multistage docker image <15 mb



## Example how to prepare and start 

### Start on host
```command
go build
DSN="user:password@tcp(10.59.0.75:3306)/modx_reception" ./goreceptionloader
source ./.env   && export DSN && ./goreceptionloader
```

### Docker

```command
docker build -t melnikov-ea/goreceptionloader:v0.1.0 .
docker run --env-file .env melnikov-ea/goreceptionloader
docker run -it --entrypoint /bin/ash  --env-file .env melnikov-ea/goreceptionloader 
```


## Comment
```json
[
{
    "requestId": "392",
    "number": "А26-05-84857211-СО1",
    "dt_modified": "2020-08-07 12:30:01.995389",
    "departmentId": "f1ae1eef-16ea-44cb-b77f-6b978ee4075d",
    "departmentName": "Администрация города Перми",
    "format": "Other",
    "formatName": "Другое",
    "isDirect": false,
    "createDate": "2020-08-03",
    "name": "Зайцева К.А.",
    "address": "Российская Федерация, Пермский край, г. Пермь",
    "email": "krisival.zaitcevu@gmail.com",
    "receiveDate": "2020-08-07 12:30:01.995389",
    "dispatchDate": null,
    "uploadDate": null,
    "exceptionMessage": null,
    "questions": [
        {
            "code": "0005.0005.0055.1122",
            "status": "NotRegistered",
            "questionStatusName": "Не зарегистрировано",
            "incomingDate": "2020-08-04"
        }
    ]
}
]
```


<b><code>goreceptionloader</code></b> remember

<table>
    <tbody>
        <tr>
            <th align="left">tttest25</th>
            <td><a href="https://github.com/tttest25">GitHub/tttest25</a></td>
            <td><a href="">no Twitter/@</a></td>
        </tr>
    </tbody>
</table>