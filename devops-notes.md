## Dockerfile🐳

```dockerfile
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

**1º** - FROM é qual é a linguagem que utilizamos no container que vamos criar, "de onde vamos partir", ele busca no dockerhub.
<br/> - Sempre utilizar uma tag de versão específica, olhar se tem muitos CVE's na versão escolhida.
<br/> - Priorizar versão alpine, que é menor/mais enxuta, possui uma superfície de ataque menor.

**2º** - Diretório de trabalho, se não definido, é adotado o diretório raiz, que não é uma boa prática. Ele cria a pasta com o nome definido.

**3º** - Copiando os arquivos do projeto go.mod e go.sum responsáveis por gerenciar as dependências do projeto, para dentro do container, jogando na raiz da nossa pasta `journey`.
<br> - Tendo os dois arquivos, vamos conseguir executar o comando que realiza esse processo. Para executar comandos no dockerfile utilizamos o comando **RUN**.

**4º** - Após instalarmos todas as dependências, precisamos buildar nossa aplicação, porém nossos dados ainda não estão dentro do container, para jogar tudo que temos no nosso projeto para o container, utilizamos o comando **COPY**, onde o primeiro `.` é para pegar tudo que temos na raiz do nosso projeto, e o segundo `.` jogar para dentro da raiz do docker, nesse caso no nosso WORKDIR /journey.

**5º** - É possível mudar de diretório a qualquer momento.

**6º** - Executando o comando de build, onde o `-o` define o output, por ser um binário vamos definir que ele vai ser gerado e adicionado na pasta /bin/journey, lembrando que o caminho `./bin/journey` sim precisa do `.`, pois ele faz referência ao diretório local(journey), e o segundo parâmetro é o `.` que faz referência ao arquivo .go que precisa ser buildado, no caso journey.go.

**7º** - A porta que será exposta pelo container, no caso a 8080 por ser onde nossa aplicação é executada, se não expormos a porta, a aplicação vai executar sem ser possível acessar.

**8º** - Entrypoint é o que esse container vai executar, qual o caminho do executável da aplicação que desejamos executar no container.

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

Podemos ver que o tamanho da nossa imagem criada está bem alto, muito do que temos em nossa imagem são arquivos do nosso projeto, porém, no final do dia, só precisamos do binário.

Temos uma imagem de API simples com quase 600MB de espaço.

![tamanho-imagem-v1](https://github.com/user-attachments/assets/10d296fb-f1b7-4ed7-bd09-641d827c926d)

Para resolver isso podemos alterar nosso dockerfile para ter várias etapas, como build e execução.

```dockerfile
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

**1º** - Demos um alias a nosso processo, tudo que tiver de operações até o próximo `FROM` eu estou chamando de `builder`.

