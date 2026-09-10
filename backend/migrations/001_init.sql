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