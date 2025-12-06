-- ============================================
-- TESTS DE PERFORMANCE WDI
-- Exécutez ces requêtes pour mesurer les performances
-- ============================================

-- Active l'affichage des temps d'exécution
\timing on

-- ===========================================
-- 1. TESTS DE REQUÊTES BASIQUES
-- ===========================================

-- CORRECTION ERREUR LIGNE 12 : Utilisation de deux apostrophes ('' ) pour échapper l'apostrophe simple (d'un)
\echo '=== TEST 1: Récupérer tous les indicateurs d''un pays (USA, 2010-2020) ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT year, indicatorname, value
FROM indicators
WHERE countrycode = 'USA' 
    AND year BETWEEN 2010 AND 2020
ORDER BY year, indicatorname;

\echo '=== TEST 2: Top 10 pays par PIB en 2020 ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT countryname, value
FROM indicators
WHERE indicatorcode = 'NY.GDP.MKTP.CD' 
    AND year = 2020
    AND value IS NOT NULL
ORDER BY value DESC 
LIMIT 10;

\echo '=== TEST 3: Évolution d''un indicateur sur le temps ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT year, AVG(value) as avg_value, COUNT(*) as country_count
FROM indicators
WHERE indicatorcode = 'SP.POP.TOTL'
    AND year BETWEEN 1990 AND 2020
    AND value IS NOT NULL
GROUP BY year
ORDER BY year;

\echo '=== TEST 4: Comparaison régionale ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT 
    c.region,
    i.year,
    AVG(i.value) as avg_gdp,
    COUNT(DISTINCT i.countrycode) as country_count
FROM indicators i
JOIN country c ON i.countrycode = c.countrycode
WHERE i.indicatorcode = 'NY.GDP.MKTP.CD'
    AND i.year >= 2010
    AND i.value IS NOT NULL
    AND c.region IS NOT NULL
GROUP BY c.region, i.year
ORDER BY c.region, i.year;

\echo '=== TEST 5: Recherche full-text d''indicateurs ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT seriescode, indicatorname, topic
FROM series
WHERE to_tsvector('english', indicatorname || ' ' || COALESCE(shortdefinition, '')) 
    @@ to_tsquery('english', 'GDP & growth')
LIMIT 20;

\echo '=== TEST 6: Utilisation de la vue matérialisée (dernières valeurs) ==='
-- NOTE: Assurez-vous que la vue mv_latest_indicators existe.
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT countryname, indicatorname, latest_year, latest_value
FROM mv_latest_indicators
WHERE indicatorcode = 'NY.GDP.MKTP.CD'
ORDER BY latest_value DESC
LIMIT 10;

\echo '=== TEST 7: Jointure complexe avec agrégation ==='
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT 
    c.shortname,
    c.region,
    c.incomegroup,
    COUNT(DISTINCT i.indicatorcode) as indicator_count,
    AVG(CASE WHEN i.indicatorcode = 'SP.POP.TOTL' THEN i.value END) as population,
    AVG(CASE WHEN i.indicatorcode = 'NY.GDP.MKTP.CD' THEN i.value END) as gdp
FROM country c
LEFT JOIN indicators i ON c.countrycode = i.countrycode
WHERE i.year = 2020
GROUP BY c.shortname, c.region, c.incomegroup
HAVING COUNT(DISTINCT i.indicatorcode) > 100
ORDER BY gdp DESC NULLS LAST
LIMIT 20;

-- ===========================================
-- TEST 8: Utilisation des fonctions personnalisées (CORRIGÉES)
-- ===========================================

\echo '=== TEST 8: Utilisation des fonctions personnalisées ==='

