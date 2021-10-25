package mapreduce

import (
	"log"
	"sync"
)

// Schedules map operations on remote workers. This will run until InputFilePathChan
// is closed. If there is no worker available, it'll block.
func (master *Master) schedule(task *Task, proc string, filePathChan chan string) int {
	var (
		wg      sync.WaitGroup
		counter int
	)

	log.Printf("Scheduling %v operations\n", proc)

	// Criamos um canal local para enfileirar as operações
	operationsQueue := make(chan *Operation, RETRY_OPERATION_BUFFER)

	counter = 0

	// Criamos, então, uma goroutine para lidar com elas
	go func() {
		// Nesta rotina, esperamos cada uma das operações ser adicionada
		// e a executamos em seguida em uma nova rotina
		for operation := range operationsQueue {
			// E executamos cada uma delas em uma goroutine separada
			go master.runOperation(operation, &wg, operationsQueue)
		}
	}()

	// Obtemos então as operações que devem ser realizadas
	for filePath := range filePathChan {
		operation := &Operation{proc, counter, filePath}
		counter++

		wg.Add(1)

		// E as colocamos no canal de operações
		operationsQueue <- operation
	}

	wg.Wait()

	// Fechamos o canal de operações para que a goroutine consiga, ao menos um
	// dia, ver a sua aposentadoria chegar
	close(operationsQueue)

	log.Printf("%vx %v operations completed\n", counter, proc)
	return counter
}

// runOperation start a single operation on a RemoteWorker and wait for it to return or fail.
func (master *Master) runOperation(
	operation *Operation,
	wg *sync.WaitGroup,
	operationsQueue chan<- *Operation,
) {
	var (
		err  error
		args *RunArgs
	)

	worker := <-master.idleWorkerChan

	log.Printf("Running %v (ID: '%v' File: '%v' Worker: '%v')\n", operation.proc, operation.id, operation.filePath, worker.id)

	args = &RunArgs{operation.id, operation.filePath}
	err = worker.callRemoteWorker(operation.proc, args, new(struct{}))

	if err != nil {
		log.Printf("Operation %v '%v' Failed. Error: %v\n", operation.proc, operation.id, err)
		// No lugar de finalizar uma task do workgroup, readicionamos a operação
		// à fila de operações
		operationsQueue <- operation
		master.failedWorkerChan <- worker
	} else {
		wg.Done()
		master.idleWorkerChan <- worker
	}
}
