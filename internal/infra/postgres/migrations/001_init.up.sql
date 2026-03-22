CREATE SCHEMA IF NOT EXISTS vacancies;

SET search_path TO vacancies;

CREATE TABLE IF NOT EXISTS activity_types (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS skills (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS vacancies (
    id                UUID PRIMARY KEY,
    employer_id       UUID         NOT NULL,
    is_archived       BOOLEAN      NOT NULL DEFAULT FALSE,
    title             VARCHAR(255) NOT NULL,
    activity_type_id  INT          NOT NULL REFERENCES activity_types(id),
    employment_type   VARCHAR(50)  NOT NULL,
    work_schedule     VARCHAR(50)  NOT NULL,
    is_verified       BOOLEAN      NOT NULL DEFAULT FALSE,
    compensation_type VARCHAR(50)  NOT NULL DEFAULT 'Без вознаграждения',
    compensation_min  NUMERIC(12,2) NOT NULL DEFAULT 0,
    compensation_max  NUMERIC(12,2) NOT NULL DEFAULT 0,
    description       TEXT         NOT NULL DEFAULT '',
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vacancy_skills (
    vacancy_id UUID NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
    skill_id   INT  NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (vacancy_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_vacancies_employer    ON vacancies(employer_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_archived    ON vacancies(is_archived);
CREATE INDEX IF NOT EXISTS idx_vacancies_created     ON vacancies(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_vacancy_skills_skill  ON vacancy_skills(skill_id);

-- Seed some activity types
INSERT INTO activity_types (name) VALUES
    ('Программирование'),
    ('Дизайн'),
    ('Маркетинг'),
    ('Аналитика'),
    ('Менеджмент'),
    ('Наука'),
    ('Образование')
ON CONFLICT (name) DO NOTHING;