-- CORRECTION FONCTION 1 : Ajustement du type countryname à TEXT 
DROP FUNCTION IF EXISTS get_top_countries(varchar, integer, integer, boolean);
CREATE OR REPLACE FUNCTION get_top_countries(
    p_indicator_code VARCHAR,
    p_year INTEGER,
    p_limit INTEGER,
    p_ascending BOOLEAN
)
RETURNS TABLE (
    countrycode CHAR(3), 
    countryname TEXT,         -- CORRECTION: Assurez-vous que le type correspond à la table
    region VARCHAR,           
    value NUMERIC
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        i.countrycode,
        i.countryname, -- La colonne 2 est retournée en tant que TEXT par la requête.
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
$$;
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT * FROM get_top_countries('NY.GDP.MKTP.CD', 2020, 10, false);


-- CORRECTION FONCTION 2 : Ajustement du type year_over_year_change à NUMERIC
DROP FUNCTION IF EXISTS get_indicator_timeline(varchar, varchar);
CREATE OR REPLACE FUNCTION get_indicator_timeline(
    p_country_code VARCHAR,
    p_indicator_code VARCHAR
)
RETURNS TABLE (
    year INTEGER,
    value NUMERIC,
    year_over_year_change NUMERIC, -- CORRECTION: NUMERIC est plus sûr pour les différences
    percent_change NUMERIC
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        i.year,
        i.value,
        (i.value - LAG(i.value) OVER (ORDER BY i.year)) as year_over_year_change,
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
$$;
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT * FROM get_indicator_timeline('USA', 'NY.GDP.MKTP.CD');


-- ===========================================
-- 2. BENCHMARKS DE PERFORMANCE
-- ===========================================

\echo '=== BENCHMARK: Scan séquentiel vs Index ==='

-- Sans index (simulation)
SET enable_indexscan = off;
SET enable_bitmapscan = off;
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT COUNT(*) FROM indicators WHERE year = 2020;

-- Avec index
SET enable_indexscan = on;
SET enable_bitmapscan = on;
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT COUNT(*) FROM indicators WHERE year = 2020;

\echo '=== BENCHMARK: Agrégation avec et sans vue matérialisée ==='

-- Sans vue matérialisée (calcul à la volée)
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT 
    countrycode,
    year,
    COUNT(*) as indicator_count
FROM indicators
WHERE year = 2020
GROUP BY countrycode, year;

-- Avec vue matérialisée
-- NOTE: Assurez-vous que la vue mv_country_year_stats existe.
EXPLAIN (ANALYZE, BUFFERS, TIMING)
SELECT countrycode, year, total_indicators
FROM mv_country_year_stats
WHERE year = 2020;

-- ===========================================
-- 3. ANALYSE DES INDEX UTILISÉS (CORRIGÉ)
-- ===========================================

\echo '=== STATISTIQUES D''UTILISATION DES INDEX ==='
-- CORRECTION: Utilisation de relname au lieu de tablename
SELECT 
    schemaname,
    relname,
    indexname,
    idx_scan as times_used,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

\echo '=== INDEX INUTILISÉS (à supprimer ?) ==='
-- CORRECTION: Utilisation de relname au lieu de tablename
SELECT 
    schemaname,
    relname,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND idx_scan = 0
    AND indexrelid IS NOT NULL
ORDER BY pg_relation_size(indexrelid) DESC;

-- ===========================================
-- 4. CACHE ET MÉMOIRE (CORRIGÉ)
-- ===========================================

\echo '=== STATISTIQUES DE CACHE ==='
-- CORRECTION: Utilisation de relname au lieu de tablename
SELECT 
    schemaname,
    relname,
    heap_blks_read as disk_reads,
    heap_blks_hit as cache_hits,
    CASE 
        WHEN (heap_blks_hit + heap_blks_read) > 0 
        THEN ROUND(100.0 * heap_blks_hit / (heap_blks_hit + heap_blks_read), 2)
        ELSE 0 
    END as cache_hit_ratio
FROM pg_statio_user_tables
WHERE schemaname = 'public'
ORDER BY heap_blks_read DESC;

-- ===========================================
-- 5. REQUÊTES LENTES
-- ===========================================

\echo '=== REQUÊTES ACTUELLEMENT EN COURS (> 1 seconde) ==='

SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    EXTRACT(EPOCH FROM (now() - query_start))::INTEGER as duration_seconds,
    state,
    LEFT(query, 100) as query_preview
FROM pg_stat_activity
WHERE state = 'active'
    AND query NOT LIKE '%pg_stat_activity%'
    AND now() - query_start > interval '1 second'
ORDER BY duration_seconds DESC;

-- ===========================================
-- 6. MAINTENANCE AUTOMATIQUE
-- ===========================================

\echo '=== SCRIPT DE MAINTENANCE QUOTIDIENNE ==='

-- Fonction de maintenance complète
CREATE OR REPLACE FUNCTION maintain_wdi_database()
RETURNS TABLE(
    step TEXT,
    status TEXT,
    duration INTERVAL
) AS $$
DECLARE
    start_time TIMESTAMP;
BEGIN
    -- Rafraîchir les vues matérialisées
    start_time := clock_timestamp();
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_country_year_stats;
    RETURN QUERY SELECT 
        'Refresh mv_country_year_stats'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_latest_indicators;
    RETURN QUERY SELECT 
        'Refresh mv_latest_indicators'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_indicator_stats;
    RETURN QUERY SELECT 
        'Refresh mv_indicator_stats'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_region_indicators;
    RETURN QUERY SELECT 
        'Refresh mv_region_indicators'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_yearly_trends;
    RETURN QUERY SELECT 
        'Refresh mv_yearly_trends'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    -- VACUUM ANALYZE
    start_time := clock_timestamp();
    VACUUM ANALYZE indicators;
    RETURN QUERY SELECT 
        'VACUUM ANALYZE indicators'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    VACUUM ANALYZE country;
    RETURN QUERY SELECT 
        'VACUUM ANALYZE country'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
    
    start_time := clock_timestamp();
    VACUUM ANALYZE series;
    RETURN QUERY SELECT 
        'VACUUM ANALYZE series'::TEXT, 
        'OK'::TEXT, 
        clock_timestamp() - start_time;
END;
$$ LANGUAGE plpgsql;

\echo 'Fonction de maintenance créée. Exécutez: SELECT * FROM maintain_wdi_database();'

-- ===========================================
-- 7. CONFIGURATION OPTIMALE
-- ===========================================

\echo '=== CONFIGURATION POSTGRESQL ACTUELLE ==='

SELECT name, setting, unit, context
FROM pg_settings
WHERE name IN (
    'shared_buffers',
    'effective_cache_size',
    'work_mem',
    'maintenance_work_mem',
    'max_parallel_workers_per_gather',
    'max_parallel_workers',
    'random_page_cost',
    'effective_io_concurrency'
)
ORDER BY name;

\echo '=== RECOMMANDATIONS DE CONFIGURATION ==='
\echo 'Pour de meilleures performances, ajoutez à postgresql.conf:'
\echo 'shared_buffers = 4GB'
\echo 'effective_cache_size = 12GB'
\echo 'work_mem = 256MB'
\echo 'maintenance_work_mem = 1GB'
\echo 'max_parallel_workers_per_gather = 4'
\echo 'max_parallel_workers = 8'
\echo 'random_page_cost = 1.1  # Pour SSD'
\echo 'effective_io_concurrency = 200'

-- ===========================================
-- 8. RAPPORT FINAL
-- ===========================================

\echo '=== RAPPORT FINAL D''OPTIMISATION ==='

SELECT 
    'Base de données' as info,
    current_database() as value
UNION ALL
SELECT 
    'Taille totale',
    pg_size_pretty(pg_database_size(current_database()))
UNION ALL
SELECT
    'Nombre de connexions actives',
    COUNT(*)::text
FROM pg_stat_activity
WHERE state = 'active'
UNION ALL
SELECT
    'Nombre d''index',
    COUNT(*)::text
FROM pg_indexes
WHERE schemaname = 'public'
UNION ALL
SELECT
    'Nombre de vues matérialisées',
    COUNT(*)::text
FROM pg_matviews
WHERE schemaname = 'public';

\echo '✅ Optimisation terminée!'
\echo '📊 Exécutez: SELECT * FROM maintain_wdi_database(); pour maintenance'