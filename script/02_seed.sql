
-- !This seed script is ai-generated
-- ============================================================
-- 02_seed.sql  –  Deterministic seed data
-- Fixed seed via setseed() ensures reproducible results.
-- ============================================================

-- Truncate all tables before seeding (respects FK constraints)
TRUNCATE TABLE user_watch_history, content, users RESTART IDENTITY CASCADE;

-- Fix the random seed for deterministic generation
SELECT setseed(0.42);

-- ============================================================
-- 1. Users  (25 records)
--    - 3 subscription tiers: free, basic, premium
--    - 10 countries (ISO 3166-1 alpha-2)
--    - Ages 18-65
-- ============================================================
INSERT INTO users (age, country, subscription_type, created_at) VALUES
(25, 'US', 'free',    '2024-01-15 08:30:00'),
(34, 'US', 'basic',   '2024-01-20 10:15:00'),
(19, 'US', 'premium', '2024-02-01 14:00:00'),
(45, 'GB', 'free',    '2024-02-05 09:00:00'),
(28, 'GB', 'basic',   '2024-02-10 11:30:00'),
(52, 'GB', 'premium', '2024-02-15 16:45:00'),
(31, 'DE', 'free',    '2024-02-20 07:00:00'),
(22, 'DE', 'basic',   '2024-03-01 13:20:00'),
(40, 'JP', 'premium', '2024-03-05 20:00:00'),
(27, 'JP', 'free',    '2024-03-10 18:10:00'),
(36, 'JP', 'basic',   '2024-03-15 12:00:00'),
(18, 'BR', 'premium', '2024-03-20 09:30:00'),
(55, 'BR', 'free',    '2024-03-25 15:00:00'),
(29, 'IN', 'basic',   '2024-04-01 08:00:00'),
(33, 'IN', 'premium', '2024-04-05 10:45:00'),
(41, 'IN', 'free',    '2024-04-10 14:30:00'),
(23, 'FR', 'basic',   '2024-04-15 11:00:00'),
(48, 'FR', 'premium', '2024-04-20 17:00:00'),
(20, 'KR', 'free',    '2024-04-25 09:15:00'),
(38, 'KR', 'basic',   '2024-05-01 13:00:00'),
(60, 'AU', 'premium', '2024-05-05 10:00:00'),
(26, 'AU', 'free',    '2024-05-10 16:30:00'),
(44, 'CA', 'basic',   '2024-05-15 08:45:00'),
(30, 'CA', 'premium', '2024-05-20 12:00:00'),
(35, 'CA', 'free',    '2024-05-25 14:00:00');

