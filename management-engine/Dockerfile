# Используем официальный образ Rust
FROM rust:1.81 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем Cargo-файлы и собираем зависимости отдельно (для кэширования)
COPY Cargo.toml Cargo.lock ./
COPY api/Cargo.toml api/
COPY clients/Cargo.toml clients/
COPY controllers/Cargo.toml controllers/
COPY models/Cargo.toml models/
COPY shared/Cargo.toml shared/

# Создаем пустые директории для модулей (чтобы корректно отработал cargo fetch)
RUN mkdir -p api/src clients/src controllers/src models/src shared/src

# Загружаем зависимости без сборки кода (ускоряет билды)
RUN cargo fetch

# Копируем исходный код
COPY . .

RUN apt-get update && apt-get install -y pkg-config libssl-dev

ENV SQLX_OFFLINE=true

# Собираем приложение в релизном режиме
RUN cargo build --release --verbose

# Запускаем минимальный образ для продакшена
FROM debian:bullseye-slim AS runtime

WORKDIR /app

# Копируем скомпилированное приложение из builder-стадии
COPY --from=builder /app/target/release/management-engine .

# Запускаем приложение
CMD ["./management-engine"]
