# 🔨Trilha DevOps

Para mais informações sobre a trilha de devops consultar seu [readme](/devops/readme.md).

---

# 🦫API Go Lang

Este projeto consiste no desenvolvimento de uma aplicação back-end utilizando a linguagem Go e o framework Go-Chi, com foco em ferramentas modernas e boas práticas para agilizar o desenvolvimento e garantir escalabilidade. API para gerenciamento de viagens desenvolvida na Nwl-Journey.

## 👨‍💻Tecnologias Utilizadas

- **Go**: Linguagem de programação.
- **Go-Chi**: Framework para roteamento HTTP.
- **Docker**: Ferramenta para criar, implantar e executar aplicações em contêineres.
- **Docker Compose**: Ferramenta para definir e executar aplicações Docker multi-contêiner.
- **SQLC**: Ferramenta de geração de código para consultas SQL em Go.
- **goapi-gen**: Ferramenta para geração de código Go a partir de especificações OpenAPI.
- **Migrations**: Gerenciamento de versões do esquema do banco de dados.
- **OpenAPI Specifications**: Especificações para a definição de interfaces de API RESTful.
- **MailPit**: Lib para testar localmente o disparo de e-mails.

## 👨‍🏫Aprendizados

- Compreenção inicial de Go Lang
- Roteamento e manipulação de requisições HTTP com Go-Chi.
- Gerenciamento de contêineres e orquestração com Docker e Docker Compose.
- Geração de código Go para acesso a banco de dados com SQLC.
- Geração de código Go a partir de especificações OpenAPI com goapi-gen.
- Migrations para gerenciamento de versões do banco de dados.
- Especificações OpenAPI para definição clara e precisa das APIs.
- Simulação local do disparo de e-mails

## 🐋Configuração do Ambiente

1. **Pré-requisitos**:
   - [Docker](https://www.docker.com/get-started)
   - [Docker Compose](https://docs.docker.com/compose/install/)
   - [Go](https://golang.org/dl/)

2. **Clone o Repositório**:
   ```bash
    git clone https://github.com/sandrolaxx/trip-pass-go.git
    cd trip-pass-go
   ```

3. **Configure as Variáveis de Ambiente**:
   - Crie um arquivo `.env`.
   - Exemplo abaixo:
    ```
    JOURNEY_DATABASE_USER=pg-test
    JOURNEY_DATABASE_PASSWORD=1329
    JOURNEY_DATABASE_HOST=localhost
    JOURNEY_DATABASE_PORT=5446
    JOURNEY_DATABASE_NAME=postgres
    ```

4. **Execute generaate**:
   ```bash
   go generate ./...
   ```
    - Esse comando vai gerar a spec com base no json do swagger, realizar a migration e a geração do código para acesso ao banco com SQLC.


5. **Inicie a Aplicação**:
   ```bash
   docker-compose up