# Rate Limiter - Instructions

### Objetivo

Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Descrição

O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

- Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.

- Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:

```
API_KEY: <TOKEN>
```

### Obs:

As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

### Requisitos:

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web

- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.

- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.

- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.

- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.

- O sistema deve responder adequadamente quando o limite é excedido:
```
Código HTTP: 429
Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
```

- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.

- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.

- A lógica do limiter deve estar separada do middleware.

### Entrega

- O código-fonte completo da implementação.

- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.

- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.

- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.

- O servidor web deve responder na porta 8080.


# Rate Limiter - Documentação

### Como o rate limiter funciona?

Este middleware é executado antes de cada chamada ao servidor web, verificando as informações armazenadas no sistema de armazenamento. Com base na configuração estabelecida, ele decide se deve bloquear ou permitir as requisições.

Em mais detalhes, a cada solicitação ao servidor, o middleware analisa se o número máximo de requisições por segundo para um determinado token ou IP, recuperados da requisição, foi excedido. Para isso, são registradas informações no Redis a cada requisição, possibilitando o monitoramento desse limite. Caso o número máximo de chamadas seja ultrapassado, um bloqueio temporário é imposto de acordo com a configuração especificada. Novas chamadas serão aceitas somente após o término do período de bloqueio da chave.

### Como ele pode ser configurado?

As configurações são baseadas em variáveis de ambiente definidas no arquivo .env, conforme descrito abaixo:

- MAX_REQUESTS_PER_IP → Número máximo de requisições permitidas por IP.
- MAX_REQUESTS_PER_TOKEN → Número máximo de requisições permitidas por token.
- BLOCK_DURATION → Duração do bloqueio quando o número máximo de requisições for excedido.
- REDIS_URL → URL de conexão com o Redis.
- REDIS_TTL → Tempo de expiração das chaves no Redis.

Exemplo de uso:

```
MAX_REQUESTS_PER_IP=2
#Limita a duas requisições por segundo para cada IP.

MAX_REQUESTS_PER_TOKEN=4
#Limita a quatro requisições por segundo para cada token.

BLOCK_DURATION=2m
#Bloqueia novas requisições por 2 minutos para qualquer chave que exceda o limite de requisições.

REDIS_URL=localhost:6379
#URL do Redis, apontando para o servidor local (localhost) na porta padrão 6379.

REDIS_TTL=1h
#Define o tempo de expiração das chaves no Redis como 1 hora.
```

### Como Rodar projeto?

Para rodar o projeto.

- Clone do projeto:
```
git clone git@github.com:lucasdpg/rate-limiter.git
```

- Acessar o diretório do projeto
```
cd rate-limiter
```

- Build e deploy local do projeto
```
docker compose up -d
```

#### Obs:

1. O app respode na porta 8080
2. Para testar o projeto pode se usar o token nas requets e validar o limiter por token, para validar os limites por IP faça as requets sem o token. 