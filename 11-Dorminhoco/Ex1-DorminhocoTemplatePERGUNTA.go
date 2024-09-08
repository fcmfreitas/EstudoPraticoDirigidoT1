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

type carta string      // carta é um string

var ch [NJ]chan carta  // NJ canais de itens tipo carta  

func embaralhaCartas() [5][5]carta {
	// Cria o baralho com 21 cartas (4 de cada letra e uma carta extra)
	baralho := []carta{
		"10", "10", "10", "10",
		"J", "J", "J", "J",
		"Q", "Q", "Q", "Q",
		"K", "K", "K", "K",
		"A", "A", "A", "A",
		"@", // carta coringa
	}

	// Embaralha o baralho
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(baralho), func(i, j int) { baralho[i], baralho[j] = baralho[j], baralho[i] })

	// Cria a matriz 5x5
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

func encontrarCartaMaisDiferente(mao []carta) carta {
	
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

func jogador(id int, mesa chan carta, cartasIniciais []carta ) {

	mao := cartasIniciais    // estado local - as cartas na mao do jogador

	for {
		if len(mao) == 5{
			
			fmt.Println(id, " joga") // escreve seu identificador
			
			cartaDescartada := encontrarCartaMaisDiferente(mao) //seleciona a melhor carta para passar
			mesa <- cartaDescartada
			mao = removerCarta(mao, cartaDescartada)
			
			fmt.Printf("Jogador %d descartou: %v\n", id, cartaDescartada)

		} else {  
		
			mao = append(mao, <-mesa)
			// recebe carta da mesa
			//  ...
			// e se alguem bate ?
		}
	}
}

func main() {
	
	matriz := embaralhaCartas() // cria baralho e embaralha
	canalMesa := make(chan carta, 1) // cria o canal em que será passada a carta

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

		go jogador(i, canalMesa, cartasIniciais)
	}

	// Mantém o programa rodando
	<-make(chan struct{}) // Bloqueia
	
}


