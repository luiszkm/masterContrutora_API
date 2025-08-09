# Guia de Deploy - Master Construtora

## Visão Geral

Este guia descreve como fazer o deploy da aplicação Master Construtora em diferentes ambientes, desde configuração local até produção completa com alta disponibilidade.

## Estratégias de Deploy

### 1. Deploy Local/Desenvolvimento
- **Banco**: Docker container local
- **Aplicação**: Execução direta com `go run`
- **Ideal para**: Desenvolvimento e testes

### 2. Deploy Docker Simples
- **Banco**: Container PostgreSQL
- **Aplicação**: Container Go
- **Ideal para**: Staging, testes de integração

### 3. Deploy Produção
- **Banco**: PostgreSQL gerenciado (AWS RDS, Google Cloud SQL)
- **Aplicação**: Múltiplas instâncias com load balancer
- **Ideal para**: Ambiente de produção

## Configuração Docker

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.18

WORKDIR /root/

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar binário do stage anterior
COPY --from=builder /app/main .

# Copiar arquivos de configuração se necessário
# COPY --from=builder /app/configs ./configs

# Expor porta
EXPOSE 8080

# Comando de execução
CMD ["./main"]
```

### docker-compose.yml para Produção

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/mastercostrutora_db?sslmode=disable
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - APP_ENV=production
    depends_on:
      - db
    restart: unless-stopped
    networks:
      - app-network

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=mastercostrutora_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init:/docker-entrypoint-initdb.d
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### Configuração Nginx

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

    server {
        listen 80;
        server_name your-domain.com;

        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        # SSL configuration
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256;

        # Security headers
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

        # API routes
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # Timeouts
            proxy_connect_timeout 30s;
            proxy_send_timeout 30s;
            proxy_read_timeout 30s;
        }

        # Static files (if serving frontend)
        location / {
            root /var/www/html;
            try_files $uri $uri/ /index.html;
        }

        # Health check
        location /health {
            proxy_pass http://app;
            access_log off;
        }
    }
}
```

## Deploy AWS

### 1. EC2 com Docker

#### Preparação da Instância

```bash
#!/bin/bash
# Script de inicialização EC2

# Atualizar sistema
sudo yum update -y

# Instalar Docker
sudo yum install -y docker
sudo service docker start
sudo usermod -a -G docker ec2-user

# Instalar Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Instalar Git
sudo yum install -y git

# Clone do repositório
git clone <YOUR_REPOSITORY_URL> /home/ec2-user/app
cd /home/ec2-user/app

# Configurar variáveis de ambiente
sudo tee /home/ec2-user/app/.env > /dev/null <<EOF
DATABASE_URL=postgres://user:${DB_PASSWORD}@${RDS_ENDPOINT}:5432/mastercostrutora_db?sslmode=require
JWT_SECRET_KEY=${JWT_SECRET}
APP_ENV=production
EOF

# Iniciar aplicação
docker-compose up -d
```

#### Security Groups

```yaml
# security-groups.yml (CloudFormation/Terraform)
SecurityGroup:
  Type: AWS::EC2::SecurityGroup
  Properties:
    GroupDescription: Master Construtora API Security Group
    VpcId: !Ref VPC
    SecurityGroupIngress:
      # HTTP
      - IpProtocol: tcp
        FromPort: 80
        ToPort: 80
        CidrIp: 0.0.0.0/0
      # HTTPS
      - IpProtocol: tcp
        FromPort: 443
        ToPort: 443
        CidrIp: 0.0.0.0/0
      # SSH (apenas para seu IP)
      - IpProtocol: tcp
        FromPort: 22
        ToPort: 22
        CidrIp: SEU.IP.AQUI/32
    SecurityGroupEgress:
      - IpProtocol: -1
        CidrIp: 0.0.0.0/0
```

### 2. ECS (Elastic Container Service)

#### Task Definition

```json
{
  "family": "master-construtora",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskRole",
  "containerDefinitions": [
    {
      "name": "app",
      "image": "YOUR_ECR_REPO/master-construtora:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "APP_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DATABASE_URL",
          "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:db-url"
        },
        {
          "name": "JWT_SECRET_KEY",
          "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:jwt-secret"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/master-construtora",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3,
        "startPeriod": 60
      }
    }
  ]
}
```

### 3. RDS Configuration