-- ============================================================
-- 2. Content  (50 records)
--    - 6 genres: action, drama, comedy, thriller, documentary, sci-fi
--    - popularity_score in range 0.0 – 1.0 (power-law / long-tail distribution)
--      A few items score very high; most score low-to-moderate.
-- ============================================================
INSERT INTO content (title, genre, popularity_score, created_at, available_countries,available_subscription) VALUES
-- Action (10)
('The Last Stand',            'action',      0.95, '2024-01-01 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'KR', 'AU', 'CA'], ARRAY['basic', 'premium']),
('Fury Road Redux',           'action',      0.89, '2024-01-02 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'AU'], ARRAY['premium']),
('Code Zero',                 'action',      0.72, '2024-01-03 00:00:00', ARRAY['US', 'GB', 'JP', 'KR', 'CA', 'IN'], ARRAY['basic', 'premium']),
('Rapid Strike',              'action',      0.54, '2024-01-04 00:00:00', ARRAY['US', 'DE', 'BR', 'FR', 'AU'], ARRAY['free', 'basic', 'premium']),
('Night Patrol',              'action',      0.36, '2024-01-05 00:00:00', ARRAY['US', 'GB', 'CA', 'KR'], ARRAY['free', 'basic', 'premium']),
('Steel Thunder',             'action',      0.22, '2024-01-06 00:00:00', ARRAY['GB', 'DE', 'JP'], ARRAY['free', 'basic', 'premium']),
('Blaze Runner',              'action',      0.15, '2024-01-07 00:00:00', ARRAY['US', 'IN', 'BR'], ARRAY['basic', 'premium']),
('Shadow Ops',                'action',      0.09, '2024-01-08 00:00:00', ARRAY['FR', 'KR'], ARRAY['premium']),
('Ground Force',              'action',      0.04, '2024-01-09 00:00:00', ARRAY['AU', 'CA'], ARRAY['free', 'basic', 'premium']),
('Rogue Agent',               'action',      0.02, '2024-01-10 00:00:00', ARRAY['DE'], ARRAY['free', 'basic', 'premium']),

-- Drama (10)
('The Quiet Mind',            'drama',       0.98, '2024-01-11 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'KR', 'AU', 'CA'], ARRAY['basic', 'premium']),
('Broken Bridges',            'drama',       0.81, '2024-01-12 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'IN', 'FR', 'CA'], ARRAY['premium']),
('Echoes of Yesterday',       'drama',       0.64, '2024-01-13 00:00:00', ARRAY['US', 'GB', 'BR', 'KR', 'AU'], ARRAY['basic', 'premium']),
('Still Waters',              'drama',       0.47, '2024-01-14 00:00:00', ARRAY['DE', 'JP', 'FR', 'CA'], ARRAY['free', 'basic', 'premium']),
('The Long Road Home',        'drama',       0.31, '2024-01-15 00:00:00', ARRAY['US', 'IN', 'AU'], ARRAY['free', 'basic', 'premium']),
('Paper Walls',               'drama',       0.20, '2024-01-16 00:00:00', ARRAY['GB', 'KR', 'BR'], ARRAY['basic', 'premium']),
('Twilight Conversations',    'drama',       0.11, '2024-01-17 00:00:00', ARRAY['US', 'FR'], ARRAY['free', 'basic', 'premium']),
('Fading Colors',             'drama',       0.06, '2024-01-18 00:00:00', ARRAY['JP', 'CA'], ARRAY['premium']),
('Letters Unsent',            'drama',       0.03, '2024-01-19 00:00:00', ARRAY['IN'], ARRAY['free', 'basic', 'premium']),
('Silent Promise',            'drama',       0.01, '2024-01-20 00:00:00', ARRAY['GB'], ARRAY['free', 'basic', 'premium']),

-- Comedy (10)
('Laugh Factory',             'comedy',      0.93, '2024-02-01 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'KR', 'CA'], ARRAY['basic', 'premium']),
('Oops I Did It Again',       'comedy',      0.78, '2024-02-02 00:00:00', ARRAY['US', 'GB', 'DE', 'BR', 'AU', 'KR'], ARRAY['premium']),
('Sitcom Life',               'comedy',      0.60, '2024-02-03 00:00:00', ARRAY['US', 'JP', 'IN', 'FR', 'CA'], ARRAY['free', 'basic', 'premium']),
('Prank Masters',             'comedy',      0.41, '2024-02-04 00:00:00', ARRAY['GB', 'DE', 'KR', 'AU'], ARRAY['basic', 'premium']),
('The Funny Side',            'comedy',      0.29, '2024-02-05 00:00:00', ARRAY['US', 'BR', 'CA'], ARRAY['free', 'basic', 'premium']),
('Comedy Nights',             'comedy',      0.16, '2024-02-06 00:00:00', ARRAY['IN', 'FR', 'JP'], ARRAY['free', 'basic', 'premium']),
('Blooper Reel',              'comedy',      0.10, '2024-02-07 00:00:00', ARRAY['GB', 'AU'], ARRAY['basic', 'premium']),
('Two and a Half Jokes',      'comedy',      0.05, '2024-02-08 00:00:00', ARRAY['US', 'DE'], ARRAY['free', 'basic', 'premium']),
('Stand Up Stories',          'comedy',      0.02, '2024-02-09 00:00:00', ARRAY['KR'], ARRAY['free', 'basic', 'premium']),
('The Gag Show',              'comedy',      0.01, '2024-02-10 00:00:00', ARRAY['BR'], ARRAY['free', 'basic', 'premium']),

-- Thriller (8)
('Dark Corridor',             'thriller',    0.91, '2024-03-01 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'KR', 'AU'], ARRAY['premium']),
('The Vanishing Point',       'thriller',    0.75, '2024-03-02 00:00:00', ARRAY['US', 'GB', 'DE', 'IN', 'CA', 'AU'], ARRAY['basic', 'premium']),
('Mind Games',                'thriller',    0.55, '2024-03-03 00:00:00', ARRAY['US', 'JP', 'FR', 'KR', 'BR'], ARRAY['basic', 'premium']),
('Whisper Network',           'thriller',    0.38, '2024-03-04 00:00:00', ARRAY['GB', 'DE', 'CA', 'AU'], ARRAY['free', 'basic', 'premium']),
('Cold Trace',                'thriller',    0.22, '2024-03-05 00:00:00', ARRAY['US', 'IN', 'JP'], ARRAY['free', 'basic', 'premium']),
('The Watcher',               'thriller',    0.12, '2024-03-06 00:00:00', ARRAY['FR', 'KR'], ARRAY['premium']),
('Nerve Endings',             'thriller',    0.06, '2024-03-07 00:00:00', ARRAY['GB', 'BR'], ARRAY['free', 'basic', 'premium']),
('Paranoia',                  'thriller',    0.02, '2024-03-08 00:00:00', ARRAY['CA'], ARRAY['free', 'basic', 'premium']),

-- Documentary (7)
('Planet Earth Revisited',    'documentary', 0.89, '2024-04-01 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'IN', 'FR', 'AU', 'CA'], ARRAY['free', 'basic', 'premium']),
('Into the Deep',             'documentary', 0.67, '2024-04-02 00:00:00', ARRAY['US', 'GB', 'BR', 'KR', 'CA'], ARRAY['basic', 'premium']),
('The Human Story',           'documentary', 0.45, '2024-04-03 00:00:00', ARRAY['DE', 'JP', 'IN', 'FR'], ARRAY['free', 'basic', 'premium']),
('Tech Revolution',           'documentary', 0.27, '2024-04-04 00:00:00', ARRAY['US', 'AU', 'KR'], ARRAY['basic', 'premium']),
('History Untold',            'documentary', 0.13, '2024-04-05 00:00:00', ARRAY['GB', 'CA'], ARRAY['free', 'basic', 'premium']),
('Wild Frontiers',            'documentary', 0.07, '2024-04-06 00:00:00', ARRAY['DE', 'BR'], ARRAY['free', 'basic', 'premium']),
('Ocean Secrets',             'documentary', 0.02, '2024-04-07 00:00:00', ARRAY['JP'], ARRAY['free', 'basic', 'premium']),

-- Sci-Fi (5)
('Galaxy Drift',              'sci-fi',      0.94, '2024-05-01 00:00:00', ARRAY['US', 'GB', 'DE', 'JP', 'BR', 'IN', 'FR', 'KR', 'AU', 'CA'], ARRAY['premium']),
('Quantum Paradox',           'sci-fi',      0.69, '2024-05-02 00:00:00', ARRAY['US', 'GB', 'JP', 'IN', 'AU', 'CA'], ARRAY['basic', 'premium']),
('Neon Horizon',              'sci-fi',      0.39, '2024-05-03 00:00:00', ARRAY['DE', 'FR', 'KR', 'BR'], ARRAY['free', 'basic', 'premium']),
('The Singularity',           'sci-fi',      0.18, '2024-05-04 00:00:00', ARRAY['US', 'GB', 'JP'], ARRAY['basic', 'premium']),
('Starfield',                 'sci-fi',      0.05, '2024-05-05 00:00:00', ARRAY['IN', 'AU'], ARRAY['free', 'basic', 'premium']);

-- ============================================================
-- 3. Watch History  (200+ records)
--    Deterministically generated using the seeded random().
--    Each row: random user (1-25) x random content (1-50)
--    watched_at spread across 2024.
-- ============================================================
INSERT INTO user_watch_history (user_id, content_id, watched_at)
SELECT
    -- deterministic user_id 1..25
    (floor(random() * 25) + 1)::BIGINT AS user_id,
    -- skew toward popular content (lower IDs = higher popularity within each genre block)
    -- using power-law: random()^2 biases toward 0, so id biases toward 1
    (floor(power(random(), 2) * 50) + 1)::BIGINT AS content_id,
    -- spread across 2024
    '2024-01-01'::TIMESTAMP + (random() * 365) * INTERVAL '1 day'
        + (random() * 24) * INTERVAL '1 hour'
        + (random() * 60) * INTERVAL '1 minute'
    AS watched_at
FROM generate_series(1, 220);

-- ============================================================
-- Verification queries (output visible in docker logs)
-- ============================================================
DO $$
DECLARE
    u_count   BIGINT;
    c_count   BIGINT;
    wh_count  BIGINT;
    g_count   BIGINT;
BEGIN
    SELECT count(*) INTO u_count  FROM users;
    SELECT count(*) INTO c_count  FROM content;
    SELECT count(*) INTO wh_count FROM user_watch_history;
    SELECT count(DISTINCT genre) INTO g_count FROM content;

    RAISE NOTICE '=== Seed Verification ===';
    RAISE NOTICE 'Users:                % (min 20)',  u_count;
    RAISE NOTICE 'Content items:        % (min 50)',  c_count;
    RAISE NOTICE 'Watch history records: % (min 200)', wh_count;
    RAISE NOTICE 'Distinct genres:      % (min 5)',   g_count;

    -- Hard assertions
    IF u_count  < 20  THEN RAISE EXCEPTION 'Not enough users: %', u_count;  END IF;
    IF c_count  < 50  THEN RAISE EXCEPTION 'Not enough content: %', c_count; END IF;
    IF wh_count < 200 THEN RAISE EXCEPTION 'Not enough watch history: %', wh_count; END IF;
    IF g_count  < 5   THEN RAISE EXCEPTION 'Not enough genres: %', g_count;  END IF;

    RAISE NOTICE '=== All checks passed ===';
END $$;
