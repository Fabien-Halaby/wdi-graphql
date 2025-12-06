-- ============================================
-- Tests de Performance WDI Database
-- ============================================

\timing on

-- Test 1: Liste des pays par région
EXPLAIN ANALYZE
SELECT countrycode, shortname, region
FROM country
WHERE region = 'Sub-Saharan Africa'
LIMIT 10;

-- Test 2: Top 10 pays pour un indicateur (GDP 2020)
EXPLAIN ANALYZE
SELECT i.countrycode, i.indicatorcode, i.year, i.value
FROM indicators i
WHERE i.indicatorcode = 'NY.GDP.MKTP.CD'
  AND i.year = 2020
ORDER BY i.value DESC NULLS LAST
LIMIT 10;

-- Test 3: Moyenne d'un indicateur par pays et année
EXPLAIN ANALYZE
SELECT 
  i.countrycode,
  c.shortname,
  i.year,
  AVG(i.value) as avg_value
FROM indicators i
JOIN country c ON i.countrycode = c.countrycode
WHERE i.indicatorcode = 'SP.URB.TOTL.IN.ZS'
GROUP BY i.countrycode, c.shortname, i.year
ORDER BY i.year DESC
LIMIT 100;

-- Test 4: Agrégation régionale
EXPLAIN ANALYZE
SELECT 
  c.region,
  COUNT(DISTINCT i.countrycode) as country_count,
  AVG(i.value) as avg_indicator_value
FROM indicators i
JOIN country c ON i.countrycode = c.countrycode
WHERE i.year >= 2010 AND i.value IS NOT NULL
GROUP BY c.region;

-- Test 5: Ranking par indicateur et année
EXPLAIN ANALYZE
SELECT 
  indicatorcode,
  year,
  countrycode,
  value,
  RANK() OVER (PARTITION BY indicatorcode, year ORDER BY value DESC NULLS LAST) as rank
FROM indicators
WHERE year = 2020 AND value IS NOT NULL
LIMIT 100;

-- Vérifier les vues matérialisées
SELECT COUNT(*) as mv_country_year_stats_count FROM mv_country_year_stats;
SELECT COUNT(*) as mv_latest_indicators_count FROM mv_latest_indicators;
SELECT COUNT(*) as mv_indicator_stats_count FROM mv_indicator_stats;

\timing off
