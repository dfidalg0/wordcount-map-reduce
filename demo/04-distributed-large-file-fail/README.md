# Resultados de execução

## Master

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `master`           |
|    `file`    | `wordcount/files/pg1342.txt` |
| `chunksize`  |           `51200`            |
| `reducejobs` |             `5`              |

## Workers

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `worker`           |
|    `port`    |           `50001`            |
|    `fail`    |             `4`              |

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `worker`           |
|    `port`    |           `50002`            |
|    `fail`    |             `2`              |

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `worker`           |
|    `port`    |           `50003`            |
