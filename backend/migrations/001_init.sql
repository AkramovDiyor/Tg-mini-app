-- CREATE TABLE masters (
--     id BIGSERIAL PRIMARY KEY,
--     telegram_id BIGINT UNIQUE NOT NULL,
--     name VARCHAR(255) NOT NULL,
--     bio TEXT,
--     address TEXT,
--     invite_link VARCHAR(50) UNIQUE NOT NULL,
--     work_hours JSONB NOT NULL DEFAULT '{}',
--     settings JSONB NOT NULL DEFAULT '{}',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );

-- CREATE TABLE services (
--     id BIGSERIAL PRIMARY KEY,
--     master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
--     name VARCHAR(255) NOT NULL,
--     duration_min INT NOT NULL,
--     price INT NOT NULL,
--     description TEXT DEFAULT '',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );

-- CREATE TABLE slots (
--     id BIGSERIAL PRIMARY KEY,
--     master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
--     start_time TIMESTAMP WITH TIME ZONE NOT NULL,
--     end_time TIMESTAMP WITH TIME ZONE NOT NULL,
--     status VARCHAR(50) NOT NULL DEFAULT 'free',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     CONSTRAINT check_slot_times CHECK (end_time > start_time)
-- );

-- CREATE TABLE clients (
--     telegram_id BIGINT PRIMARY KEY,
--     username VARCHAR(255) DEFAULT '',
--     first_name VARCHAR(255) NOT NULL DEFAULT '',
--     last_name VARCHAR(255) DEFAULT '',
--     strikes_count INT NOT NULL DEFAULT 0,
--     is_banned BOOLEAN NOT NULL DEFAULT FALSE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );
-- CREATE TABLE bookings (
--     id BIGSERIAL PRIMARY KEY,
--     slot_id BIGINT NOT NULL UNIQUE REFERENCES slots(id) ON DELETE RESTRICT,
--     service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
--     client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE RESTRICT,
--     client_name VARCHAR(255) NOT NULL, -- Обязательное поле для имени Васи!
--     price_locked INT NOT NULL,
--     status VARCHAR(50) NOT NULL DEFAULT 'active',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
-- );
-- CREATE TABLE IF NOT EXISTS waitlist (
--     id BIGSERIAL PRIMARY KEY,
--     master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
--     client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE CASCADE,
--     desired_date DATE NOT NULL,
--     status VARCHAR(50) NOT NULL DEFAULT 'waiting', -- waiting, offered, fulfilled, expired
--     offered_slot_id BIGINT REFERENCES slots(id) ON DELETE SET NULL,
--     created_at TIMESTAMPTZ DEFAULT NOW()
-- );
-- ALTER TABLE waitlist ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();


-- -- 1. Создаем мастера Педро с нормальной ссылкой (без слэшей)
-- INSERT INTO masters (telegram_id, name, bio, address, invite_link, work_hours, settings)
-- VALUES (9999991, 'Педро Барбер', 'Барбер · 6 лет опыта', 'ул. Центральная, 1', 'LINK1', '{}', '{}')
-- ON CONFLICT (telegram_id) DO UPDATE SET invite_link = 'LINK1';

-- -- 2. Создаем услугу (master_id ищется автоматически по ссылке, чтобы не гадать с цифрами)
-- INSERT INTO services (master_id, name, duration_min, price)
-- SELECT id, 'Мужская стрижка', 45, 1500 FROM masters WHERE invite_link = 'LINK1';

-- -- 3. Создаем слот на СЕГОДНЯ в 15:00 (используем NOW(), чтобы не менять даты руками)
-- -- Смещение +3 часа для МСК, если что поправь под свой часовой пояс
-- INSERT INTO slots (master_id, start_time, end_time, status)
-- SELECT id, NOW() + INTERVAL '3 hours', NOW() + INTERVAL '3 hours 45 minutes', 'free'
-- FROM masters WHERE invite_link = 'LINK1';

-- -- 4. Создаем тестовых клиентов (чтобы обойти Foreign Key при бронировании)
-- INSERT INTO clients (telegram_id, first_name) VALUES (777111222, 'Вася 1') ON CONFLICT DO NOTHING;
-- INSERT INTO clients (telegram_id, first_name) VALUES (888999000, 'Вася 2') ON CONFLICT DO NOTHING;




-- -- 1. Добавляем услуги (ON CONFLICT гарантирует, что при повторном запуске не будет дубликатов)
-- INSERT INTO services (master_id, name, duration_min, price)
-- SELECT id, 'Мужская стрижка', 45, 1500 FROM masters WHERE invite_link = 'LINK1'
-- ON CONFLICT DO NOTHING;

