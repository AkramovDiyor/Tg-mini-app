#  Slotify — Telegram Mini App для записи к мастерам

<div align="center">

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-Bot-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)

**Современная система онлайн-записи для мастеров с уведомлениями в Telegram**

[Возможности](#-возможности) • [Архитектура](#-архитектура) • [Быстрый старт](#-быстрый-старт) • [Deploy](#-deploy)

Бот: (https://t.me/myslotify_bot)

</div>

---

##  Возможности

###  Для мастеров
-  **Персональная ссылка** — уникальный invite-линк для клиентов
-  **Панель управления** — расписание, услуги, настройки
-  **Портфолио** — загрузка фотографий работ
-  **Уведомления** — автоматические пуши при записи/отмене
-  **Напоминания** — за 2 часа до визита (фоновый worker)

### Для клиентов
-  **Мгновенная запись** — без регистрации в боте
-  **Удобный календарь** — выбор даты и времени
-  **Подтверждение** — просмотр и отмена своих записей
-  **Mini App** — работает прямо внутри Telegram

###  Технологический стек

| Слой | Технологии |
|------|-----------|
| **Backend** | Go 1.21, Chi Router, PostgreSQL (pgx/v5) |
| **Frontend** | React 18, Vite, Tailwind CSS, TanStack Query, Zustand |
| **Интеграции** | Telegram Bot API, Telegram WebApp SDK |
| **Инфраструктура** | ngrok (dev), Vercel + Railway (prod) |

---

##  Архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                         Telegram                            │
│   Мастер/Клиент → Named Web App → initData с подписью       │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (React + Vite)                    │
│   • TanStack Query (server state)                           │
│   • Zustand (UI state)                                      │
│   • React.lazy (code splitting)                             │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTTPS
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go)                             │
│   • HTTP API (Chi Router)                                   │
│   • Telegram Bot Handler                                    │
│   • Worker (Scheduler для напоминаний)                      │
│   • Auth Middleware (проверка HMAC подписи)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
                    ┌──────────────┐
                    │  PostgreSQL  │
                    └──────────────┘
```

###  Принцип "Сервер — владыка"

Frontend **не принимает решений** о роли пользователя. При каждом открытии приложения:

1. Frontend отправляет `initData` на `GET /api/v1/me`
2. Backend расшифровывает криптоподпись Telegram
3. Backend проверяет `telegram_id` в таблице `masters`
4. Возвращает: `{"role": "master"|"client", "invite_link": "..."}`
5. Frontend рендерит нужный интерфейс

---

##  Быстрый старт

###  Требования

- **Go 1.21+** — [установить](https://go.dev/doc/install)
- **Node.js 18+** — [установить](https://nodejs.org/)
- **PostgreSQL 14+** — [установить](https://www.postgresql.org/download/) или использовать [Supabase](https://supabase.com/)
- **ngrok** — [установить](https://ngrok.com/download) и [зарегистрироваться](https://dashboard.ngrok.com/signup)
- **Telegram Bot** — создать через [@BotFather](https://t.me/BotFather)

---

###  Шаг 1: Клонирование и установка

```bash
git clone https://github.com/yourusername/slotify.git
cd slotify

# Backend
cd backend
go mod download

# Frontend (в новом терминале)
cd frontend
npm install
```

---

###  Шаг 2: Настройка базы данных

Создай базу данных и выполни миграции:

```sql
-- Основные таблицы
CREATE TABLE masters (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    address TEXT,
    invite_link VARCHAR(8) UNIQUE NOT NULL,
    work_hours JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    master_id BIGINT REFERENCES masters(id),
    name VARCHAR(255) NOT NULL,
    duration_min INTEGER NOT NULL,
    price INTEGER NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE slots (
    id SERIAL PRIMARY KEY,
    master_id BIGINT REFERENCES masters(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'free',
    UNIQUE(master_id, start_time)
);

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    slot_id BIGINT REFERENCES slots(id),
    service_id BIGINT REFERENCES services(id),
    master_id BIGINT REFERENCES masters(id),
    client_telegram_id BIGINT NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    price_locked INTEGER NOT NULL,
    status VARCHAR(30) DEFAULT 'active',
    reminder_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE master_photos (
    id SERIAL PRIMARY KEY,
    master_id BIGINT REFERENCES masters(id),
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE waitlist (
    id SERIAL PRIMARY KEY,
    master_id BIGINT REFERENCES masters(id),
    client_telegram_id BIGINT NOT NULL,
    service_id BIGINT REFERENCES services(id),
    desired_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_bookings_reminder ON bookings (status, reminder_sent, start_time) 
WHERE status = 'active';
CREATE INDEX idx_slots_master_date ON slots (master_id, start_time);
CREATE INDEX idx_masters_invite ON masters (invite_link);
```

---

###  Шаг 3: Переменные окружения

Создай файл `backend/.env`:

```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/slotify?sslmode=disable

# Telegram
TG_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
WEB_APP_URL=https://t.me/YourBotName/slotify

# Server
PORT=8080
```

---

###  Шаг 4: Настройка Telegram бота

#### 4.1. Создай бота

1. Открой [@BotFather](https://t.me/BotFather)
2. Отправь `/newbot`
3. Следуй инструкциям
4. Скопируй **токен** в `.env` как `TG_BOT_TOKEN`

#### 4.2. Настрой Named Web App

1. В BotFather: `/mybots` → выбери бота → `Bot Settings`
2. `Menu Button` → `Configure Menu Button`
3. URL: `https://t.me/YourBotName/slotify` (оставь пока так, обновим после ngrok)
4. Text: `🛠 Открыть приложение`

#### 4.3. Создай Web App через BotFather

```
/mybots → Your Bot → Bot Settings → Menu Button → Edit Menu Button
```

Или используй прямые ссылки:
```
https://t.me/YourBotName/slotify
```

---

###  Шаг 5: Запуск с 2 ngrok (КРИТИЧНО ВАЖНО!)

Для локальной разработки нужны **2 отдельных ngrok туннеля**:

#### Терминал 1: Frontend (Vite)

```bash
cd frontend
npm run dev
```

Увидишь:
```
  VITE v5.0.0  ready in 500ms
  ➜  Local:   http://localhost:5173/
```

#### Терминал 2: ngrok для Frontend

```bash
ngrok http http://localhost:5173
```

Увидишь:
```
Forwarding  https://abcd-12-34-56-78.ngrok-free.app → http://localhost:5173
```

 **Скопируй этот HTTPS URL** — это твой `FRONTEND_URL`

#### Терминал 3: Backend (Go)

```bash
cd backend
go run server/main.go
```

Увидишь:
```
 Bot authorized as @YourBotName
 HTTP server started on :8080
 Reminder worker started
```

#### Терминал 4: ngrok для Backend

```bash
ngrok http http://localhost:8080
```

Увидишь:
```
Forwarding  https://xyz-98-76-54-32.ngrok-free.app → http://localhost:8080
```

 **Скопируй этот HTTPS URL** — это твой `BACKEND_URL`

---

###  Шаг 6: Связываем всё вместе

#### 6.1. Frontend → Backend

В файле `frontend/.env.local`:

```env
VITE_API_BASE=https://xyz-98-76-54-32.ngrok-free.app/api/v1
VITE_STATIC_BASE=https://xyz-98-76-54-32.ngrok-free.app
```

**Перезапусти Vite** (Ctrl+C → `npm run dev`)

#### 6.2. Backend → Telegram (Named Web App)

Обнови `backend/.env`:

```env
WEB_APP_URL=https://abcd-12-34-56-78.ngrok-free.app
```

**Перезапусти Go сервер** (Ctrl+C → `go run server/main.go`)

#### 6.3. Обнови BotFather

1. `/mybots` → Your Bot → `Bot Settings` → `Menu Button`
2. URL: `https://abcd-12-34-56-78.ngrok-free.app` (твой FRONTEND_URL)
3. Text: `🛠 Открыть приложение`

---

---

## 📁 Структура проекта

```
slotify/
├── backend/
│   ├── server/
│   │   └── main.go                 # Точка входа
│   ├── internal/
│   │   ├── config/                 # Конфигурация
│   │   ├── delivery/
│   │   │   ├── httpapi/           # HTTP handlers + middleware
│   │   │   └── telegram/          # Telegram bot handler
│   │   ├── repositories/          # Работа с БД
│   │   ├── services/              # Бизнес-логика
│   │   └── worker/                # Фоновые задачи (scheduler)
│   ├── pkg/
│   │   ├── database/             # Подключение к PostgreSQL
│   │   └── telegram/             # Валидация initData
│   └── .env
│
└── frontend/
    ├── src/
    │   ├── components/           # UI компоненты
    │   ├── hooks/                # Custom hooks (React Query)
    │   ├── lib/                  # Утилиты (telegram.js, http.js)
    │   ├── page/                 # Страницы (Master, Services, Booking)
    │   ├── store/                # Zustand store
    │   └── App.jsx               # Главный компонент
    ├── public/
    └── .env.local
```

---

##  API Endpoints

### Публичные

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/v1/invite/{link}/info` | Инфо о мастере |
| `GET` | `/api/v1/invite/{link}/services` | Услуги мастера |
| `GET` | `/api/v1/invite/{link}/slots` | Свободные слоты |

### Защищённые (требуют `X-Telegram-Init-Data`)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/v1/me` | Определение роли пользователя |
| `POST` | `/api/v1/book` | Создать запись |
| `GET` | `/api/v1/client/bookings` | Мои записи |
| `POST` | `/api/v1/client/bookings/{id}/cancel` | Отменить запись |
| `GET/PUT` | `/api/v1/master/profile` | Профиль мастера |
| `GET/POST/PUT/DELETE` | `/api/v1/master/services` | CRUD услуг |
| `PUT` | `/api/v1/master/settings` | Настройки |
| `GET/POST/DELETE` | `/api/v1/master/photos` | Фото работ |

---

##  Deploy в продакшен

### Backend (Railway / Render / Fly.io)

1. Залей код на GitHub
2. Подключи репозиторий к Railway/Render
3. Добавь переменные окружения (без ngrok, с реальными URL)
4. Deploy!

### Frontend (Vercel)

```bash
cd frontend
vercel --prod
```

Добавь в Vercel Environment Variables:
```
VITE_API_BASE=https://your-backend.railway.app/api/v1
VITE_STATIC_BASE=https://your-backend.railway.app
```

### BotFather (финальная настройка)

Обнови Named Web App URL:
```
https://your-app.vercel.app
```

---

<div align="center">

**Сделано с ❤️ для мастеров и их клиентов**

[⬆ Вернуться наверх](#-slotify--telegram-mini-app-для-записи-к-мастерам)

</div>