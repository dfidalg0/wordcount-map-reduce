# Wordcount - MapReduce

Este repositório contém alterações de uma implementação prévia do framework MapReduce escrita em Go e um programa que utiliza esta implementação para contar a ocorrência de cada uma das palavras únicas presentes em um arquivo de texto qualquer.

**Autor**: Diego Teixeira Nogueira Fidalgo

## Compilação

O executável `wordcount` pode ser compilado através do comando

```bash
make
```

Que o salvará no interior da pasta `bin`, criando-a caso esta não exista.

## Implementação

### `wordcount`

A implementação do algoritmo `wordcount` se resume à implementação das funções `map` e `reduce`

A função `map` extrai do input original (ou de um fragmento deste) as ocorrências individuais de cada palavra e as mapeia em um par chave-valor cuja chave é a palavra e o valor é o número 1, significando uma ocorrência.

A função `reduce`, por sua vez, soma os valores de todas as chaves iguais, resultando em um novo conjunto de chaves e valores onde a chave é uma das palavras que ocorre no texto e o valor é o número de ocorrências desta.

### `mapreduce`

Para o módulo `mapreduce`, foi adicionado um mecanismo de remoção de workers quando estes falham. Para isso, foi necessário apenas remover a chave `worker.id` do mapa de workers na struct `master`.

Além disso, foi modificado o agendador do master para que novas tentativas sejam feitas no caso de operações falhas. Isso pôde ser feito através da criação de um canal **local** que armazena todas as operações e uma `goroutine` que lê todas as operações. Assim, quando uma operação é executada sem sucesso, basta apenas readicioná-la na fila para que a execução desta seja tentada novamente.

Além disso, modificou-se a lógica de uso do `WorkGroup` para que o método `Done` só seja chamado quando a operação for de fato concluída. Desta forma, o método `schedule` só termina quando todas as operações agendadas finalizam sua execução com sucesso, independentemente do número de vezes em que foram executadas.

Finalmente, quando todas as operações são executadas com sucesso, o canal de operações é fechado, de forma que a `goroutine` que o processa é devidamente interrompida.

## Execuções

Foram testadas várias execuções do programa final, as quais resultaram nos arquivos presentes no diretório `/demo` deste repositório. Em cada um dos sub-diretórios deste, encontram-se os resultados de cada uma das seguintes execuções, além de um README que descreve os parâmetros específicos utilizados para obtenção de tais resultados.

Para a verificação da corretude dos resultados, foi utilizada a ferramenta de busca de palavra inteira do editor de texto *Visual Studio Code* para contar o número de ocorrências da palavra cercada por underlines e sem underlines. Assim, é possível fazer a comparação dos resultados como na figura abaixo.