-- INSERT INTO services (master_id, name, duration_min, price)
-- SELECT id, 'Стрижка машинкой', 20, 800 FROM masters WHERE invite_link = 'LINK1'
-- ON CONFLICT DO NOTHING;

-- INSERT INTO services (master_id, name, duration_min, price)
-- SELECT id, 'Оформление бороды', 30, 1000 FROM masters WHERE invite_link = 'LINK1'
-- ON CONFLICT DO NOTHING;

-- INSERT INTO services (master_id, name, duration_min, price)
-- SELECT id, 'Детская стрижка', 30, 700 FROM masters WHERE invite_link = 'LINK1'
-- ON CONFLICT DO NOTHING;


-- -- 2. Генерируем слоты на СЕГОДНЯ и ЗАВТРА (по 9 слотов в день, с 10:00 до 18:00)
-- -- Статус выбирается случайно: 60% - free, 20% - pending, 20% - booked
-- INSERT INTO slots (master_id, start_time, end_time, status)
-- SELECT 
--     m.id, 
--     CURRENT_DATE + d AS day_start, 
--     CURRENT_DATE + d + INTERVAL '45 minutes' AS day_end,
--     CASE 
--         WHEN random() < 0.6 THEN 'free' 
--         WHEN random() < 0.5 THEN 'pending' 
--         ELSE 'booked' 
--     END AS status
-- FROM masters m
-- CROSS JOIN (VALUES (0), (1)) AS days(d) 
-- CROSS JOIN generate_series(10, 18) AS hour(h) 
-- WHERE m.invite_link = 'LINK1'
-- ON CONFLICT DO NOTHING;




-- ALTER TABLE bookings ADD COLUMN IF NOT EXISTS master_id INTEGER REFERENCES masters(id);
-- ALTER TABLE bookings ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();


-- -- Добавляем недостающие поля
-- ALTER TABLE waitlist 
-- ADD COLUMN IF NOT EXISTS client_name VARCHAR(255),
-- ADD COLUMN IF NOT EXISTS client_phone VARCHAR(50),
-- ADD COLUMN IF NOT EXISTS service_name VARCHAR(255),
-- ADD COLUMN IF NOT EXISTS desired_time TIME;

-- -- Заполнить master_id у всех существующих записей
-- UPDATE bookings b
-- SET master_id = s.master_id
-- FROM slots s
-- WHERE b.slot_id = s.id 
--   AND b.master_id IS NULL;

-- -- Проверить результат
-- SELECT id, master_id, status FROM bookings WHERE master_id IS NULL;


-- -- Добавляем колонку is_deleted
-- ALTER TABLE services 
-- ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT FALSE;

-- -- Устанавливаем значение FALSE для всех существующих записей
-- UPDATE services SET is_deleted = FALSE WHERE is_deleted IS NULL;


-- CREATE TABLE IF NOT EXISTS master_photos (
--     id BIGSERIAL PRIMARY KEY,
--     master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
--     url TEXT NOT NULL,
--     created_at TIMESTAMPTZ DEFAULT NOW()
-- );

-- 1. Мастера
CREATE TABLE IF NOT EXISTS masters (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    address TEXT,
    invite_link VARCHAR(50) UNIQUE NOT NULL,
    work_hours JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Услуги (с учетом is_deleted)
CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    duration_min INT NOT NULL,
    price INT NOT NULL,
    description TEXT DEFAULT '',
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Слоты
CREATE TABLE IF NOT EXISTS slots (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'free',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT check_slot_times CHECK (end_time > start_time)
);

-- 4. Клиенты
CREATE TABLE IF NOT EXISTS clients (
    telegram_id BIGINT PRIMARY KEY,
    username VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) DEFAULT '',
    strikes_count INT NOT NULL DEFAULT 0,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Записи (с учетом master_id)
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    slot_id BIGINT NOT NULL UNIQUE REFERENCES slots(id) ON DELETE RESTRICT,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
    client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE RESTRICT,
    client_name VARCHAR(255) NOT NULL,
    price_locked INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Лист ожидания (с учетом новых полей)
CREATE TABLE IF NOT EXISTS waitlist (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE CASCADE,
    client_name VARCHAR(255),
    client_phone VARCHAR(50),
    service_name VARCHAR(255),
    desired_date DATE NOT NULL,
    desired_time TIME,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting',
    offered_slot_id BIGINT REFERENCES slots(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Фото работ
CREATE TABLE IF NOT EXISTS master_photos (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);