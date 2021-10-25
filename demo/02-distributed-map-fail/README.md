# Resultados de execução

## Master

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `master`           |
|    `file`    |  `wordcount/files/test.txt`  |
| `chunksize`  |            `100`             |
| `reducejobs` |             `5`              |

## Workers

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `worker`           |
|    `port`    |           `50001`            |
|    `fail`    |             `1`              |

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `worker`           |
|    `port`    |           `50002`            |
