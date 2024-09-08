// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.


package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5           // numero de jogadores
const M = 4            // numero de cartas
var bateu = false;
var sequenciaBater []int 

type carta string      // carta é um string

var turno = make(chan int, 1)
var mesa = make(chan carta, 1)
var dorminhoco = make(chan struct{}) //sinaliza o fim do jogo

func embaralhaCartas() [5][5]carta {
	// Cria o baralho com 21 cartas (4 de cada letra e uma carta extra)
	baralho := []carta{
		"10", "10", "10", "10",
		"J", "J", "J", "J",
		"Q", "Q", "Q", "Q",
		"K", "K", "K", "K",
		"A", "A", "A", "A",
		"@", 
	}

	// Embaralha o baralho
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(baralho), func(i, j int) { baralho[i], baralho[j] = baralho[j], baralho[i] })

	var matriz [5][5]carta

	// Preenche a matriz com as cartas embaralhadas
	for i := 0; i < 5; i++ {
		if i == 0 {
			for j := 0; j < 5; j++ {
				index := i*5 + j
				if index < len(baralho) {
					matriz[i][j] = baralho[index]
				}
			}
		} else{
			for j := 0; j < 4; j++ {
				index := i*4 + j
				if index < len(baralho) {
					matriz[i][j] = baralho[index]
				}
			}
		}
			
	}

	// Imprime a matriz
	imprimirMatriz(matriz)

	return matriz
}

func imprimirMatriz(matriz [5][5]carta) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(matriz[i][j], " ")
		}
		fmt.Println()
	}
}

func encontrarMaisDiferente(mao []carta) carta {
	
	frequencias := make(map[carta]int)
	
	for i := 0; i < len(mao); i++ {
		cartaAtual := mao[i]
		frequencias[cartaAtual]++
	}
	var cartaMenosFrequente carta
	minFrequencia := len(mao) + 1
	
	for carta, frequencia := range frequencias {
		if frequencia < minFrequencia {
			minFrequencia = frequencia
			cartaMenosFrequente = carta
		}
	}
	
	return cartaMenosFrequente
}

func cartasIguais(mao []carta) bool {

	primeiraCarta := mao[0]
	for _, c := range mao {
		if c != primeiraCarta {
			return false
		}
	}
	return true
}

func printMao(id int, mao []carta) {
	fmt.Printf("Mão do Jogador %d: ", id)
	for _, c := range mao {
		fmt.Print(c, " ")
	}
	fmt.Println()
}

func jogador(id int, cartasIniciais []carta) {

	mao := cartasIniciais    // mão do jogador

	for {
		turnoAtual := <- turno

		if bateu {
			sequenciaBater = append(sequenciaBater, id)
			if len(sequenciaBater) == NJ {dorminhoco <- struct{}{}}
			return
		}
		
		if turnoAtual == id{

			if len(mao) != NJ {
				mao = append(mao, <-mesa)
			}  
			
			fmt.Println(id, " jogando...") 
				
				cartaDescartada := encontrarMaisDiferente(mao) //seleciona a melhor carta para passar
				for i, c := range mao {
					if c == cartaDescartada {
						mao = append(mao[:i], mao[i+1:]...)
						break
					}
				}
				mesa <- cartaDescartada
				
				printMao(id, mao)
				fmt.Printf("Jogador %d descartou: %v\n\n", id, cartaDescartada)

				if cartasIguais(mao) {
					fmt.Println("BATEU!")
					sequenciaBater = append(sequenciaBater, id)
					bateu = true;
					close(turno)
					return
				}
	
				turnoAtual = (id + 1) % NJ
				
		}
		turno <- turnoAtual
	}
}

func main() {
	
	matriz := embaralhaCartas() // cria baralho e embaralha
	turno <- 0

	// distribui as cartas e cria os jogadores
	for i := 0; i < NJ; i++ {
		var cartasIniciais []carta
		
		if i == 0{
			for j := 0; j < M+1; j++ {
				cartasIniciais = append(cartasIniciais, matriz[i][j])
			}
		} else {
			for j := 0; j < M; j++ {
				cartasIniciais = append(cartasIniciais, matriz[i][j])
			}
		}

		go jogador(i, cartasIniciais)
	}

	// Mantém o programa rodando
	<-dorminhoco // Bloqueia
	fmt.Printf("\nSequencia de batida: %d, %d, %d, %d, %d\n", sequenciaBater[0], sequenciaBater[1], sequenciaBater[2], sequenciaBater[3], sequenciaBater[4])
	fmt.Printf("Jogador %d levou ROLHADA!", sequenciaBater[4])
	
}