```yaml
# RDS Configuration (CloudFormation)
Database:
  Type: AWS::RDS::DBInstance
  Properties:
    DBInstanceIdentifier: master-construtora-db
    DBInstanceClass: db.t3.micro  # Para produção: db.t3.small ou maior
    Engine: postgres
    EngineVersion: '16.1'
    MasterUsername: masteruser
    MasterUserPassword: !Ref DBPassword
    AllocatedStorage: 20
    StorageType: gp2
    StorageEncrypted: true
    VPCSecurityGroups:
      - !Ref DatabaseSecurityGroup
    DBSubnetGroupName: !Ref DBSubnetGroup
    BackupRetentionPeriod: 7
    DeletionProtection: true
    MultiAZ: true  # Para alta disponibilidade
    
DatabaseSecurityGroup:
  Type: AWS::EC2::SecurityGroup
  Properties:
    GroupDescription: RDS Security Group
    VpcId: !Ref VPC
    SecurityGroupIngress:
      - IpProtocol: tcp
        FromPort: 5432
        ToPort: 5432
        SourceSecurityGroupId: !Ref AppSecurityGroup
```

## Deploy Google Cloud Platform

### 1. Cloud Run

#### Build e Deploy

```bash
# Configurar gcloud
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# Build da imagem
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/master-construtora

# Deploy no Cloud Run
gcloud run deploy master-construtora \
    --image gcr.io/YOUR_PROJECT_ID/master-construtora \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated \
    --set-env-vars APP_ENV=production \
    --set-env-vars DATABASE_URL="postgres://user:pass@/db?host=/cloudsql/PROJECT:REGION:INSTANCE" \
    --add-cloudsql-instances PROJECT:REGION:INSTANCE
```

#### cloudbuild.yaml

```yaml
steps:
  # Build da aplicação
  - name: 'golang:1.23'
    script: |
      go mod download
      CGO_ENABLED=0 go build -o main ./cmd/server
    env:
      - 'GOOS=linux'

  # Build da imagem Docker
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/master-construtora:$COMMIT_SHA', '.']

  # Push da imagem
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/master-construtora:$COMMIT_SHA']

  # Deploy no Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    entrypoint: gcloud
    args:
      - 'run'
      - 'deploy'
      - 'master-construtora'
      - '--image'
      - 'gcr.io/$PROJECT_ID/master-construtora:$COMMIT_SHA'
      - '--region'
      - 'us-central1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'

images:
  - 'gcr.io/$PROJECT_ID/master-construtora:$COMMIT_SHA'
```

### 2. Cloud SQL

```bash
# Criar instância Cloud SQL
gcloud sql instances create master-construtora-db \
    --database-version=POSTGRES_16 \
    --tier=db-f1-micro \
    --region=us-central1 \
    --storage-auto-increase \
    --backup-start-time=03:00

# Criar banco de dados
gcloud sql databases create mastercostrutora_db \
    --instance=master-construtora-db

# Criar usuário
gcloud sql users create appuser \
    --instance=master-construtora-db \
    --password=STRONG_PASSWORD
```

## Deploy Kubernetes

### Kubernetes Manifests

#### Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: master-construtora
  labels:
    app: master-construtora
spec:
  replicas: 3
  selector:
    matchLabels:
      app: master-construtora
  template:
    metadata:
      labels:
        app: master-construtora
    spec:
      containers:
      - name: app
        image: your-registry/master-construtora:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

#### Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: master-construtora-service
spec:
  selector:
    app: master-construtora
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

#### Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: master-construtora-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.construtora.com
    secretName: api-tls
  rules:
  - host: api.construtora.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: master-construtora-service
            port:
              number: 80
```

#### Secrets

```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  database-url: <base64-encoded-url>
  jwt-secret: <base64-encoded-secret>
```

### Deploy no Kubernetes

```bash
# Aplicar manifests
kubectl apply -f k8s/

# Verificar status
kubectl get pods -l app=master-construtora
kubectl get svc master-construtora-service
kubectl get ing master-construtora-ingress

# Ver logs
kubectl logs -f deployment/master-construtora

# Escalar aplicação
kubectl scale deployment master-construtora --replicas=5
```

## Configuração de SSL/TLS

### Let's Encrypt com Certbot

```bash
# Instalar certbot
sudo snap install core; sudo snap refresh core
sudo snap install --classic certbot

# Obter certificado
sudo certbot certonly --standalone -d api.construtora.com

# Renovação automática
sudo crontab -e
# Adicionar linha:
0 12 * * * /usr/bin/certbot renew --quiet
```

### Configuração SSL no Nginx

```nginx
server {
    listen 443 ssl http2;
    server_name api.construtora.com;

    ssl_certificate /etc/letsencrypt/live/api.construtora.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.construtora.com/privkey.pem;
    
    # SSL configuration
    ssl_session_cache shared:le_nginx_SSL:10m;
    ssl_session_timeout 1440m;
    ssl_session_tickets off;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;
    
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;
    
    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    
    # Locations...
}
```

## Monitoramento e Observabilidade

### Health Check Endpoint

A aplicação já inclui um endpoint de health check em `/health`:

```json
{
  "status": "ok",
  "timestamp": "2024-02-20T14:30:00Z",
  "database": "connected",
  "version": "1.0.0"
}
```

### Prometheus Metrics

Para adicionar métricas Prometheus:

```go
// main.go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Métricas customizadas
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// No router, adicionar:
r.Handle("/metrics", promhttp.Handler())
```

### Logging para Produção

```go
// Configuração de logging para produção
func setupLogger() *slog.Logger {
    if os.Getenv("APP_ENV") == "production" {
        return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        }))
    }
    return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
}
```

## Backup e Disaster Recovery

### Backup Automatizado

```bash
#!/bin/bash
# backup-script.sh

