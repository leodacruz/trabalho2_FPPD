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
type mensagem struct {
	tipo  int    				// tipo da mensagem
	corpo [4]int 				// conteudo da mensagem para ser usado na eleição
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
	temp.tipo = 6 //esse tipo é para fazer todos os processos sairem do loop infinito
	chans[0] <- temp
	chans[1] <- temp
	chans[2] <- temp
	chans[3] <- temp
}

/*
Explicação Auto detecção, processo de eleição, processo de decisão do novo lider, processo de avisar quem é o novo lider
1 - a mensagem de auto detecção fica rodando no anel
2 - todos os processos falhos só passam a mensagem para o proximo, não interagem com ela
3 -  quando o lider recebe a mensagem ele muda o valor que está registrado no seu indice,
caso esteja 0 ele muda para 1, caso esteja 1 ele muda para 0, ele faz isso toda vez que fica 
recebendo a mensagem durante a auto detecção
4 - processos ativos comparam a informação que está salva no corpo da mensagem, na posição
do lider, com a variavel local de status do lider que eles tem, caso sejam diferentes está
tudo certo, caso estejam iguais, o processo inicia a eleição. 
*Como o lider fica mudando este valor a cada rodada, o unico motivo de um processo ter a mesmo
valor, é que o lider está falho, pois como um processo falho só passa a informação adiante o 
valor que o lider sempre muda seria o mesmo
5 - Ao iniciar uma eleição, se cria uma variavel temp1, para zerar o corpo da mensagem, por isso
não se usa o temp. Apos, o processo que iniciou a eleição anota 2 no seu indice do corpo da mensagem, 
e envia a mensagem de votação(3) para o proximo processo
6 - Ao receber o tipo votação (3) o processo que está ativo vai colocar 1 no corpo da mensagem e 
enviar para o Proximo
7 - Processo não ativos só passaram a mensagem adiante, assim em suas posições ficara o valor 0
8 - Para saber quando acaba a votação, ao receber a mensagem o processo ve se o valor 2 está no 
seu indice, caso esteja se inicia o processo de "contagem de votos"
9 - Contagem de votos é saber quais processos estão aptos para concorrer a eleição e depois 
aplicar a regra que o processo que tenha o menor indice ganhe
10 - como o menor indice ganha, basta fazer um for, que percorra o corpo da mensagem
começando de menor ao maior, e o primeiro que tiver o valor 1, é o novo lider
11 - Apos o processo salvar quem é o novo lider ele passa a informação pelo Anel
12- O corpo da  mensagem do novo lider: Indice 0, vai o Id do novo lider
Indice 1 vai o Id de quem começou a enviar a mensagem de aviso do novo lider




*/





func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done() 

						// para ficar sempre ativos os processos, controle do loop infinito
	var tes bool 		// declaração da variavel no Go
	tes = true   		// atribuindo o true na variavel

	var actualLeader int
	var actualStatusLeader int

	var bFailed bool = false 	// todos inciam sem falha

	actualLeader = leader 		// indicação do lider veio por parâmatro

	if(TaskId == leader){
		actualStatusLeader = 1      // usado na autodetecção do lider falho 	
	}else{
		actualStatusLeader = 0      // usado na autodetecção do lider falho 
	}
	

	//loop infinito
	for tes {

		temp := <-in // ler mensagem
		fmt.Printf("Processo %2d: recebi mensagem %d, [ %d, %d, %d, %d ] Atual lider %d e falho = %v\n", 
		TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3],actualLeader,bFailed)

		switch temp.tipo {

			case 0: // tipo 0: Força a falha de um processo
				{
					bFailed = true
					fmt.Printf("Processo %2d: Status falho - %v ", TaskId, bFailed)
					controle <- 0
				}

			case 1: // tipo 1: Religa um Processo falho
				{
					bFailed = false
					fmt.Printf("Processo %2d: Status falho - %v ", TaskId, bFailed)
					controle <- 0
				}

			case 2: // tipo 2: Inicia a auto detcção de falha do lider atual
				{
					//caso o processo esteja falho ele passa adiante a mensagem
					if(bFailed){ 
						out <- temp 
					}else{
						
						// sou o lider, e preciso trocar o valor da mensagem que esta no corpo, no meu indice
                        if(TaskId == actualLeader){ 
							if(temp.corpo[TaskId] == 0){
								temp.corpo[TaskId] = 1
								out <- temp 	
							}else{
								temp.corpo[TaskId] = 0
								out <- temp 
							}
						}

						//Sou um processo e preciso ver se o lider está ativo, caso não
						if(temp.corpo[actualLeader] == actualStatusLeader ){ 
							var temp1 mensagem 
							temp1.tipo = 3 
							temp1.corpo[TaskId] = 2
							out <- temp1
						}else{
							actualStatusLeader = temp.corpo[TaskId]
							out <- temp
						}
					}
				}

			case 3: //tipo 3: Usado para a "votação" da eleição
				{
					//se sou um processo falho
					if bFailed { 
						out <- temp  //caso seja um processo falho, só passo a mensagem adiante sem interferir ela
					} else {

						//aqui significa que a votação ja passou por todo mundo do Anel
						if(temp.corpo[TaskId]==2){
							temp.corpo[TaskId]=1 //mudo para todo mundo apto estiver com 1 no indice do corpo da mensagem
						}else{
							//continuo a votação
							temp.corpo[TaskId]=1
							out <- temp
						}

						//se passar para cá eu sei que acabou a votação
                        
						//Algotimo de decisão do novo lider, lembrando que o menor vence
						for i := 0; i < 4; i++ {
							if temp.corpo[i] == 1 {
								actualStatusLeader = i
								break
							}
						}

						// enviando a mensagem de novo lider ao demais no anel 
                        temp.tipo= 4
						temp.corpo[0]= actualLeader
						temp.corpo[1]=TaskId
						out <- temp

					}
				}

			case 4: //tipo 4: Usado para informar quem ganhou a eleição
				{ 
					if temp.corpo[1] == TaskId {
						controle <- 0
					} else {
						actualLeader = temp.corpo[0]
						out <- temp
					}
				}

			case 5: //tipo 5: Força uma eleição
				{
					var temp1 mensagem 
					temp1.tipo = 3 
					temp1.corpo[TaskId] = 2
					out <- temp1
				}

			case 6://tipo 6: Finalizar o processo
				{
					fmt.Printf("Processo %2d: Finalizando ...  \n", TaskId, actualLeader)
					tes = false
				}

			default:
				{
					fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
					
					controle <- -5
				}
		}

	}
	
	fmt.Printf("Processo %2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Add a count of four, one for each goroutine

	// criar os processo do anel de eleicao
	go ElectionStage(0, chans[0], chans[1], 0) //Processo 0 - este é o lider
	go ElectionStage(1, chans[1], chans[2], 0) //Processo 1  
	go ElectionStage(2, chans[2], chans[3], 0) //Processo 2 
	go ElectionStage(3, chans[3], chans[0], 0) //Processo 3 
	fmt.Println("\n   Anel de processos criado")

	// criar o processo controlador
	go ElectionControler(controle)
	fmt.Println("\n   Processo controlador criado")

	wg.Wait() // Wait for the goroutines to finish\

	fmt.Println("\nPrograma Encerrado")

}
