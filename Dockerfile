# Используем базовый образ Ubuntu
FROM ubuntu:latest

# Устанавливаем необходимые инструменты для сборки Go
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wget \
    ca-certificates \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем переменную окружения для версии Go
ENV GOLANG_VERSION=1.23.6

# Загружаем и устанавливаем Go
RUN wget -q https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

# Обновляем PATH
ENV PATH=$PATH:/usr/local/go/bin

# Создаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./ 
# Установка нужных модулей
RUN go mod download
# Копируем остальной код
COPY . .
#COPY *.db ./

# Собираем программу
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Объявляем порт, который будет использоваться
EXPOSE ${TODO_PORT}

# Переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
ENV TODO_PASSWORD="123"

# Запускаем программу
CMD ["./main"]
