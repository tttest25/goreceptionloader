# #goreceptionloader
![Much Win!](https://golang.org/lib/godoc/images/home-gopher.png)

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