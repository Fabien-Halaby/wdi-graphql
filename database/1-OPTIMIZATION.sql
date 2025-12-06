-- ============================================
-- SCRIPT D'OPTIMISATION POSTGRESQL WDI
-- À exécuter sur votre base existante
-- ============================================

-- 1. AJOUTER LES CLÉS PRIMAIRES ET ÉTRANGÈRES
-- ============================================

-- Vérifier et ajouter les clés primaires si manquantes
DO $$ 
BEGIN
    -- Clé primaire sur country
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'country_pkey'
    ) THEN
        ALTER TABLE country ADD PRIMARY KEY (countrycode);
    END IF;

    -- Clé primaire sur series
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'series_pkey'
    ) THEN
        ALTER TABLE series ADD PRIMARY KEY (seriescode);
    END IF;

    -- Ajouter ID auto-incrémenté sur indicators si pas présent
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'indicators' AND column_name = 'id'
    ) THEN
        ALTER TABLE indicators ADD COLUMN id BIGSERIAL PRIMARY KEY;
    END IF;
END $$;

-- Ajouter contrainte d'unicité sur indicators
ALTER TABLE indicators 
DROP CONSTRAINT IF EXISTS unique_indicator;

ALTER TABLE indicators 
ADD CONSTRAINT unique_indicator 
UNIQUE(countrycode, indicatorcode, year);

-- Ajouter les clés étrangères
ALTER TABLE indicators 
DROP CONSTRAINT IF EXISTS fk_indicators_country;

ALTER TABLE indicators 
ADD CONSTRAINT fk_indicators_country 
FOREIGN KEY (countrycode) REFERENCES country(countrycode) ON DELETE CASCADE;

ALTER TABLE indicators 
DROP CONSTRAINT IF EXISTS fk_indicators_series;

ALTER TABLE indicators 
ADD CONSTRAINT fk_indicators_series 
FOREIGN KEY (indicatorcode) REFERENCES series(seriescode) ON DELETE CASCADE;

-- 2. CRÉER DES INDEX OPTIMISÉS
-- ============================================

-- Index composites pour requêtes fréquentes (TRÈS IMPORTANT)
DROP INDEX IF EXISTS idx_indicators_country_year;
CREATE INDEX idx_indicators_country_year 
ON indicators(countrycode, year) 
INCLUDE (value, indicatorcode, indicatorname);

DROP INDEX IF EXISTS idx_indicators_indicator_year;
CREATE INDEX idx_indicators_indicator_year 
ON indicators(indicatorcode, year) 
INCLUDE (value, countrycode, countryname);

DROP INDEX IF EXISTS idx_indicators_year_value;
CREATE INDEX idx_indicators_year_value 
ON indicators(year, value) 
WHERE value IS NOT NULL;

-- Index pour recherches par valeur
DROP INDEX IF EXISTS idx_indicators_value_range;
CREATE INDEX idx_indicators_value_range 
ON indicators(value) 
WHERE value IS NOT NULL;

-- Index pour jointures optimisées
DROP INDEX IF EXISTS idx_indicators_full_composite;
CREATE INDEX idx_indicators_full_composite 
ON indicators(countrycode, indicatorcode, year, value);

-- Index sur country pour jointures
DROP INDEX IF EXISTS idx_country_region_income;
CREATE INDEX idx_country_region_income 
ON country(region, incomegroup);

DROP INDEX IF EXISTS idx_country_shortname;
CREATE INDEX idx_country_shortname 
ON country(shortname);

-- Index sur series pour recherche
DROP INDEX IF EXISTS idx_series_topic;
CREATE INDEX idx_series_topic 
ON series(topic);

-- Index full-text pour recherche d'indicateurs
DROP INDEX IF EXISTS idx_series_fulltext;
CREATE INDEX idx_series_fulltext 
ON series USING GIN(to_tsvector('english', 
    COALESCE(indicatorname, '') || ' ' || 
    COALESCE(shortdefinition, '') || ' ' || 
    COALESCE(topic, '')
));

-- Index pour countrynotes, seriesnotes, footnotes
DROP INDEX IF EXISTS idx_countrynotes_country;
CREATE INDEX idx_countrynotes_country ON countrynotes(countrycode);

DROP INDEX IF EXISTS idx_countrynotes_series;
CREATE INDEX idx_countrynotes_series ON countrynotes(seriescode);

DROP INDEX IF EXISTS idx_seriesnotes_series;
CREATE INDEX idx_seriesnotes_series ON seriesnotes(seriescode);

DROP INDEX IF EXISTS idx_footnotes_country_series;
CREATE INDEX idx_footnotes_country_series ON footnotes(countrycode, seriescode);

-- 3. CRÉER DES VUES MATÉRIALISÉES POUR AGRÉGATIONS
-- ============================================