**2º** - Nesse segundo estágio, inicio nosso container com base no [scratch](https://hub.docker.com/_/scratch), uma imagem docker que tem apenas o básico para um sistema ser executado.

**3º** - Conseguimos se comunicar entre estágios, onde copiamos o que tem disponível em /bin/journey do estágio `builder` e colamos na raiz do nosso novo estágio.

Após essa alteração, podemos criar nossa imagem novamente e visualizarmos se houve alteração no tamanho de nossa imagem.

![tamanho-imagem-v2](https://github.com/user-attachments/assets/d137bb58-96c3-4b16-8950-49d026535bff)

Tivemos uma diminuição de 96.95% no tamanho da nossa imagem.

---

## CI com github actions

É possível criar esse arquivo pelo github, mas caso queiramos criar a pasta e o arquivo `.github/workflows/main.yml` também é possível.

```yaml
#1º
name: CI

#2º
on:
  push:
    branches:
      - master

#3º
jobs:
  #4º  
  build-and-push:
    name: "Build and Push"
    runs-on:  ubuntu-latest

    #5º
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Build docker image
        run: docker build -t sandrolax/api-journey:latest .

```

**1º** - Nome do nosso workflow

**2º** - Define quando é trigado nosos workflow, no caso, quando houver um push na branch master.

**3º** - O job é uma máquina em execução e essa máquina tem vários steps, podemos definir vários jobs com vários steps(teste/build/etc).

**4º** - Nome do meu job e onde ele vai ser executado.

**5º** - Steps são os passos que desejo realizar quando o job processar, no nosso caso fazemos o `Checkout` que é um [actions](https://github.com/actions/checkout)(steps pré-prontos) que basicamente é um check out da branch no workspace, após isso realizamos o step de buildar a imagem da nossa api.

Como o objetivo é em breve enviar para o container registry do dockerhub, precisamos colocar o nosso nome de usuário do dockerhub na frente do nome da imagem.

---

### Melhorias na action

**Gerar TAG imagem com base hash do commit**: Para isso criamos um step anterior a criação da imagem e utilizando a váriavel $GITHUB_SHA que está presente no contexto, tempos acesso ao hash, após isso pegamos os 7 primeiros caracteres.

step:
```yaml
- name: Generate SHA
  id: generate_sha
  run: |
    SHA=$(echo $GITHUB_SHA | head -c7)
    echo "sha=$SHA" >> $GITHUB_OUTPUT
```

* Id para identificar os valores criados nesse passo.
* O pipe é utilizado para definirmos comandos que tenham mais de uma linha.
* SHA recebe os 7 primeiros caracteres do hash do commit
* Criamos uma variável para adicionar o valor de SHA no output desse step. Todo step tem o output do anterior, uma maneira centralizada de ir passando informações entre os steps.

Após criado e adicionado na variável `GITHUB_OUTPUT`, podemos utilizar para definir a tag da criação da nossa imagem. Abaixo temos um exemplo de como acessar esse valor no step de build.

```yaml
- name: Build docker image
  run: docker build -t sandrolax/api-journey:${{ steps.generate_sha.outputs.sha }} .
```

**Login no container registry**: Para fazer isso, vamos utilizar o action [Docker Login](https://github.com/marketplace/actions/docker-login), nele precisamos passar o usuário e senha(token) do registry que vamos utilizar.

Utilizando o DockerHub passo o nome do meu user e o token é gerado no dockerhub em Account settings > Security > New Access Token.

Para utilizar essas informações, como uma boa prática, vamos utilizar os secrects do github para setar os dados, isso está disponível em Settings(Do repositório) > Security > Actions > New Repository secret.

Abaixo o step criado:
```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}
```

**Enviando imagem para container registry**: Criamos o step e adicionamos o docker push, comando que envia a imagem criada para o dockerhub. Quanto à tag podemos defini-la utilizando o comando docker tag.

Step criado:
```yaml
- name: Push to registry
  run: | 
    docker push sandrolax/api-journey:${{ steps.generate_sha.outputs.sha }}
    docker tag sandrolax/api-journey:${{ steps.generate_sha.outputs.sha }} sandrolax/api-journey:latest
    docker push sandrolax/api-journey:latest
```
> Boa prática

**Utilizando action para realizar o build e push**: O que fizemos manualmente até o momento funciona, porém, não é a melhor prática e está bem verboso. Vamos utilizar a [action ](https://github.com/marketplace/actions/build-and-push-docker-images) para melhorar essa parte do nosso workflow.

Revisando, ficará da seguinte forma:
```yaml
- name: Build and push
  uses: docker/build-push-action@v6
  with:
    context: .
    push: true
    tags: |
    sandrolax/api-journey:${{ steps.generate_sha.outputs.sha }}
    sandrolax/api-journey:latest
```

**Steps adicionais**: Poderiamos adicionar também um step para realizar os testes unitários, para isso é necessário ter o go na máquina onde roda o job, então precisamos instalar o go e então rodar os testes unitários.

Exemplo abaixo:
```
- name: Setup Go
uses: actions/setup-go@v5
with:
    go-version: "1.22.5"

- name: Run tests
run: go test
```

---

## Kubernetes⚓

### Arquitetura kube

![Components of kube](https://kubernetes.io/images/docs/components-of-kubernetes.svg)

**Node**: São os componentes de trabalho, os nós se comunicam com o control plane através do kubelet e tem uma camada de rede que é o kube-proxy. Conceitos como deployment, pod, deamonSet, service, ingress estão nesse componente. Para mais informações, consulte a [documentação](https://kubernetes.io/docs/concepts/architecture/nodes/).
**Control Plane**: É o que gerencia globalmente nosso cluster, componentes de rede, scheduler, api, etcd(banco chave valor).  É o cara que, se cair, temos um baita problema. Para mais informações, consulte a [documentação](https://dockerlabs.collabnix.com/kubernetes/beginners/Kubernetes_Control_Plane.html).
**Scheduler**: É quem tenta alocar nossa aplicação em um determinado nó. [Documentação](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/)

---

### Namespace

É uma divisão lógica para garantir uma melhor organização na execução dos nossos pods. [Mais detalhes](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/namespace-v1/).

Criando via comand line:
```
kubectl create namespace journey
```

---

### Secret

São objetos onde podemos adicionar dados sensíveis, possuem uma estrutura chave/valor e encoda em base64 nossos segredos. Por não ser uma estrutura que criptografa os dados, não é recomendada a utilização em produção. Para mais detalhes, acesse a [documentação](https://kubernetes.io/docs/concepts/configuration/secret/).

Exemplo de secret:
```yaml
apiVersion: v1
kind: Service

metadata:
  name: journey-service
  labels:
    app: journey

spec:
  selector:
    app: journey
  type: ClusterIP
  ports:
    - name: journey-service
      port: 80
      targetPort: 8080
      protocol: TCP
```

Para aplicar o secret podemos utilizar o comando:
```
kubectl apply -f k8s/secret.yaml -n journey
```

> Realizar no diretório que possui o yaml e definir o namespace

---

### Deployment

É a forma declarativa de definir o funcionamento de um replicaset e seus respectivos pods. Nele podemos utilizar outros recursos criados em nosso cluster, como, por exemplo, os secrets. Para mais sobre, [documentação](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/).

Exemplo deployment:
```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: journey-deployment
  labels:
    app: journey

spec:
  replicas: 5
  selector:
    matchLabels:
      app: journey
  template:
    metadata:
      labels:
        app: journey
    spec:
      containers:
        - name: api-journey
          image: sandrolax/api-journey:4388865
          env:
            - name: JOURNEY_DATABASE_USER
              valueFrom:
                secretKeyRef:
                  name: db-connection
                  key: db_user
            - name: JOURNEY_DATABASE_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db-connection
                  key: db_password
            - name: JOURNEY_DATABASE_HOST
              valueFrom:
                secretKeyRef:
                  name: db-connection
                  key: db_host
            - name: JOURNEY_DATABASE_PORT
              valueFrom:
                secretKeyRef:
                  name: db-connection
                  key: db_port
            - name: JOURNEY_DATABASE_NAME
              valueFrom:
                secretKeyRef:
                  name: db-connection
                  key: db_name
          ports:
            - containerPort: 8080
          resources: 
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 128Mi
```

Para aplicar é similar a secret:
```
kubectl apply -f k8s -n journey
```

> Por padrão, ele busca o arquivo deployment na pasta e o executa.

---

### Service

É uma maneira de expor a rede do cluster, para conseguirmos acessar a aplicação que está nos pods, na confirguração do service precisamos definir como irá funcionar a rede interna deles, após isso realizar um [port-forward](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) para externalizar a rede interna para o da máquina que executa o k8s. Documentação sobre [service](https://kubernetes.io/docs/concepts/services-networking/service/).

Exemplo de service:
```yaml
apiVersion: v1
kind: Service

metadata:
  name: journey-service
  labels:
    app: journey

spec:
  selector:
    app: journey
  type: ClusterIP
  ports:
    - name: journey-service
      port: 80
      targetPort: 8080
      protocol: TCP
```

Aplicando service:
```
kubectl apply -f k8s/service.yaml -n journey
```

Executando o port-forward:
```
kubectl port-forward svc/journey-service 8080:80 -n journey
```

> Sempre importante lembrar de passar o namespace na execução dos comandos

Para mais informações sobre o kube consultar o repositório [study-k8s](https://github.com/Sandrolaxx/study-k8s) e também a [documentação oficial](https://kubernetes.io/docs/home).