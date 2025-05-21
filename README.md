# HSE Helper
    HSE Helper — это учебный проект, разработанный для автоматизации обработки и анализа студенческих письменных работ.
† created by **Mikael Oganesian** aka **15lu.Akari** †

___
## Оглавление

- [File Analysis Service](#file-analysis-service)
  - [Описание](#описание)
  - [Архитектура](#архитектура)
  - [API](#api)
    - [POST /analyze](#post-analyze)
    - [GET /reportsid](#get-reportsid)
- [File Storing Service](#file-storing-service)
  - [Описание](#описание-1)
  - [Архитектура](#архитектура-1)
  - [API](#api-1)
    - [POST /upload](#post-upload)
    - [GET /files](#get-files)
    - [GET /filesid](#get-filesid)

___

# File Analysis Service

## Описание
Микросервис для анализа студенческих отчётов в формате текста (.txt).  
Основные функции:  
- Подсчёт статистики по тексту (абзацы, слова, символы)  
- Проверка на 100% плагиат по хэшу текста  
- Сохранение результатов анализа в базу данных  

---

## Архитектура
- Входящие данные: JSON с текстом отчёта и именем файла  
- Обработка текста: анализ статистики и вычисление хэша  
- Проверка существующих отчётов с таким же хэшем  
- Сохранение новых результатов в базу PostgreSQL  
- Возврат результата анализа с пометкой о плагиате

---

## API

### POST `/analyze`

Принимает JSON с полями:

```json
{
  "text": "текст отчёта студента",
  "file_name": "имя_файла.txt"
}
```

Возвращает:
```json
{
  "paragraphs": 3,
  "words": 150,
  "characters": 900,
  "is_plagiarized": false
}
```

Если текст уже есть в базе (плагиат по хешу), is_plagiarized будет true.

### GET `/reports/:id`

Возвращает данные анализа отчета по ID:

```json
{
  "id": 1,
  "file_name": "report1.txt",
  "paragraphs": 3,
  "words": 150,
  "characters": 900,
  "hash": "a3b1c...",
  "created_at": "2025-05-20T12:34:56Z"
}
```
___

# File Storing Service

## Описание
Микросервис для хранения файлов студентов в базе данных PostgreSQL.  
Основные функции:  
- Приём и сохранение загружаемых файлов (имя, содержимое, метаданные)  
- Получение списка всех файлов (без содержимого)  
- Получение полного содержимого конкретного файла по ID  
- Надёжное хранение данных с использованием UUID/serial ID и временных меток  

---

## Архитектура
- Приём файлов через HTTP multipart/form-data запросы  
- Сохранение файла в таблицу `files` с полями: ID, имя, содержимое (bytea), дата создания и др.  
- Возвращение информации о загруженных файлах без содержимого для быстрого просмотра  
- Отдельный эндпоинт для загрузки полного файла с контентом по ID  
- Взаимодействие с PostgreSQL через пул соединений (pgxpool)

---

## API

### POST `/upload`

Принимает multipart/form-data с файлом под ключом `file`.

Успешный ответ:

```json
{
  "message": "File uploaded successfully",
  "id": 1,
  "name": "example.txt",
  "created": "2025-05-20T14:00:00Z"
}
```
### GET `/files`

Возвращает список всех файлов без содержимого:

```json
[
  {
    "id": 1,
    "name": "example.txt",
    "created_at": "2025-05-20T14:00:00Z"
  },
  {
    "id": 2,
    "name": "report.txt",
    "created_at": "2025-05-19T18:30:00Z"
  }
]
```

### GET `/files/:id`

Возвращает полный файл с содержимым по ID:

```json
{
  "id": 1,
  "name": "example.txt",
  "content": "SGVsbG8gd29ybGQhCg==", 
  "created_at": "2025-05-20T14:00:00Z"
}

```

___

# Gateway Proxy

## Описание
Микросервис gateway служит единым API для внешнего взаимодействия с системой. Он принимает входящие запросы от frontend'а, делегирует задачи соответствующим микросервисам и возвращает результаты пользователю.
- Основные задачи:
- Роутинг запросов к file-storing и file-analysis
- Агрегация данных
- Унификация формата ответов
- Обработка ошибок на уровне API


---

## Архитектура
- POST /upload → перенаправляет файл в file-storing
- GET /files → получает список файлов без содержимого
- GET /files/:id → получает конкретный файл с содержимым
- POST /analyze → отправляет текст в file-analysis
- GET /reports/:id → получает отчёт анализа по ID

---

## API

### 📁 Загрузка файла

**POST** `/upload`  
Загружает файл в микросервис хранения.

**Request:** `multipart/form-data`

```
file: <binary file>
```

**Response:**
```json
{
  "id": 1,
  "name": "example.txt",
  "created_at": "2025-05-20T12:00:00Z"
}
```

---

### 📄 Получение списка файлов

**GET** `/files`  
Возвращает список загруженных файлов (без содержимого).

**Response:**
```json
[
  {
    "id": 1,
    "name": "example.txt",
    "created_at": "2025-05-20T12:00:00Z"
  }
]
```

---

### 📄 Получение конкретного файла

**GET** `/files/:id`  
Возвращает файл с содержимым:

**Response:**
```json
{
  "id": 1,
  "name": "example.txt",
  "content": "base64-encoded-content-or-plain-text",
  "created_at": "2025-05-20T12:00:00Z"
}
```

---

**GET** `/files/`  
Возвращает список всех файлов:

**Response:**
```json
[
  {
    "id": 1,
    "name": "example.txt",
    "content": "base64-encoded-content-or-plain-text",
    "created_at": "2025-05-20T12:00:00Z"
  },
]
```

---

### 🧠 Анализ текста

**POST** `/analyze`  
Отправляет текст на анализ (подсчёт статистики, проверка на плагиат).

**Request:**
```json
{
  "text": "Текст студенческого отчета",
  "file_name": "report1.txt"
}
```

**Response:**
```json
{
  "paragraphs": 4,
  "words": 200,
  "characters": 1200,
  "is_plagiarized": false
}
```

---

### 📊 Получение отчёта

**GET** `/reports/:id`  
Получение результатов анализа текста по ID.

**Response:**
```json
{
  "id": 1,
  "file_name": "report1.txt",
  "paragraphs": 3,
  "words": 150,
  "characters": 900,
  "hash": "a3b1c...",
  "created_at": "2025-05-20T12:34:56Z"
}
```
