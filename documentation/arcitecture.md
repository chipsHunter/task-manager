# Архитектура проекта Task Manager

## Часть 1. Проектирование архитектуры («To Be»)

### 1. Определение типа приложения

Task Manager представляет собой Web приложение, включающее:

- **Backend** – Golang и Rust для обработки запросов и взаимодействия с базой данных.
- **Frontend** – веб-приложение на JavaScript для работы через браузер.

### 2. Стратегия развёртывания

- **Backend**: контейнеризация с Docker, деплой на облачные платформы (YandexCloud), база данных PostgreSQL, кэширование Redis.
- **Frontend**: разрабатывается с использованием JS, контейнеризация с Docker, деплой на облачные платформы (YandexCloud).

### 3. Обоснование выбора технологий

#### Backend:

- **Rust** – высокопроизводительный язык для задач, требующих высокой надежности и безопасности.
- **Golang** – компилируемый язык с возможностью разработки различных Backend сервисов и асинхронный приложений.
- **PostgreSQL** – база данных.
- **Redis** – кэширование данных.
- **JWT** – аутентификация.

#### Frontend:

- **HTML** – язык разметки.
- **CSS** – язык описания внешнего вида.
- **JavaScript** – это высокоуровневый интерпретируемый язык программирования, в основном, использующийся в браузерах для придания интерактивности веб-страницам.

### 4. Показатели качества

- **Производительность** – кэширование, клиентская оптимизация.
- **Масштабируемость** – контейнеризация и адаптация под облачные сервисы.
- **Безопасность** – JWT-аутентификация, защищённые соединения.
- **Надёжность** – логирование и обработка ошибок.

### 5. Реализация сквозной функциональности

- **Backend**: логирование (`slog`), кэширование (`Redis`), обработка ошибок, аутентификация с помощью JWT.
- **Frontend**: работа с API, кэширование (`In memmory`).

### 6. Структурная схема приложения
![Архитектурная диаграмма ToBe](https://github.com/chipsHunter/task-manager/blob/main/documentation/pictures/ToBe.png)

## Часть 2. Анализ архитектуры («As Is»)

### 1. Анализ текущей структуры кода

- **Backend**:
  - **Golang**:
    - **http-server/** – бизнес-логика.
    - **config/** – конфигурация.
    - **model/** – модели данных.
    - **storage/** – взаимодействие с PostgreSQL.
  - **Rust**:
    - **api/** – бизнес-логика Rust.
    - **models/** – модели данных.
    - **clients/** – взаимодействие с PostgreSQL.
- **Frontend**:
  - **assets/** – асеты для сайта
  - **css/** – файлы стилей
  - **img/** –изображения для сайта

### 2. Структурная схема текущей реализации

![Архитектурная диаграмма AsIs](https://github.com/chipsHunter/task-manager/blob/main/documentation/pictures/IsAs.png)


## Часть 3. Сравнение и рефакторинг

### 1. Сравнение «To Be» и «As Is»

- Архитектурное видение включает чёткое разделение backend и frontend, а также использование микросервисов для большей гибкости и масштабируемости.
- В текущей реализации присутствуют некоторые ключевые модули, но требуется улучшение взаимодействия между микросервисами и клиентом.

### 2. Выявленные отличия

- Необходимо оптимизировать взаимодействие backend и frontend (стандартизировать API, добавить API-Gateway, улучшить обработку ошибок).
- Кэширование Redis требуется добавть, особенно для сессионных данных и часто запрашиваемых объектов.
- Функциональность сервиса задач следует расширить.
- Доработать архитектуру frontend, внедрив улучшенный интерфейс с усовершенствованным дизайном.

### 3. Пути улучшения архитектуры

- Оптимизация API взаимодействия между backend и frontend.
- Доработка функционала backend.
- Введение кэширования для backend, такого как Redis и улучшенное управление состоянием в frontend (Redux, Context API).
- Автоматизированное тестирование и CI/CD для стабильности развертывания.
