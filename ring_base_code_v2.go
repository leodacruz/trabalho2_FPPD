// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2022

package main

import (
	"fmt"
	"sync"
	"time"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [6]int // conteudo da mensagem para colocar os ids, no caso tem mais um espaço que está sendo
	//usado para salvar quem iniciou a eleição, e o ultimo para salvar o atual lider
}

var (
	chans = []chan mensagem{ // vetor de canias para formar o anel de eleicao - chan[0], chan[1] and chan[2] ...
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	}
	controle = make(chan int)
	wg       sync.WaitGroup // wg is used to wait for the program to finish
)

func ElectionControler(in chan int) {
	//esse defer faz o que está do lado dele apos a conclusão de todo metodo, meio que se tu
	//botar o wg.Done() la no final da no mesmo
	defer wg.Done()   //isso aqui é o waitGroup, depois que o metodo acaba ele fala que este processo ja finalizou
	var temp mensagem //uma declaração de uma variavel em Go


	//inicio a autoEleição, os processo ficam monitorando o lider
    temp.tipo = 1
	chans[2] <- temp // enviando a mensagem pro processo 0
    time.Sleep(2 * time.Microsecond)
	// mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	temp.tipo = 0
	chans[1] <- temp // enviando a mensagem pro processo 0
	fmt.Printf("Controle: mudar o processo 0 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,
   // <- in //isso para esperar eles encontrarem sozinho o erro

	temp.tipo = 0
	chans[0] <- temp // enviando a mensagem pro processo 0
	fmt.Printf("Controle: mudar o processo 0 para falho\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,
    <- in //isso para esperar eles encontrarem sozinho o erro



    //reativo os processo
//	temp.tipo = 1
//	chans[2] <- temp


    // mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	temp.tipo = 4
	chans[0] <- temp // enviando a mensagem pro processo 0
	fmt.Printf("Controle: mudar o processo 0 para ativo\n")
	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,

	   // mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	   temp.tipo = 7
	   chans[0] <- temp // enviando a mensagem pro processo 0
	   fmt.Printf("Controle: forca uma nova eleicao\n")
	   fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,
    


   // <- in //espero o acabar a eleição
	 // mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)
	/* temp.tipo = 4
	 chans[2] <- temp // enviando a mensagem pro processo 0
	 fmt.Printf("Controle: mudar o processo 0 para falho\n")
	 fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,
     <- in */




	//pedir para o processo 1 inicie a eleição
//	temp.tipo = 1
//	temp.corpo[4] = 1 //eu to salvando na mensagem que quem iniciou a eleição foi o processo do indice 1
//	chans[1] <- temp  // enviando mensagem pro processo 1
//	fmt.Printf("Controle: mudar o processo 1 para iniciar a eleicao\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

//	fmt.Println(" Controle: qual seu status")
//	temp.tipo = 6 //esse tipo é para fazer todos os processos sairem do loop infinito
//	chans[3] <- temp
//	<- in
//	chans[0] <- temp
//	<- in
//	chans[1] <- temp
//	<- in
//	chans[2] <- temp
//	<- in


//	temp.tipo = 0
//	chans[2] <- temp // enviando a mensagem pro processo 0
//	fmt.Printf("Controle: mudar o processo 0 para falho\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,

//	temp.tipo = 1
//	temp.corpo[4] = 1 //eu to salvando na mensagem que quem iniciou a eleição foi o processo do indice 1
//	chans[3] <- temp  // enviando mensagem pro processo 1
//	fmt.Printf("Controle: mudar o processo 1 para iniciar a eleicao\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação
    
//	temp.tipo = 4
//	chans[0] <- temp // enviando a mensagem pro processo 0
//	fmt.Printf("Controle: mudar o processo 0 para ativo\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação,

//	temp.tipo = 1
//	temp.corpo[4] = 1 //eu to salvando na mensagem que quem iniciou a eleição foi o processo do indice 1
//	chans[1] <- temp  // enviando mensagem pro processo 1
//	fmt.Printf("Controle: mudar o processo 1 para iniciar a eleicao\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	fmt.Println(" Processo controlador finalizando canais")
	temp.tipo = 5 //esse tipo é para fazer todos os processos sairem do loop infinito
	chans[3] <- temp
	chans[0] <- temp
	chans[1] <- temp
	chans[2] <- temp
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done() // isso ele faz depois que acabar o metodo

	//para ficar sempre ativos os processos, controle do loop infinito
	var tes bool // declaração da variavel no Go
	tes = true   // atribuindo o true na variavel

	var actualLeader int
	var bFailed bool = false // todos inciam sem falha

	actualLeader = leader // indicação do lider veio por parâmatro

	//loop infinito
	for tes {

		temp := <-in // ler mensagem
		fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d, %d, %d ] Atual lider %d e falho = %v\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3], temp.corpo[4],actualLeader,bFailed)

		switch temp.tipo {

		case 0: // controle avisa para este processo falhar
			{
				bFailed = true
				fmt.Printf("{ %2d: falho %v ", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d }\n", TaskId, actualLeader)
				controle <- -5
			}

		case 1: // fica monitorando o lider automagicamente para ver se ele está falho
			{
                //parte para anotar os erros de cada processo
				//var temp2 mensagem
				if(bFailed){ //falhou é 99 no seu espaco da mensagem
                temp.corpo[TaskId]=99
				out <- temp
				}else{
					//situação que caso eu não seja o lider veja se ele nao esta falho
					if(temp.corpo[actualLeader]==99 && actualLeader!=TaskId ){ //isso significa que o lider falhou e esse processo vai iniciar a eleição
						var temp1 mensagem
						temp1 = temp //copio a mensagem que o controle me enviou, tem o id de quem iniciou a eleição 
						temp1.tipo = 2 
						temp1.corpo[TaskId] = TaskId
						temp1.corpo[4] = TaskId
						out <- temp1
					}else{
					temp.tipo=1
					//temp.corpo[TaskId]=TaskId
					out <- temp
					}
				}
			}

		case 2:
			{
				//estou fazendo a eleicao, e dou a volta no Anel 

				//caso o processo esteja com falha ele "pula a eleição"
				if bFailed { 
					var temp1 mensagem
					temp1 = temp
					temp1.tipo = 2
					temp1.corpo[TaskId] = 99 // isso é pra esse processo nao interferir na eleição
					out <- temp1
				} else {
					//aqui o processo está apto a partcipar da eleição


					// aqui é pra caso ja deu a volta no Anel e ja da pra ver quem é o lider
					if temp.corpo[4] == TaskId {
						
						var temp1 mensagem
						temp1 = temp
						temp1.tipo = 3 // esse é o tipo pra corfirmar que a eleicao acabou

						//aqui eu percorro e vejo o menor id para ser o novo lider
						//processos falhos terao seu id como 99 por padrão assim nao corre o
						//risco de eleger um processo falho como lider
						var contador int
						contador = 4 // numero maior que o id max que seria o 3

						for i := 0; i < 4; i++ {
							if temp.corpo[i] < contador {
								contador = temp.corpo[i]
							}
						}

						//salvo no primeiro indice quem venceu e passo pelo anel a informação
						temp1.corpo[5] = contador
						//aqui eu salvo quem fez a eleicao pra depois saber que enviou o novo lider para todos do anel
						temp1.corpo[4] = TaskId
						//salvo quem é o novo lider
						actualLeader = contador
						//envio a informação pra todo mundo
						out <- temp1
					} else {
						//aqui é pra cada processo botar o Id no Anel e passar adiante
						var temp1 mensagem
						temp1 = temp
						temp1.tipo = 2
						temp1.corpo[TaskId] = TaskId
						out <- temp1
					}

				}

			}

		case 3:
			{ //aqui eu att o lider e passo a informação a diante

				//para saber se ja posso informar ao controlador que a eleição ja acabou e todos ja
				//estao cientes do novo lider
				if temp.corpo[4] == TaskId {
					controle <- 1
				} else {
					actualLeader = temp.corpo[5]
					out <- temp
				}
			}

		case 4:
			{
				bFailed = false
				fmt.Printf("{ %2d: falho %v ", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d }\n", TaskId, actualLeader)
				controle <- -5
			}

		case 5:
			{
				fmt.Printf("%2d: eu sei que o lider eh %d   \n", TaskId, actualLeader)
				tes = false
			}
		case 6:
			{
				fmt.Printf("{ %2d: falho %v ", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d }\n", TaskId, actualLeader)
				controle <- 0
			}
		case 7: //força uma eleição
			{
				var temp1 mensagem
						temp1 = temp //copio a mensagem que o controle me enviou, tem o id de quem iniciou a eleição 
						temp1.tipo = 2 
						temp1.corpo[TaskId] = TaskId
						temp1.corpo[4] = TaskId
						out <- temp1


			}

		

		default:
			{
				fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
				
				controle <- -5
			}
		}

	}
	//fmt.Print("teste ")

	// variaveis locais que indicam se este processo é o lider e se esta ativo

	fmt.Printf("%2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao

	go ElectionStage(0, chans[0], chans[1], 0) // este é o lider
	go ElectionStage(1, chans[1], chans[2], 0) // não é lider, é o processo 0
	go ElectionStage(2, chans[2], chans[3], 0) // não é lider, é o processo 0
	go ElectionStage(3, chans[3], chans[0], 0) // não é lider, é o processo 0

	fmt.Println("\n   Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado\n")

	wg.Wait() // Wait for the goroutines to finish\
}
