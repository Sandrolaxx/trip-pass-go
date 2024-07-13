## Dockerfile🐳

```
#1
FROM golang:1.22.5-alpine

#2
WORKDIR /journey

#3º
COPY go.mod go.sum ./

RUN go mod download && go mod verify

#4º
COPY . .

#5º
#WORKDIR /outrodir/app

#6º
RUN go build -o ./bin/journey .

#7º
EXPOSE 8080

#8º
ENTRYPOINT [ "/journey/bin/journey" ]
```

**1º** - FROM é qual é a linguagem de utilizamos no container que vamos criar, "de onde vamos partir", ele busca no dockerhub.
<br/> - Sempre utilizar uma tag de versão especifica, olhar se tem muitos CVE's na versão escolhida.
<br/> - Priorizar versão alpine, que é menor/mais enchuta, possui uma superficie de ataque menor.

**2º** - Diretório de trabalho, se não definido é asumido o diretório raiz, que não é uma boa prática. Ele cria a pasta com o nome definido.

**3º** - Copiando os arquivos do projeto go.mod e go.sum responsáveis por gerenciar as dependências do projeto, para dentro do container, jogando eles na raiz da nossa pasta `journey`.
<br> - Tendo os dois arquivos vamos conseguir executar o comando que realiza esse processo, para executar comandos no dockerfile utilizamos o comando **RUN**.

**4º** - Após instalarmos todas as dependências, precisamos buildar nossa aplicação, porém nossos dados ainda não estão dentro do container, para jogar tudo que temos no nosso projeto para o container, utilizamos o comando **COPY**, onde o primeiro `.` é para pegar tudo que temos na raiz do nosso projeto, e o segundo `.` jogar para dentro da raiz do docker, nesse caso no nosso WORKDIR /journey.

**5º** - É possível mudar de diretório a qualquer momento.

**6º** - Executando o comando de build, onde o `-o` define o output, por ser um binário vamos definir que ele vai ser gerado e adicionado na pasta /bin/journey, lembrando que o caminho `./bin/journey` sim precisa do `.`, pois ele faz referência ao diretório local(journey), e o segundo parâmetro é o `.` que faz referência ao arquivo .go que precisa ser buildado, no caso journey.go.

**7º** - A porta que será exposta pelo container, no caso a 8080 por ser onde nossa aplicação é executada, se não expormos a porta a aplicação vai executar sem ser possível acessar.

**8º** - Entrypoint é o que esse container vai executar, qual o caminho do executavel da aplicação que desejamos executar no container.

---

## Buildando dockerfile🏗

Comando docker para criar imagem
```
docker build -t api-journey:v1 .
```

Podemos listar as imagens
```
docker image ls ou docker images
```

Executar a criação do nosso container
```
docker run --name api-journey -d -p 8080:8080 api-journey:v1
```
> Comando `-d` é para rodar o container detached do terminal.

Verificar os containers rodando
```
docker ps -a
```

Se nosso container não estiver rodando, podemos ver seu log utilizando o container ID
```
docker logs 514aa176536b
```

---

## Multi-Stage build✨

Podemos ver que o tamanho da nossa imagem criada está bem alto, muito do que temos em nossa imagem são arquivos do nosso projeto, porém no final do dia só precisamos do binário.

Temos uma imagem de API simples com quase 600MB de espaço.

![tamanho-imagem-v1](https://github.com/user-attachments/assets/10d296fb-f1b7-4ed7-bd09-641d827c926d)

Para resolver isso podemos alterar nosso dockerfile para ter várias etapas, como build e execução.

```
#1º
FROM golang:1.22.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bin/journey .

#2º
FROM scratch

WORKDIR /app

#3º
COPY --from=builder /bin/journey .

EXPOSE 8080

ENTRYPOINT [ "./journey" ]
```

**1º** - Demos um alias a nosso processo, tudo que tiver de operações até o proximo `FROM` eu estou chamando de `builder`.

**2º** - Nesse segundo estágio inicio nosso container com base no [scratch](https://hub.docker.com/_/scratch), uma imagem docker que tem apenas o básico para um sistema ser executado.

**3º** - Conseguimos se comunicar entre estágios, onde copiamos o que tem disponivel em /bin/journey do estágio `builder` e colamos na raiz do nosso novo estágio.

Após essa alteração podemos criar nossa imagem novamente e visualizarmos se ouve alteração no tamanho de nossa imagem.

![tamanho-imagem-v2](https://github.com/user-attachments/assets/d137bb58-96c3-4b16-8950-49d026535bff)

Tivemos uma diminuição de 96.95% no tamanho da nossa imagem.

