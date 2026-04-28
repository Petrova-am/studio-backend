CREATE TABLE IF NOT EXISTS trainers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    specialty VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO trainers (name, specialty) VALUES
    ('Анна Петрова', 'Основатель студии, инструктор по растяжке'),
    ('Елена Смирнова', 'Инструктор по йоге и пилатесу'),
    ('Мария Иванова', 'Инструктор по функциональному тренингу');