DB_HOST="your-db-host"
DB_NAME="mastercostrutora_db"
DB_USER="user"
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Criar backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Upload para S3 (AWS)
aws s3 cp $BACKUP_DIR/backup_$DATE.sql.gz s3://your-backup-bucket/database/

# Limpar backups locais mais antigos que 7 dias
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

# Enviar notificação de sucesso
echo "Backup completed successfully at $(date)" | mail -s "Database Backup Success" admin@construtora.com
```

### Cron Job para Backup

```bash
# Adicionar ao crontab
0 2 * * * /scripts/backup-script.sh >> /var/log/backup.log 2>&1
```

## CI/CD Pipeline

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.23
    
    - name: Run tests
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable
        JWT_SECRET_KEY: test-secret
      run: |
        go mod download
        go test -v ./...

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-east-1
    
    - name: Build and push Docker image
      run: |
        docker build -t master-construtora .
        docker tag master-construtora:latest $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/master-construtora:latest
        aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
        docker push $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/master-construtora:latest
    
    - name: Deploy to ECS
      run: |
        aws ecs update-service --cluster production --service master-construtora --force-new-deployment
```

## Variáveis de Ambiente para Produção

### Arquivo .env.production

```env
# Database
DATABASE_URL=postgres://user:password@production-db:5432/mastercostrutora_db?sslmode=require

# JWT
JWT_SECRET_KEY=super-secret-jwt-key-256-bits-minimum

# App
APP_ENV=production
PORT=8080

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Security
CORS_ALLOWED_ORIGINS=https://app.construtora.com,https://admin.construtora.com
RATE_LIMIT_REQUESTS_PER_SECOND=10
RATE_LIMIT_BURST=20

# External Services
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your-sendgrid-api-key

# Monitoring
PROMETHEUS_ENABLED=true
HEALTH_CHECK_TIMEOUT=5s
```

## Checklist de Deploy

### Pré-Deploy

- [ ] Testes passando
- [ ] Code review aprovado
- [ ] Backup do banco de dados
- [ ] Variáveis de ambiente configuradas
- [ ] SSL/TLS configurado
- [ ] Monitoring configurado

### Durante o Deploy

- [ ] Build bem-sucedido
- [ ] Containers iniciando corretamente
- [ ] Health check respondendo
- [ ] Logs sem erros críticos
- [ ] Conectividade com banco de dados

### Pós-Deploy

- [ ] Smoke tests passando
- [ ] Performance dentro do esperado
- [ ] Métricas normais
- [ ] Logs sendo coletados
- [ ] Backup funcionando
- [ ] Alertas configurados

## Rollback Strategy

### Estratégia de Rollback

1. **Identificar problema**
2. **Parar novos deploys**
3. **Executar rollback**
4. **Verificar estabilidade**
5. **Comunicar incidente**

### Comandos de Rollback

```bash
# Docker Compose
docker-compose down
docker-compose pull app:previous-version
docker-compose up -d

# Kubernetes
kubectl rollout undo deployment/master-construtora
kubectl rollout status deployment/master-construtora

# ECS
aws ecs update-service --cluster production --service master-construtora --task-definition master-construtora:PREVIOUS_REVISION
```

## Segurança em Produção

### Hardening do Sistema

1. **Firewall**: Apenas portas necessárias abertas
2. **Updates**: Sistema sempre atualizado
3. **Users**: Usuários sem privilégios root
4. **SSH**: Chaves em vez de senhas
5. **Fail2ban**: Proteção contra força bruta

### Segurança da Aplicação

1. **HTTPS**: Sempre obrigatório
2. **HSTS**: Headers de segurança
3. **Rate Limiting**: Proteção contra DoS
4. **Input Validation**: Validação rigorosa
5. **Secrets**: Nunca em código fonte

### Monitoramento de Segurança

```bash
# Logs de segurança
tail -f /var/log/nginx/access.log | grep -E "(40[0-9]|50[0-9])"

# Tentativas de login
journalctl -u ssh -f | grep "Failed password"

# Uso de recursos
htop
iotop
```