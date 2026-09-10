#  Slotify — Telegram Mini App для записи к мастерам

<div align="center">

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-Bot-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)

**Современная система онлайн-записи для мастеров с уведомлениями в Telegram**

[Возможности](#-возможности) • [Архитектура](#-архитектура) • [Быстрый старт](#-быстрый-старт) • [Deploy](#-deploy)

Бот: https://t.me/myslotify_bot

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

<div align="center">

**Сделано с ❤️ для мастеров и их клиентов**

[⬆ Вернуться наверх](#-slotify--telegram-mini-app-для-записи-к-мастерам)

</div>