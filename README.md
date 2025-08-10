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

COMPANY_ZIP_URL="https://arquivos.receitafederal.gov.br/dados/cnpj/dados_abertos_cnpj/2025-05/Empresas0.zip"
COMPANY_STORAGE_PATH=data
DATABASE_URL="postgresql://admin:1234567@localhost:5432/arko_tech_challenge"
LOCATION_API_URL="https://servicodados.ibge.gov.br/api/v1/localidades"
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
make run
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

<img width="1906" height="990" alt="image" src="https://github.com/user-attachments/assets/1f899e21-9e4a-4997-a458-29a9e7fa2630" />

## Interface Web

O projeto também inclui uma interface web simples para visualização dos dados:

- Estados: `/estados/`
- Cidades: `/cidades/`
- Distritos: `/distritos/`
- Empresas: `/empresas/`

<img width="1919" height="991" alt="image" src="https://github.com/user-attachments/assets/06aafb77-91f9-4522-a584-f137be72fde2" />

## Comandos Úteis

- `make db-up`: Inicia o banco de dados e executa migrações
- `make db-down`: Para e remove o container do banco
- `python manage.py migrate`: Executa migrações pendentes
- `python manage.py createsuperuser`: Cria um usuário administrador
- `python manage.py collectstatic`: Coleta arquivos estáticos

## Extração e Atualização de Dados

O projeto utiliza um sistema de upsert (update + insert) para manter os dados atualizados. com golang e goroutines O processo funciona da seguinte forma:

1. **Estados**:

   - Dados são obtidos da API do IBGE
   - Inserção usando `ON CONFLICT (id) DO UPDATE`
   - Atualiza nome e sigla se o estado já existir

2. **Cidades**:

   - Dados são obtidos da API do IBGE
   - Inserção usando `ON CONFLICT (id) DO UPDATE`
   - Atualiza nome e estado_id se a cidade já existir
   - Verifica se o estado existe antes de inserir/atualizar

3. **Distritos**:

   - Dados são obtidos da API do IBGE
   - Inserção usando `ON CONFLICT (id) DO UPDATE`
   - Atualiza nome e cidade_id se o distrito já existir
   - Verifica se a cidade existe antes de inserir/atualizar

4. **Empresas**:
   - Dados são obtidos do arquivo CSV da Receita Federal
   - Inserção usando `ON CONFLICT (cnpj) DO UPDATE`
   - Atualiza todos os campos exceto CNPJ se a empresa já existir

Esse processo é feito utilizando batch processing e paralelização, utilize o comando `make db-up` para rodar o banco junto com o algoritmo de extração dos dados
O processo de upsert garante que:

- Dados não são duplicados
- Registros existentes são atualizados com novas informações
- Integridade referencial é mantida (estados -> cidades -> distritos)
- Operação é idempotente (pode ser executada múltiplas vezes com o mesmo resultado)

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
