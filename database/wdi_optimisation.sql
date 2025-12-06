-- ============================================
-- Optimisation WDI Database
-- ============================================

-- INDEX COMPOSITES pour time-series (countrycode, indicatorcode, year)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_indicators_country_indicator_year
  ON indicators(countrycode, indicatorcode, year);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_indicators_indicator_year
  ON indicators(indicatorcode, year);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_indicators_country_year
  ON indicators(countrycode, year);

-- INDEX sur value pour les rankings
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_indicators_value
  ON indicators(value DESC NULLS LAST)
  WHERE value IS NOT NULL;

-- INDEX sur les tables de lookup
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_country_code
  ON country(countrycode);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_series_code
  ON series(seriescode);

-- VUES MATÉRIALISÉES pour requêtes lourdes
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_country_year_stats AS
SELECT 
  c.countrycode,
  c.shortname,
  i.year,
  COUNT(i.indicatorcode) as indicator_count,
  COUNT(DISTINCT i.indicatorcode) as unique_indicators,
  COUNT(CASE WHEN i.value IS NOT NULL THEN 1 END) as non_null_values
FROM country c
LEFT JOIN indicators i ON c.countrycode = i.countrycode
GROUP BY c.countrycode, c.shortname, i.year;

CREATE UNIQUE INDEX IF NOT EXISTS mv_country_year_stats_idx 
  ON mv_country_year_stats(countrycode, year);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_latest_indicators AS
SELECT DISTINCT ON (countrycode, indicatorcode)
  countrycode,
  indicatorcode,
  year,
  value
FROM indicators
WHERE value IS NOT NULL
ORDER BY countrycode, indicatorcode, year DESC;

CREATE UNIQUE INDEX IF NOT EXISTS mv_latest_indicators_idx 
  ON mv_latest_indicators(countrycode, indicatorcode);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_indicator_stats AS
SELECT 
  indicatorcode,
  year,
  COUNT(*) as country_count,
  AVG(value) as avg_value,
  MIN(value) as min_value,
  MAX(value) as max_value,
  STDDEV_POP(value) as stddev_value
FROM indicators
WHERE value IS NOT NULL
GROUP BY indicatorcode, year;

CREATE UNIQUE INDEX IF NOT EXISTS mv_indicator_stats_idx 
  ON mv_indicator_stats(indicatorcode, year);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_region_indicators AS
SELECT 
  c.region,
  i.year,
  i.indicatorcode,
  AVG(i.value) as avg_value,
  COUNT(*) as country_count
FROM country c
LEFT JOIN indicators i ON c.countrycode = i.countrycode
WHERE i.value IS NOT NULL
GROUP BY c.region, i.year, i.indicatorcode;

CREATE INDEX IF NOT EXISTS mv_region_indicators_idx 
  ON mv_region_indicators(region, indicatorcode, year);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_yearly_trends AS
SELECT 
  indicatorcode,
  countrycode,
  year,
  value,
  LAG(value) OVER (PARTITION BY countrycode, indicatorcode ORDER BY year) as prev_value,
  (value - LAG(value) OVER (PARTITION BY countrycode, indicatorcode ORDER BY year)) as year_change
FROM indicators
WHERE value IS NOT NULL;

CREATE INDEX IF NOT EXISTS mv_yearly_trends_idx 
  ON mv_yearly_trends(countrycode, indicatorcode, year);
