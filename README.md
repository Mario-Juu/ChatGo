###  Projeto Blog em Go
## Descrição
Projeto em Go para demonstração do uso de WebSockets dentro de aplicações web, onde desenvolvi um chat em tempo real com múltiplos usuários, tendo um sistema de mensageria próprio processado pelo servidor.

## Tecnologias
- Go
- HTML
- CSS
- VueJS
- OpenSSL
- WebSockets

## Como usar
1. Abra o VSCode ou sua IDE compatível de preferência
2. Clone o repositório
```sh
git clone https://github.com/Mario-Juu/WebDevGo.git
```
3. Crie as chaves públicas e privadas com o OpenSSL
```sh
# gera a chave privada
openssl genrsa -out server.key 2048

# gera a chave pública a partir da chave privada
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
```

4. Execute a aplicação
```sh
go run .
```
5. Visite o site na porta 8080 (https://localhost:8080)


## Objetivo 
Sendo uma aplicação fullstack, tem como objetivo demonstrar a utilidade dos websockets dentro de aplicações web, criando um chat dentre os usuários logados em que as informações são passadas em tempo real, criando um vínculo entre usuário-servidor mais dinâmico do que requisições http padrão. Para o site manter-se com as informações atualizadas, utilizei do framework VueJS para atualizar o html conforme necessário. 

