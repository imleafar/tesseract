package main

import (
	"math/rand"
	"fmt"
	"os"
	"os/exec"
	"net/http"
	"io"
	"encoding/json"
	"io/ioutil"
	"time"
	)
	
func main() {
	http.HandleFunc("/", obterImagemDaUrl)
	http.HandleFunc("/erros", obterErros)
	http.ListenAndServe(":3000", nil)
}

func catalogarErro(mensagem string) {
	errosAntigos := obterMensagensDosErros()
    bytes := []byte(time.Now().Format("2006-01-02 15:04:05")+ " - " + mensagem+ "\n" + errosAntigos)
    ioutil.WriteFile("erros.txt", bytes, 0)
}

func obterMensagensDosErros() (texto string) {
	binario, erro := ioutil.ReadFile("erros.txt")
	if erro != nil {
		texto = ""
        return
    }
	texto = string(binario)
	return
}

func obterErros(w http.ResponseWriter, r *http.Request) {
	erros := obterMensagensDosErros()
 	w.Write([]byte(erros))
}

func obterImagemDaUrl(w http.ResponseWriter, r *http.Request) {
	imagens, ok := r.URL.Query()["image"]
    if !ok || len(imagens[0]) < 1 {
        json.NewEncoder(w).Encode("Informe o link da imagem. Exemplo: http://localhost:3000/?image=hello.png")
        return
	}

	imagem := string(imagens[0])
	nomeDaImagem := gerarNomeDaImagem(8)
	downloadDaImagem(imagem, nomeDaImagem)
	executarOCR(nomeDaImagem, w)
	os.Remove(nomeDaImagem)
}

func executarOCR(nomeDaImagem string, w http.ResponseWriter) {
	resultadoDoTexto := ""
	tesseract := fmt.Sprintf("tesseract %s stdout -l por", nomeDaImagem)
	cmd := exec.Command("sh", "-c", tesseract)
	out, err := cmd.CombinedOutput()
    if err != nil {
		catalogarErro("erro na imagem: "+nomeDaImagem)
		return
    }
	resultadoDoTexto = string(out)
	json.NewEncoder(w).Encode(resultadoDoTexto)
}

func gerarNomeDaImagem(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)+".jpg"
}

func downloadDaImagem(url string, nomeDaImagem string) {
	output, _ := os.Create(nomeDaImagem)
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro de download", url, "-", err)
		return
	}
	defer response.Body.Close()

	io.Copy(output, response.Body)
}