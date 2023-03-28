package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ResponseAPICEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ResponseViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	canal1 := make(chan *ResponseViaCEP)
	canal2 := make(chan *ResponseAPICEP)
	var cep string
	for _, value := range os.Args[1:] {
		cep = value
	}

	var inputApiCEP string = cep[:5] + "-" + cep[5:]

	//API 1
	go func() {
		resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var c ResponseViaCEP
		err = json.Unmarshal(body, &c)
		if err != nil {
			panic(err)
		}
		canal1 <- &c
	}()

	//API 2
	go func() {
		resp, err := http.Get(fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", inputApiCEP))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var c ResponseAPICEP
		err = json.Unmarshal(body, &c)
		if err != nil {
			panic(err)
		}
		canal2 <- &c
	}()

	select {

	case msg := <-canal1:
		fmt.Printf("ViaCEP: %v\n", *msg)

	case msg := <-canal2:
		fmt.Printf("ApiCEP: %v\n", *msg)

	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}
