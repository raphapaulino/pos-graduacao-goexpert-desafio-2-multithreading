package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CEPAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state,omitempty"`
	Estado       string `json:"uf,omitempty"`
	City         string `json:"city,omitempty"`
	Cidade       string `json:"localidade,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	Bairro       string `json:"bairro,omitempty"`
	Street       string `json:"street,omitempty"`
	Rua          string `json:"logradouro,omitempty"`
}

const PROTOCOL = "https://"

const VIA_CEP_DOMAIN = "viacep.com.br"

const BRASIL_API_DOMAIN = "brasilapi.com.br"

func SearchAddress(apiURL string, cep string, ch chan<- *CEPAPI) {
	url := apiURL + cep
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro na requisição para %s: %s\n", apiURL, err)
		ch <- nil
		return
	}
	defer resp.Body.Close()
	var address CEPAPI
	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		fmt.Printf("Erro ao decodificar resposta da %s: %s\n", apiURL, err)
		ch <- nil
		return
	}
	ch <- &address
}

func main() {
	cep := "14092000"

	ch1 := make(chan *CEPAPI)
	ch2 := make(chan *CEPAPI)

	go SearchAddress(PROTOCOL+BRASIL_API_DOMAIN+"/api/cep/v1/", cep, ch1)
	go SearchAddress(PROTOCOL+VIA_CEP_DOMAIN+"/ws/", cep+"/json/", ch2)

	select {
	case address1 := <-ch1:
		printResult(PROTOCOL+BRASIL_API_DOMAIN, address1)
	case address2 := <-ch2:
		printResult(PROTOCOL+VIA_CEP_DOMAIN, address2)
	case <-time.After(time.Second):
		fmt.Println("Timeout excedido. Nenhuma resposta recebida dentro do tempo limite.")
	}
}

func printResult(apiURL string, address *CEPAPI) {
	if apiURL == PROTOCOL+BRASIL_API_DOMAIN {
		fmt.Printf("Resposta da API %s:\n", apiURL)
		fmt.Printf("CEP: %s\n", address.Cep)
		fmt.Printf("Rua: %s\n", address.Street)
		fmt.Printf("Bairro: %s\n", address.Neighborhood)
		fmt.Printf("Cidade: %s\n", address.City)
		fmt.Printf("UF: %s\n", address.State)
	} else {
		fmt.Printf("Resposta da API %s:\n", apiURL)
		fmt.Printf("CEP: %s\n", address.Cep)
		fmt.Printf("Rua: %s\n", address.Rua)
		fmt.Printf("Bairro: %s\n", address.Bairro)
		fmt.Printf("Cidade: %s\n", address.Cidade)
		fmt.Printf("UF: %s\n", address.Estado)
	}
}
