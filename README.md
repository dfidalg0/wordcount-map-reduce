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

### [Execução sequencial com um arquivo pequeno](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/00-sequential-small-file)

Esta execução mostra o funcionamento básico do algoritmo. Nesta, pôde-se comparar os resultados obtidos para cada uma das palavras com o resultado obtido utilizando-se da função de busca em um editor de texto como o *Visual Studio Code* e assim, observar que o algoritmo funciona corretamente.

### [Execução sequencial com um arquivo grande](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/01-sequential-large-file)

Esta execução mostra o correto funcionamento do algoritmo para arquivos grandes, demonstrando o poder de paralelismo do framework MapReduce.

### [Execução distribuída com falha no map](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/02-distributed-map-fail)

Esta execução mostra que, mesmo após uma falha na operação map, o sistema consegue se recuperar como um todo ao alocar todas as tarefas restantes ao worker que ainda funciona.

### [Execução distribuída com falha no reduce](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/03-distributed-reduce-fail)

Esta execução mostra que o sistema funciona também mesmo após a falha de uma operação de reduce.

### [Execução distribuída com duas falhas](https://github.com/diegood12/wordcount-map-reduce/tree/main/demo/04-distributed-large-file-fail)

Esta execução mostra que, mesmo após duas falhas dos workers, contanto que ainda exista um worker funcional, o algoritmo consegue se recuperar destas falhas e exibir o resultado correto no final sem grandes problemas.
