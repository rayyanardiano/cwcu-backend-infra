pipeline {
    agent any

    environment {
        APP_NAME = 'cwcu-backend'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                sh '''
                    # Build image menggunakan backend/Dockerfile dengan context backend
                    docker build -t ${APP_NAME}:latest -f backend/Dockerfile backend
                '''
            }
        }

        stage('Deploy Container') {
            steps {
                sh '''
                    # Hentikan dan hapus container backend lama jika ada
                    docker stop cwcu_backend || true
                    docker rm cwcu_backend || true

                    # Jalankan container backend di network docker compose
                    docker run -d \
                      --name cwcu_backend \
                      --network cwcu_cwcu_net \
                      --network-alias backend \
                      --restart unless-stopped \
                      -e APP_ENV=development \
                      -e APP_PORT=8080 \
                      -e APP_NAME=CWCU \
                      -e JWT_SECRET=dev-secret-key-change-in-production-at-least-32-chars \
                      -e DB_HOST=postgres \
                      -e DB_PORT=5432 \
                      -e DB_USER=cwcu_user \
                      -e DB_PASSWORD=cwcu_password \
                      -e DB_NAME=cwcu_db \
                      -e DB_SSLMODE=disable \
                      -e REDIS_HOST=redis \
                      -e REDIS_PORT=6379 \
                      -e AMQP_URL=amqp://cwcu_user:cwcu_password@lavinmq:5672/ \
                      -e CORS_ALLOWED_ORIGINS=https://siakad.samuelms.web.id,http://localhost:5173,http://localhost:80 \
                      ${APP_NAME}:latest
                '''
            }
        }
    }

    post {
        success {
            echo "Backend CWCU berhasil di-deploy!"
        }
        failure {
            echo "Deployment Backend Gagal!"
        }
    }
}