-- Vue: Statistiques par pays et année
DROP MATERIALIZED VIEW IF EXISTS mv_country_year_stats CASCADE;
CREATE MATERIALIZED VIEW mv_country_year_stats AS
SELECT 
    i.countrycode,
    i.countryname,
    c.region,
    c.incomegroup,
    i.year,
    COUNT(*) as total_indicators,
    COUNT(i.value) as indicators_with_data,
    COUNT(DISTINCT i.indicatorcode) as unique_indicators
FROM indicators i
LEFT JOIN country c ON i.countrycode = c.countrycode
GROUP BY i.countrycode, i.countryname, c.region, c.incomegroup, i.year;

CREATE UNIQUE INDEX idx_mv_country_year_unique 
ON mv_country_year_stats(countrycode, year);

CREATE INDEX idx_mv_country_year_region 
ON mv_country_year_stats(region, year);

-- Vue: Dernières valeurs par indicateur et pays
DROP MATERIALIZED VIEW IF EXISTS mv_latest_indicators CASCADE;
CREATE MATERIALIZED VIEW mv_latest_indicators AS
SELECT DISTINCT ON (countrycode, indicatorcode)
    countrycode,
    countryname,
    indicatorcode,
    indicatorname,
    year as latest_year,
    value as latest_value
FROM indicators
WHERE value IS NOT NULL
ORDER BY countrycode, indicatorcode, year DESC;

CREATE UNIQUE INDEX idx_mv_latest_unique 
ON mv_latest_indicators(countrycode, indicatorcode);

CREATE INDEX idx_mv_latest_indicator 
ON mv_latest_indicators(indicatorcode);

-- Vue: Statistiques globales par indicateur
DROP MATERIALIZED VIEW IF EXISTS mv_indicator_stats CASCADE;
CREATE MATERIALIZED VIEW mv_indicator_stats AS
SELECT 
    indicatorcode,
    MIN(indicatorname) as indicatorname,
    COUNT(*) as total_records,
    COUNT(DISTINCT countrycode) as country_count,
    MIN(year) as first_year,
    MAX(year) as last_year,
    COUNT(DISTINCT year) as year_count,
    MIN(value) as min_value,
    MAX(value) as max_value,
    AVG(value) as avg_value,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY value) as median_value,
    STDDEV(value) as stddev_value
FROM indicators
WHERE value IS NOT NULL
GROUP BY indicatorcode;

CREATE UNIQUE INDEX idx_mv_indicator_stats_unique 
ON mv_indicator_stats(indicatorcode);

-- Vue: Top indicateurs par région
DROP MATERIALIZED VIEW IF EXISTS mv_region_indicators CASCADE;
CREATE MATERIALIZED VIEW mv_region_indicators AS
SELECT 
    c.region,
    i.indicatorcode,
    MIN(i.indicatorname) as indicatorname,
    i.year,
    AVG(i.value) as avg_value,
    COUNT(*) as country_count
FROM indicators i
JOIN country c ON i.countrycode = c.countrycode
WHERE i.value IS NOT NULL AND c.region IS NOT NULL
GROUP BY c.region, i.indicatorcode, i.year;

CREATE INDEX idx_mv_region_indicators 
ON mv_region_indicators(region, indicatorcode, year);

-- Vue: Tendances temporelles (évolution année par année)
DROP MATERIALIZED VIEW IF EXISTS mv_yearly_trends CASCADE;
CREATE MATERIALIZED VIEW mv_yearly_trends AS
SELECT 
    indicatorcode,
    MIN(indicatorname) as indicatorname,
    year,
    COUNT(DISTINCT countrycode) as country_count,
    AVG(value) as global_avg,
    MIN(value) as global_min,
    MAX(value) as global_max
FROM indicators
WHERE value IS NOT NULL
GROUP BY indicatorcode, year;

CREATE INDEX idx_mv_yearly_trends 
ON mv_yearly_trends(indicatorcode, year);

-- 4. FONCTIONS OPTIMISÉES POUR REQUÊTES FRÉQUENTES
-- ============================================

