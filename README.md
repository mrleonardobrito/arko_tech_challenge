# Arko Tech Challenge

API REST para gerenciamento de estados, cidades, distritos e empresas.

## Requisitos

- Python 3.12+
- PostgreSQL 16+
- Docker e Docker Compose

## Configuração do Ambiente

1. Clone o repositório:

```bash
git clone <url-do-repositorio>
cd arko_tech_challenge
```

2. Crie um ambiente virtual e ative-o:

```bash
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
# ou
.venv\Scripts\activate  # Windows
```

3. Instale as dependências:

```bash
pip install -r requirements.txt
```

4. Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
POSTGRES_DB=arko_tech_challenge
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

## Executando o Projeto

1. Inicie o banco de dados:

```bash
make db-up
```

Este comando irá:

- Iniciar o container do PostgreSQL
- Aguardar o banco ficar disponível
- Executar as migrações

2. Execute o servidor de desenvolvimento:

```bash
python manage.py runserver
```

O servidor estará disponível em `http://localhost:8000`.

## Estrutura do Projeto

```
arko_tech_challenge/
├── apps/
│   ├── api/          # Configurações da API (exceções, mensagens)
│   ├── city/         # App de cidades
│   ├── company/      # App de empresas
│   ├── district/     # App de distritos
│   ├── state/        # App de estados
│   └── web/          # App web (templates, static files)
├── arko_tech_challenge/
│   ├── settings.py   # Configurações do Django
│   └── urls.py       # URLs principais
├── docker-compose.yml
├── Makefile
├── requirements.txt
└── README.md
```

## Endpoints da API

### Estados

- `GET /api/states/` - Lista estados (paginado)
  - Parâmetros:
    - `page`: Número da página (default: 1)
    - `page_size`: Registros por página (default: 10, max: 100)

### Cidades

- `GET /api/cities/` - Lista cidades (paginado)
  - Parâmetros:
    - `page`: Número da página (default: 1)
    - `page_size`: Registros por página (default: 10, max: 100)

### Distritos

- `GET /api/districts/` - Lista distritos (paginado)
  - Parâmetros:
    - `page`: Número da página (default: 1)
    - `page_size`: Registros por página (default: 10, max: 100)

### Empresas

- `GET /api/companies/` - Lista empresas (paginado)
  - Parâmetros:
    - `page`: Número da página (default: 1)
    - `page_size`: Registros por página (default: 10, max: 100)

## Documentação da API

A documentação completa da API está disponível em:

- Swagger UI: `/api/docs/`
- ReDoc: `/api/redoc/`
- OpenAPI Schema: `/api/schema/`

## Interface Web

O projeto também inclui uma interface web simples para visualização dos dados:

- Estados: `/estados/`
- Cidades: `/cidades/`
- Distritos: `/distritos/`
- Empresas: `/empresas/`

## Comandos Úteis

- `make db-up`: Inicia o banco de dados e executa migrações
- `make db-down`: Para e remove o container do banco
- `python manage.py migrate`: Executa migrações pendentes
- `python manage.py createsuperuser`: Cria um usuário administrador
- `python manage.py collectstatic`: Coleta arquivos estáticos

## Tratamento de Erros

A API retorna erros no seguinte formato:

```json
{
  "error": {
    "type": "TipoDoErro",
    "message": "Mensagem amigável do erro"
  }
}
```

Códigos de status HTTP:

- 200: Sucesso
- 400: Erro de validação ou requisição inválida
- 404: Recurso não encontrado
- 500: Erro interno do servidor
