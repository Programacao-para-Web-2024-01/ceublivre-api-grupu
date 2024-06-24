
# Microsserviço de Avaliações, Comentários, Perguntas e Respostas

Este microsserviço gerencia avaliações, comentários, perguntas e respostas sobre produtos em um marketplace. Inclui funcionalidades de gerenciamento e moderação de conteúdo gerado pelos usuários.

## Funcionalidades

- **Gerenciamento de Avaliações**: Permite que os usuários avaliem produtos com nota e comentário.
- **Gerenciamento de Comentários**: Permite que os usuários comentem nas avaliações de outros usuários.
- **Gerenciamento de Perguntas**: Permite que os usuários façam perguntas sobre os produtos para os vendedores.
- **Gerenciamento de Respostas**: Permite que os vendedores respondam às perguntas dos usuários.
- **Controle de Conteúdo**: Verifica se o conteúdo contém palavras banidas.
- **Moderação de Conteúdo**: Fornece ferramentas para que os moderadores possam revisar e remover conteúdo inadequado.

## Endpoints

### Adicionar uma Avaliação

**Endpoint**: `POST /avaliacoes`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "nota": 5,
    "comentario": "Ótimo produto!"
}
```

### Adicionar uma Avaliação com Palavras Banidas

**Endpoint**: `POST /avaliacoes`

**Body** (JSON):
```json
{
    "id_produto": "2",
    "id_usuario": "usuario2",
    "nota": 1,
    "comentario": "Produto horrível!"
}
```

### Adicionar um Comentário a uma Avaliação

**Endpoint**: `POST /avaliacoes/comentar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario2",
    "texto": "Concordo, é ótimo!"
}
```

### Adicionar uma Pergunta

**Endpoint**: `POST /perguntas/adicionar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "id_vendedor": "vendedor1",
    "duvida": "Este produto está disponível em outras cores?"
}
```

### Adicionar uma Resposta a uma Pergunta

**Endpoint**: `POST /perguntas/responder`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "duvida": "Este produto está disponível em outras cores?",
    "id_vendedor": "vendedor1",
    "resposta": "Sim, está disponível em vermelho, azul e verde."
}
```

### Listar Todas as Avaliações

**Endpoint**: `GET /avaliacoes`

### Listar Todas as Perguntas

**Endpoint**: `GET /perguntas`

### Marcar uma Avaliação como Inadequada

**Endpoint**: `POST /avaliacoes/marcar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1"
}
```

### Marcar uma Pergunta como Inadequada

**Endpoint**: `POST /perguntas/marcar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1"
}
```

### Listar Todas as Avaliações Marcadas para Moderação

**Endpoint**: `GET /avaliacoes/moderar`

### Moderar uma Avaliação (Aprovar)

**Endpoint**: `POST /avaliacoes/moderar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "acao": "aprovar"
}
```

### Moderar uma Avaliação (Remover)

**Endpoint**: `POST /avaliacoes/moderar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "acao": "remover"
}
```

### Listar Todas as Perguntas Marcadas para Moderação

**Endpoint**: `GET /perguntas/moderar`

### Moderar uma Pergunta (Aprovar)

**Endpoint**: `POST /perguntas/moderar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "acao": "aprovar"
}
```

### Moderar uma Pergunta (Remover)

**Endpoint**: `POST /perguntas/moderar`

**Body** (JSON):
```json
{
    "id_produto": "1",
    "id_usuario": "usuario1",
    "acao": "remover"
}
```

## Como Executar

1. **Instale o Go**: Certifique-se de que o Go está instalado no seu sistema.
2. **Clone o repositório**: `git clone <URL_DO_REPOSITORIO>`
3. **Navegue até o diretório do projeto**: `cd <NOME_DO_DIRETORIO>`
4. **Execute o servidor**: `go run main.go`

O servidor estará rodando na porta `8080`.

## Arquivo de Palavras Banidas

Crie um arquivo chamado `palavras.txt` no mesmo diretório do projeto. Este arquivo deve conter uma palavra banida por linha. Exemplo:

```
palavra1
palavra2
palavra3
```

## Notas

- Certifique-se de que o arquivo `palavras.txt` existe e está no mesmo diretório do projeto.
- As palavras banidas são verificadas tanto nas avaliações quanto nas perguntas e respostas.
- O serviço de moderação permite que moderadores aprovem ou removam conteúdo marcado como inadequado.