-- Fonction: Obtenir les indicateurs d'un pays avec cache
CREATE OR REPLACE FUNCTION get_country_indicators(
    p_country_code VARCHAR,
    p_year_start INTEGER DEFAULT NULL,
    p_year_end INTEGER DEFAULT NULL
)
RETURNS TABLE(
    indicator_code VARCHAR,
    indicator_name TEXT,
    year INTEGER,
    value NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.indicatorcode,
        i.indicatorname,
        i.year,
        i.value
    FROM indicators i
    WHERE i.countrycode = p_country_code
        AND (p_year_start IS NULL OR i.year >= p_year_start)
        AND (p_year_end IS NULL OR i.year <= p_year_end)
        AND i.value IS NOT NULL
    ORDER BY i.year DESC, i.indicatorname;
END;
$$ LANGUAGE plpgsql STABLE;

-- Fonction: Comparer plusieurs pays sur un indicateur
CREATE OR REPLACE FUNCTION compare_countries(
    p_indicator_code VARCHAR,
    p_country_codes VARCHAR[],
    p_year INTEGER
)
RETURNS TABLE(
    country_code VARCHAR,
    country_name TEXT,
    value NUMERIC,
    rank INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.countrycode,
        i.countryname,
        i.value,
        RANK() OVER (ORDER BY i.value DESC)::INTEGER
    FROM indicators i
    WHERE i.indicatorcode = p_indicator_code
        AND i.year = p_year
        AND i.countrycode = ANY(p_country_codes)
        AND i.value IS NOT NULL
    ORDER BY i.value DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- Fonction: Top N pays pour un indicateur
CREATE OR REPLACE FUNCTION get_top_countries(
    p_indicator_code VARCHAR,
    p_year INTEGER,
    p_limit INTEGER DEFAULT 10,
    p_ascending BOOLEAN DEFAULT FALSE
)
RETURNS TABLE(
    country_code VARCHAR,
    country_name TEXT,
    region VARCHAR,
    value NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.countrycode,
        i.countryname,
        c.region,
        i.value
    FROM indicators i
    JOIN country c ON i.countrycode = c.countrycode
    WHERE i.indicatorcode = p_indicator_code
        AND i.year = p_year
        AND i.value IS NOT NULL
    ORDER BY 
        CASE WHEN p_ascending THEN i.value END ASC,
        CASE WHEN NOT p_ascending THEN i.value END DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql STABLE;

-- Fonction: Évolution temporelle d'un indicateur pour un pays
CREATE OR REPLACE FUNCTION get_indicator_timeline(
    p_country_code VARCHAR,
    p_indicator_code VARCHAR
)
RETURNS TABLE(
    year INTEGER,
    value NUMERIC,
    year_over_year_change NUMERIC,
    percent_change NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.year,
        i.value,
        i.value - LAG(i.value) OVER (ORDER BY i.year) as year_over_year_change,
        CASE 
            WHEN LAG(i.value) OVER (ORDER BY i.year) > 0 
            THEN ((i.value - LAG(i.value) OVER (ORDER BY i.year)) / LAG(i.value) OVER (ORDER BY i.year)) * 100
            ELSE NULL
        END as percent_change
    FROM indicators i
    WHERE i.countrycode = p_country_code
        AND i.indicatorcode = p_indicator_code
        AND i.value IS NOT NULL
    ORDER BY i.year;
END;
$$ LANGUAGE plpgsql STABLE;

-- 5. OPTIMISATIONS DE CONFIGURATION
-- ============================================

-- Augmenter les statistiques pour de meilleures estimations
ALTER TABLE indicators ALTER COLUMN countrycode SET STATISTICS 1000;
ALTER TABLE indicators ALTER COLUMN indicatorcode SET STATISTICS 1000;
ALTER TABLE indicators ALTER COLUMN year SET STATISTICS 1000;
ALTER TABLE indicators ALTER COLUMN value SET STATISTICS 1000;

-- 6. ANALYSER ET VACUUM
-- ============================================

-- Nettoyer et analyser toutes les tables
VACUUM (ANALYZE, VERBOSE) country;
VACUUM (ANALYZE, VERBOSE) series;
VACUUM (ANALYZE, VERBOSE) indicators;
VACUUM (ANALYZE, VERBOSE) countrynotes;
VACUUM (ANALYZE, VERBOSE) seriesnotes;
VACUUM (ANALYZE, VERBOSE) footnotes;

-- Analyser les vues matérialisées
ANALYZE mv_country_year_stats;
ANALYZE mv_latest_indicators;
ANALYZE mv_indicator_stats;
ANALYZE mv_region_indicators;
ANALYZE mv_yearly_trends;

-- 7. RAPPORT D'OPTIMISATION
-- ============================================

SELECT '=== STATISTIQUES DE LA BASE ===' as info;

SELECT 
    'Total pays' as metric,
    COUNT(*)::text as value
FROM country
UNION ALL
SELECT 
    'Total indicateurs',
    COUNT(*)::text
FROM indicators
UNION ALL
SELECT 
    'Années couvertes',
    MIN(year)::text || ' - ' || MAX(year)::text
FROM indicators
UNION ALL
SELECT 
    'Indicateurs uniques',
    COUNT(DISTINCT indicatorcode)::text
FROM indicators
UNION ALL
SELECT
    'Pays uniques',
    COUNT(DISTINCT countrycode)::text
FROM indicators;

SELECT '=== TAILLE DES TABLES ===' as info;

SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - 
                   pg_relation_size(schemaname||'.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

SELECT '=== INDEX CRÉÉS ===' as info;

SELECT 
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_indexes
JOIN pg_stat_user_indexes USING (schemaname, tablename, indexname)
WHERE schemaname = 'public'
ORDER BY tablename, indexname;