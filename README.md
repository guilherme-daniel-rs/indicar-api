# Indicar API

API para serviço de avaliação de veículos com integração completa ao Amazon S3.

## Funcionalidades

- ✅ Autenticação JWT
- ✅ Gerenciamento de usuários e avaliadores
- ✅ Sistema de avaliações de veículos
- ✅ Upload de fotos para S3 (JPEG, PNG, GIF, WebP)
- ✅ Upload de relatórios PDF para S3
- ✅ URLs pré-assinadas para download seguro
- ✅ Validação de tipos e tamanhos de arquivo
- ✅ Documentação Swagger completa

## Configuração Rápida

### 1. Configurar Variáveis de Ambiente

Copie o arquivo de exemplo e configure suas credenciais:

```bash
cp .env-development .env
```

Edite o arquivo `.env` com suas configurações:

```bash
# Database
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=indicar_db

# JWT
JWT_SECRET=your-super-secret-jwt-key

# AWS S3 (credentials handled automatically by IAM roles)
AWS_REGION=us-east-1
AWS_S3_BUCKET=indicar-evaluation-photos
```

### 2. Executar Migrações

```bash
go run main.go -migrate
```

### 3. Iniciar Servidor

```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

## Documentação

- [Documentação Completa da API](API_DOCUMENTATION.md)
- [Configuração do S3](S3_SETUP.md)
- [Swagger UI](http://localhost:8080/swagger/index.html)

## Endpoints Principais

### Autenticação
- `POST /auth/signup` - Registrar usuário
- `POST /auth/login` - Fazer login
- `POST /auth/refresh` - Renovar token

### Avaliações
- `POST /evaluations` - Criar avaliação
- `GET /evaluations` - Listar avaliações
- `POST /evaluations/{id}/photos` - Upload de foto
- `GET /evaluations/{id}/photos` - Listar fotos

### Relatórios
- `POST /reports` - Criar relatório
- `POST /reports/{id}/file` - Upload de PDF
- `GET /reports/{id}/file` - Download de PDF

## Tecnologias

- **Go 1.24+** - Linguagem principal
- **Gin** - Framework web
- **GORM** - ORM para banco de dados
- **AWS S3** - Armazenamento de arquivos
- **JWT** - Autenticação
- **Swagger** - Documentação da API

## Estrutura do Projeto

```
indicar-api/
├── internal/
│   ├── application/
│   │   ├── controllers/     # Controladores HTTP
│   │   └── services/        # Lógica de negócio
│   ├── domain/
│   │   └── entities/        # Entidades do domínio
│   └── infrastructure/
│       ├── aws/            # Integração S3
│       ├── database/       # Conexão e migrações
│       ├── middleware/     # Middlewares
│       └── routes/         # Definição de rotas
├── configs/                # Configurações
├── docs/                   # Documentação Swagger
└── main.go                 # Ponto de entrada
```

## Desenvolvimento

### Executar em Modo Desenvolvimento

```bash
# Instalar dependências
go mod tidy

# Executar migrações
go run main.go -migrate

# Iniciar servidor
go run main.go
```

### Executar Testes

```bash
go test ./...
```

### Gerar Documentação Swagger

```bash
swag init
```

## Deploy

### GitHub Actions (Recomendado)

A aplicação está configurada para deploy automático via GitHub Actions:

1. **Configure os secrets** no repositório GitHub (veja [GITHUB_SECRETS.md](GITHUB_SECRETS.md))
2. **Faça push** para a branch `main`
3. **O deploy acontece automaticamente**

### Docker

```bash
# Build da imagem
docker build -t indicar-api .

# Executar container
docker run -p 8080:8080 --env-file .env indicar-api
```

### Docker Compose

```bash
docker-compose up -d
```

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.