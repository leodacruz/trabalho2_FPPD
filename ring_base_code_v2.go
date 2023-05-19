// Código exemplo para o trabaho de sistemas distribuidos (eleicao em anel)
// By Cesar De Rose - 2022

// para não fechar o deadlock não pode deixar um chanal aberto, faz na sequencia

package main

import (
	"fmt"
	"sync"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo [5]int // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
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
	defer wg.Done()

	var temp mensagem

	// comandos para o anel iciam aqui

	// mudar o processo 0 - canal de entrada 3 - para falho (defini mensagem tipo 2 pra isto)

	temp.tipo = 2
	chans[0] <- temp
	fmt.Printf("Controle: mudar o processo 0 para falho\n")

	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	//pedir para o processo 2 inicia a leição
	temp.tipo = 1
	temp.corpo[4]=1
	chans[1] <- temp
	fmt.Printf("Controle: mudar o processo 1 para iniciar a eleicao\n")

	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// mudar o processo 1 - canal de entrada 0 - para falho (defini mensagem tipo 2 pra isto)

	//temp.tipo = 3
	//chans[2] <- temp
	//fmt.Printf("Controle: mudar o processo 1 para falho\n")
//	fmt.Printf("Controle: confirmação %d\n", <-in) // receber e imprimir confirmação

	// matar os outrs processos com mensagens não conhecidas (só pra cosumir a leitura)

//	temp.tipo = 4
	///chans[3] <- temp
	//fmt.Printf("Controle: confirmação %d\n", <-in)
	//fmt.Println("\n   Processo controlador concluído\n")

	fmt.Println("\n   Processo controlador finalizando canais\n")

	temp.tipo = 99
	chans[3] <- temp
	chans[0] <- temp
	chans[1] <- temp
	chans[2] <- temp

}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done() // isso ele faz depois que acabar o metodo, 
	
	//para ficar sempre ativos os processos 
	var tes bool
	tes = true

	
   

	    var actualLeader int
		var bFailed bool = false // todos inciam sem falha

		actualLeader = leader // indicação do lider veio por parâmatro
        
		
		//while true, o controlador envia um mensagem pra todos que os fazer sair deste loop
		for tes {
			//fmt.Printf("%2d: recebi",TaskId)
	//	fmt.Printf("   teste 2 = %d\n", TaskId)
		temp := <-in // ler mensagem, eles iniciam mas ficam travados aqui esperando uma mensagem
		fmt.Printf("%2d: recebi mensagem %d, [ %d, %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2],temp.corpo[3],temp.corpo[4])

		switch temp.tipo {
		case 0:
			{
                //estou fazendo a eleicao
				//	percorro todo o anel 
                
				//vejo se o processo está falhado ou ativo
				if(bFailed){//aqui ele está falhado dai só passo a mensagem adiante
					var temp1 mensagem
					temp1 = temp
					temp1.tipo = 0
					temp1.corpo[TaskId] = 99
                    out <- temp1
				}else{

				//condição para saber quando acabar a eleicao, dai o vetor tem 5 posiçoes
				//uma para cada processo e a ultima para saber quem iniciou a eleição
                if(temp.corpo[4]==TaskId){
					fmt.Printf(" vou decidir o novo lider \n");	
				//acabou a eleicao e precisamos passar a informação pro outros processos 	
                    var temp1 mensagem
					temp1 = temp
					temp1.tipo = 8 // esse é o tipo pra corfirmar que a eleicao acabou

					//aqui eu percorro e vejo o menor id para ser o novo lider
					//processos falhos terao seu id como 5 por padrão assim nao corre o
					//risco de eleger um processo falho como lider
                    var contador int
					contador=5

					for i := 0; i < 4 ; i++ {
						fmt.Printf("\n\nDados: \ntemp1= %d  \ni= %d  \ncontador %d", temp1.corpo[i],i,contador);		
					if(temp1.corpo[i]<contador){
						contador	= temp1.corpo[i]
						fmt.Printf("\n entrou aqui  temp1 %d  i= %d  contador %d", temp1.corpo[i],i,contador);	
					}	
					}

					//salvo no primeiro indice qcluem venceu e passo pelo anel a informação
					temp1.corpo[0] = contador
					//aqui eu salvo quem fez a eleicao pra depois saber que enviou o novo lider para todos do anel
					temp1.corpo[1] = TaskId
					//envio a informação pra todo mundo
					leader= contador
					out <- temp1
				}else{
					var temp1 mensagem
					temp1 = temp
					temp1.tipo = 0
					temp1.corpo[TaskId] = TaskId
					out <- temp1
				}

			}


			}
		case 1:{
			var temp1 mensagem
					temp1 = temp
					temp1.tipo = 0
					temp1.corpo[TaskId] = TaskId
					out <- temp1
		}
		case 2:
			{
				bFailed = true
				fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
				controle <- -5
			}
		case 3:
			{
				bFailed = false
				fmt.Printf("%2d: falho %v \n", TaskId, bFailed)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
				controle <- -5
			}

		case 8:
			{ //aqui eu att o lider e passo a informação a diante

				//´para saber se ja posso informar ao controlador que a eleição ja acabou e todos ja 
				//estpa cientes do novo lider
				if(temp.corpo[1]==TaskId){
                controle <- 1
				}else{
					fmt.Printf("%2d: eu sei que o novo lider eh %d  e que o processo %d  iniciou a eleicao \n", TaskId, temp.corpo[0],temp.corpo[1])
					leader = temp.corpo[0]
					out <- temp
				}
			}

		case 99:
			{
				tes = false
			}

		default:
			{
				fmt.Printf("%2d: não conheço este tipo de mensagem\n", TaskId)
				fmt.Printf("%2d: lider atual %d\n", TaskId, actualLeader)
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