![Comparação de resultados](https://i.imgur.com/XdFDOiR.png)

**Nota**: em **TODAS** as execuções distribuídas, todos os *workers* foram executados antes do *master*.

### Execução sequencial com um arquivo pequeno

Esta execução mostra o funcionamento básico do algoritmo. Nesta, pôde-se comparar os resultados obtidos para todas as palavras com o resultado obtido utilizando-se a função de busca do editor de texto e assim, observar que o algoritmo funciona corretamente.

#### Parâmetros

|  Parâmetro   |           Valor            |
|:------------:|:--------------------------:|
|    `mode`    |        `sequential`        |
|    `file`    | `wordcount/files/test.txt` |
| `chunksize`  |           `100`            |
| `reducejobs` |            `2`             |

#### Logs de Execução

![Logs](https://i.imgur.com/f5qZ5i4.png)

#### Resultados

Todos os resultados brutos desta execução estão neste repositório, em [/demo/00-sequential-small-file](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/00-sequential-small-file). A imagem abaixo mostra um resumo destes.

Nesta figura, pode-se observar a divisão do input para os jobs de map, os inputs mapeados com a chave igual à palavra lida e o valor igual ao número 1 para os jobs de reduce e, finalmente, a redução de todos esses registros para um mapeamento chave-valor correspondente à contagem de cada palavra.

![Resultados](https://i.imgur.com/SBtMa6j.png)

### Execução sequencial com um arquivo grande

Esta execução mostra o correto funcionamento do algoritmo para arquivos grandes, demonstrando o poder de paralelismo do framework MapReduce.

#### Parâmetros

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `sequential`         |
|    `file`    | `wordcount/files/pg1342.txt` |
| `chunksize`  |          `100000`            |
| `reducejobs` |             `5`              |

#### Logs de Execução

Nos logs de execução, vemos que, mesmo com um arquivo grande (~700 kB), o algoritmo foi capaz de realizar a contagem das palavras em 3 segundos.

![Logs](https://i.imgur.com/ZeZ8tKj.png)

#### Resultados

Todos os resultados brutos desta execução estão neste repositório, em [/demo/01-sequential-large-file](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/01-sequential-large-file).

Para a verificação da corretude dos resultados, não foi possível comparar cada um dos resultados manualmente com a ferramenta de busca, mas todos os casos testados mostraram resultados corretos.

### Execução distribuída com falha no map

Esta execução mostra que, mesmo após uma falha na operação map, o sistema consegue se recuperar como um todo ao alocar todas as tarefas restantes ao worker que ainda funciona.

#### Parâmetros

##### Master

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `master`           |
|    `file`    |  `wordcount/files/test.txt`  |
| `chunksize`  |            `100`             |
| `reducejobs` |             `5`              |

##### Workers

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

#### Logs de Execução

Nos logs desta execução, pode-se verificar a ocorrência de falha na operaçao de map com id 1 no worker 0, que é então removido da lista de workers disponíveis pelo master. Após isso, pode-se ver, em seguida, o worker 1 executando a operação de map com id 1, comprovando a capacidade de recuperação do sistema em uma falha no map.

![Logs](https://i.imgur.com/RLVnvco.png)

#### Resultados

Todos os resultados brutos desta execução estão neste repositório, em [/demo/02-distributed-map-fail](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/02-distributed-map-fail).


### Execução distribuída com falha no reduce

Esta execução mostra que o sistema funciona também mesmo após a falha de uma operação de reduce.

#### Parâmetros

##### Master

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `master`           |
|    `file`    |  `wordcount/files/test.txt`  |
| `chunksize`  |            `100`             |
| `reducejobs` |             `5`              |

##### Workers

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

#### Logs de Execução

Nos logs de execução, pode-se ver a falha do worker ao realizar *reduce job* de id 3, seguida pela sua remoção da lista do *master* e da delegação do *job* para o outro worker, comprovando assim a capacidade de recuperação do sistema na operação de reduce.

![Logs](https://i.imgur.com/CupsNSw.png)

#### Resultados

Todos os resultados brutos desta execução estão neste repositório, em [/demo/03-distributed-reduce-fail](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/03-distributed-reduce-fail).

### Execução distribuída com duas falhas

Esta execução mostra que, mesmo após duas falhas dos workers, contanto que ainda exista um worker funcional, o algoritmo consegue se recuperar destas falhas e exibir o resultado correto no final sem grandes problemas.

#### Parâmetros

##### Master

|  Parâmetro   |            Valor             |
|:------------:|:----------------------------:|
|    `mode`    |         `distributed`        |
|    `type`    |           `master`           |
|    `file`    | `wordcount/files/pg1342.txt` |
| `chunksize`  |           `51200`            |
| `reducejobs` |             `5`              |

##### Workers

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

#### Logs de Execução

Nestes logs, é possível ver a falha dos workers nas operações de map 4 e 5 seguidas pela remoção dos workers falhos da lista do *master* e pela delegação dessas operações ao único worker funcional, o que comprova a capacidade de recuperação do sistema como um todo após mais de uma falha.

![Logs](https://i.imgur.com/5lAwBHD.png)

#### Resultados

Todos os resultados brutos desta execução estão neste repositório, em [/demo/04-distributed-large-file-fail](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/04-distributed-large-file-fail).

Para a verificação da corretude dos resultados, não foi possível comparar cada um dos resultados manualmente com a ferramenta de busca, mas todos os casos testados mostraram resultados corretos e coerentes com o resultado da execução sequencial para o mesmo arquivo.
