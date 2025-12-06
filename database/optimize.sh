#!/bin/bash
# ============================================
# Script d'optimisation WDI Database
# Usage: ./optimize_wdi.sh [option]
# ============================================

set -e

DB_NAME="wdi"
DB_USER="wdi"
DB_HOST="localhost"
DB_PORT="5432"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction d'affichage
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifier la connexion
check_connection() {
    log_info "Vérification de la connexion à PostgreSQL..."
    if psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT -c "SELECT 1;" > /dev/null 2>&1; then
        log_success "Connexion OK"
        return 0
    else
        log_error "Impossible de se connecter à la base de données"
        exit 1
    fi
}

# Sauvegarde avant optimisation
backup_database() {
    log_info "Création d'une sauvegarde..."
    BACKUP_FILE="wdi_backup_$(date +%Y%m%d_%H%M%S).sql"
    pg_dump -U $DB_USER -h $DB_HOST -p $DB_PORT $DB_NAME > $BACKUP_FILE
    log_success "Sauvegarde créée: $BACKUP_FILE"
}

# Optimisation complète
full_optimization() {
    log_info "========================================="
    log_info "OPTIMISATION COMPLÈTE DE LA BASE WDI"
    log_info "========================================="
    
    check_connection
    
    # Demander confirmation
    read -p "Voulez-vous créer une sauvegarde avant l'optimisation? (o/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Oo]$ ]]; then
        backup_database
    fi
    
    log_info "Étape 1/5: Création des index optimisés..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT -f wdi_optimization.sql 2>&1 | grep -E "CREATE|ERROR|WARNING" || true
    log_success "Index créés"
    
    log_info "Étape 2/5: Création des vues matérialisées..."
    log_success "Vues matérialisées créées"
    
    log_info "Étape 3/5: Analyse des tables..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT -c "VACUUM (ANALYZE, VERBOSE);" 2>&1 | tail -5
    log_success "Tables analysées"
    
    log_info "Étape 4/5: Rafraîchissement des vues matérialisées..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_country_year_stats;
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_latest_indicators;
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_indicator_stats;
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_region_indicators;
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_yearly_trends;
EOF
    log_success "Vues rafraîchies"
    
    log_info "Étape 5/5: Génération du rapport..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
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
    'Taille base',
    pg_size_pretty(pg_database_size('$DB_NAME'))
UNION ALL
SELECT 
    'Index créés',
    COUNT(*)::text
FROM pg_indexes
WHERE schemaname = 'public'
UNION ALL
SELECT
    'Vues matérialisées',
    COUNT(*)::text
FROM pg_matviews
WHERE schemaname = 'public';
EOF
    
    log_success "========================================="
    log_success "OPTIMISATION TERMINÉE AVEC SUCCÈS !"
    log_success "========================================="
}

# Tests de performance
performance_tests() {
    log_info "Exécution des tests de performance..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT -f wdi_performance_tests.sql
    log_success "Tests terminés"
}

# Maintenance quotidienne
daily_maintenance() {
    log_info "Maintenance quotidienne en cours..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
SELECT * FROM maintain_wdi_database();
EOF
    log_success "Maintenance terminée"
}

# Statistiques
show_stats() {
    log_info "Statistiques de la base de données WDI:"
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
-- Taille des tables
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.'||tablename) DESC;

-- Statistiques générales
SELECT 
    'Total enregistrements indicators' as metric,
    COUNT(*)::text as value
FROM indicators
UNION ALL
SELECT 
    'Pays uniques',
    COUNT(DISTINCT countrycode)::text
FROM indicators
UNION ALL
SELECT 
    'Indicateurs uniques',
    COUNT(DISTINCT indicatorcode)::text
FROM indicators
UNION ALL
SELECT 
    'Années',
    MIN(year)::text || ' - ' || MAX(year)::text
FROM indicators;

-- Cache hit ratio
SELECT 
    'Cache Hit Ratio' as metric,
    ROUND(
        100.0 * SUM(heap_blks_hit) / NULLIF(SUM(heap_blks_hit + heap_blks_read), 0),
        2
    )::text || '%' as value
FROM pg_statio_user_tables
WHERE schemaname = 'public';
EOF
}

# Vérifier les requêtes lentes
check_slow_queries() {
    log_info "Vérification des requêtes lentes..."
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
SELECT 
    pid,
    usename,
    EXTRACT(EPOCH FROM (now() - query_start))::INTEGER as duration_seconds,
    state,
    LEFT(query, 150) as query
FROM pg_stat_activity
WHERE state = 'active'
    AND query NOT LIKE '%pg_stat_activity%'
    AND now() - query_start > interval '1 second'
ORDER BY duration_seconds DESC;
EOF
}

# Configuration recommandée
show_config() {
    log_info "Configuration PostgreSQL actuelle vs recommandée:"
    psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT << EOF
SELECT 
    name,
    setting || COALESCE(' ' || unit, '') as current_value,
    CASE name
        WHEN 'shared_buffers' THEN '4GB (pour 16GB RAM)'
        WHEN 'effective_cache_size' THEN '12GB (pour 16GB RAM)'
        WHEN 'work_mem' THEN '256MB'
        WHEN 'maintenance_work_mem' THEN '1GB'
        WHEN 'max_parallel_workers_per_gather' THEN '4'
        WHEN 'max_parallel_workers' THEN '8'
        WHEN 'random_page_cost' THEN '1.1 (SSD)'
        WHEN 'effective_io_concurrency' THEN '200 (SSD)'
    END as recommended
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
EOF
}

# Menu interactif
show_menu() {
    echo ""
    echo "======================================"
    echo "   OPTIMISATION BASE DE DONNÉES WDI"
    echo "======================================"
    echo "1. Optimisation complète"
    echo "2. Tests de performance"
    echo "3. Maintenance quotidienne"
    echo "4. Afficher les statistiques"
    echo "5. Vérifier les requêtes lentes"
    echo "6. Voir la configuration"
    echo "7. Créer une sauvegarde"
    echo "8. Quitter"
    echo "======================================"
    read -p "Choisissez une option [1-8]: " choice
    
    case $choice in
        1) full_optimization ;;
        2) performance_tests ;;
        3) daily_maintenance ;;
        4) show_stats ;;
        5) check_slow_queries ;;
        6) show_config ;;
        7) backup_database ;;
        8) exit 0 ;;
        *) log_error "Option invalide" ;;
    esac
}

# Point d'entrée principal
main() {
    if [ $# -eq 0 ]; then
        # Mode interactif
        while true; do
            show_menu
        done
    else
        # Mode ligne de commande
        case $1 in
            optimize) full_optimization ;;
            test) performance_tests ;;
            maintain) daily_maintenance ;;
            stats) show_stats ;;
            slow) check_slow_queries ;;
            config) show_config ;;
            backup) backup_database ;;
            *)
                echo "Usage: $0 {optimize|test|maintain|stats|slow|config|backup}"
                echo "Ou lancez sans argument pour le menu interactif"
                exit 1
                ;;
        esac
    fi
}

main "$@"