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
- [Gateway Proxy](#gateway-proxy)
  - [Описание](#описание-2)
  - [Архитектура](#архитектура-2)
  - [API](#api-2)
    - [POST /upload](#post-upload-1)
    - [GET /files](#get-files-1)
    - [GET /filesid](#get-filesid-1)
    - [GET /analyze](#get-analyze)
    - [GET /reportsid](#get-reportsid-1)

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

### POST `/upload`  
Загружает файл в микросервис хранения, анализирует загруженный файл и возвращает результат анализа:

**Request:** `multipart/form-data`

```
file: <binary file>
```

**Response:**
```json
{
    "analysis": {
        "paragraphs": 1,
        "words": 6,
        "characters": 31,
        "is_plagiarized": true
    },
    "created_at": "2025-05-22T08:35:57.392859Z",
    "file_id": 28,
    "file_name": "finaltest.txt"
}
```

---

### GET `/files`  
Возращает список всех загруженных файлов:

```json
[
    {
        "id": 30,
        "name": "finaltest.txt",
        "created_at": "2025-05-22T08:38:21.11115Z"
    },
    {
        "id": 29,
        "name": "finaltest.txt",
        "created_at": "2025-05-22T08:38:20.136184Z"
    },
]
```

---


### GET `/files/id`  
Загружает файл с сервера:

```
file: <binary file>
```

---

### GET `/analyze`  
Принимает текст для анализа и возвращает результаты анализа:

**Request:**
```json
{
    "file_name": "testboba.txt",
    "text": "alarmf"
}
```

**Responce:**
```json
{
    "paragraphs": 1,
    "words": 1,
    "characters": 6,
    "is_plagiarized": true
}
```
---


### GET `/reports/:id`  
Возращает проверку конкретной работы:

```json
{
    "id": 14,
    "file_name": "testboba.txt",
    "paragraphs": 1,
    "words": 1,
    "characters": 6,
    "hash": "6d730bbd6d77b45734a7647d26ef23b13ce40cbe0864c3e57455030a6815d381",
    "created_at": "2025-05-22T08:46:31.765357Z"
}
```

---