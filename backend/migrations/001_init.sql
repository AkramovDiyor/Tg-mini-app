CREATE TABLE masters (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    address TEXT,
    invite_link VARCHAR(50) UNIQUE NOT NULL,
    work_hours JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    duration_min INT NOT NULL,
    price INT NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE slots (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'free',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT check_slot_times CHECK (end_time > start_time)
);

CREATE TABLE clients (
    telegram_id BIGINT PRIMARY KEY,
    username VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) DEFAULT '',
    strikes_count INT NOT NULL DEFAULT 0,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    slot_id BIGINT NOT NULL UNIQUE REFERENCES slots(id) ON DELETE RESTRICT,
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
    client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE RESTRICT,
    client_name VARCHAR(255) NOT NULL, -- Обязательное поле для имени Васи!
    price_locked INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS waitlist (
    id BIGSERIAL PRIMARY KEY,
    master_id BIGINT NOT NULL REFERENCES masters(id) ON DELETE CASCADE,
    client_telegram_id BIGINT NOT NULL REFERENCES clients(telegram_id) ON DELETE CASCADE,
    desired_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting', -- waiting, offered, fulfilled, expired
    offered_slot_id BIGINT REFERENCES slots(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
ALTER TABLE waitlist ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
