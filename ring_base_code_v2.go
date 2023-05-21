// Código para o trabaho de sistemas distribuidos (eleicao em anel)
// Grupo 12:  Daniela Suttoff, Leonardo Cruz e Nathalia Rodrigues

/*
A ordem dos canais são a seguinte
processo 0: in - canal[0]; out - canal[1]
processo 1: in - canal[1]; out - canal[2]
processo 2: in - canal[2]; out - canal[3]
processo 3: in - canal[3]; out - canal[0]
*/




package main

import (
	"fmt"
	"sync"
)

type mensagem struct {
	tipo  int    				// tipo da mensagem
	corpo [6]int 				// conteudo da mensagem para ser usado na eleição
}

/*
Tipo das mensagens
Controle manda -> 			tipo 0: Força a falha de um processo
Controle manda -> 			tipo 1: Religa um Processo falho
Controle manda -> 			tipo 2: Inicia a auto detcção de falha do lider atual
Processo para processo -> 	tipo 3: Usado para a "votação" da eleição
Processo para processo ->	tipo 4: Usado para informar quem ganhou a eleição
Controle manda -> 			tipo 5: Força uma eleição
Controle manda -> 			tipo 6: Finalizar o processo
*/






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
    //o controlador vai fazer os seguintes testes
	/*
	1- testar o auto detector de falha do lider
	{
		- mandar a mensagem do tipo 2 para qualquer processo do anel 
		- mandar a mensagem do tipo 0 para o lider 
	}
	2- falhar um processo que nao seja o lider e iniciar uma eleição
	{
		- mandar uma mensagem do tipo 0 para um processo que nao seja o lider 
		- mandar uma mensagem do tipo 5 para o atual lider 
	}
	3- religar um processo e iniciar a eleição
	    - mandar a mensagem do tipo 1 para um processo que esteja falho
		- mandar a mensagem do tipo 5 para o processo que foi religado 
	}
	4- finalizar todos os processos
	{
		-mandar a mensagem do tipo 6 para todos os processos  
	}
	*/


									//esse defer faz o que está do lado dele apos a conclusão de todo metodo, meio que se tu
									//botar o wg.Done() la no final da no mesmo
	defer wg.Done()   				//isso aqui é o waitGroup, depois que o metodo acaba ele fala que este processo ja finalizou
	
	
	var temp mensagem //uma declaração de uma variavel em Go, esta variavel sera usada em todos os testes

	//-------------------------------------------------------------------- Teste 1 -----------------------------------------------------------
	
    temp.tipo = 2
	chans[2] <- temp  	//mando a mensagem do tipo 2 para um processo qualquer 
    fmt.Println("Controle: Iniciar Auto detecção de Falha do lider")
	temp.tipo = 0
	chans[0] <- temp  	//mando a mensagem do tipo 0 para um processo lider
	fmt.Println("Controle: mudar o processo 0 para falho")
	fmt.Printf("Controle: confirmação de falha %d\n", <-in) 
    fmt.Printf("Controle: confirmação de fim da eleição %d\n", <-in) 


	//-------------------------------------------------------------------- Teste 2 -----------------------------------------------------------
	
	temp.tipo = 0
	chans[2] <- temp // enviando a mensagem pro processo 0
	fmt.Printf("Controle: mudar o processo 2 para falho\n")
	fmt.Printf("Controle: confirmação de falha %d\n", <-in) 
    temp.tipo = 5
	chans[1] <- temp // enviando a mensagem pro processo 1
	fmt.Printf("Controle: Iniciar eleição\n")
	fmt.Printf("Controle: confirmação de fim da eleição %d\n", <-in) 


	//-------------------------------------------------------------------- Teste 3 -----------------------------------------------------------

	temp.tipo = 1
	chans[0] <- temp // enviando a mensagem pro processo 0
	fmt.Printf("Controle: mudar o processo 0 para ativo\n")
	fmt.Printf("Controle: confirmação de ativação %d\n", <-in) 
    temp.tipo = 5
	chans[0] <- temp // enviando a mensagem pro processo 1
	fmt.Printf("Controle: Iniciar eleição\n")
	fmt.Printf("Controle: confirmação de fim da eleição %d\n", <-in) 

	//-------------------------------------------------------------------- Teste 4 -----------------------------------------------------------

	fmt.Println("Controle: finalizando todos os processos")
	temp.tipo = 5 //esse tipo é para fazer todos os processos sairem do loop infinito
	chans[0] <- temp
	chans[1] <- temp
	chans[2] <- temp
	chans[3] <- temp
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done() 

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

	fmt.Println("\n   Processo controlador criado")

	wg.Wait() // Wait for the goroutines to finish\
}
