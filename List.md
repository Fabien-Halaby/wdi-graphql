# Fonctionnalités possibles avec le dataset WDI

## 1. Page d'accueil (Dashboard global)

- Liste/interactif de **tous les pays** avec recherche et filtres.
- Carte ou top 10 des pays selon un indicateur (ex : PIB, population).
- Statistiques globales : nombre total de pays, d’indicateurs, répartition par région/groupe de revenu.

---

## 2. Exploration de pays

- **Voir le profil détaillé d’un pays** :
  - Infos de base (nom, région, groupe de revenu, devise, etc.).
  - Liste des indicateurs disponibles pour ce pays (avec années couvertes).
  - Graphique de série temporelle pour un indicateur sélectionné (ex : évolution PIB, population, etc.).
  - Tableau de toutes les valeurs pour un indicateur.
- **Filtrer et rechercher un pays** par nom, région, groupe de revenu, etc.
- **Accéder directement à un pays via son code** (URL paramètre).

---

## 3. Exploration des indicateurs (Series)

- **Parcourir tous les indicateurs** avec recherche et filtres (par thème : démographie, économie, environnement...).
- Accéder à la définition complète, l’unité, la source et la méthodologie d’un indicateur.
- **Filtrer les indicateurs** disponibles en fonction de critères (thème, unité, etc.).

---

## 4. Visualisation de séries temporelles

- **Tracer le graphique** d’un indicateur pour un pays donné (ex : population France 1960-2023).
- **Télécharger les données** en CSV/JSON pour un indicateur/pays donné.
- Naviguer dans les années avec une plage personnalisée.

---

## 5. Comparaison multi-pays

- **Comparer plusieurs pays** sur le même indicateur (curve multiligne : ex Population ou PIB de 5 pays de 1990 à 2023).
- **Comparer un même pays** sur plusieurs indicateurs (ex : croissance vs inflation).
- Tableau comparatif pour voir les valeurs côte à côte.

---

## 6. Statistiques avancées et agrégées (Analytics)

- **Analyse régionale** :
  - Moyenne d’un indicateur pour chaque région/économie.
  - Répartition par groupe de revenu.
  - Top N pays sur un indicateur particulier (best/worst).
- Statistiques et agrégations temporelles (évolution moyenne, CAGR, etc.)

---

## 7. Recherches personnalisées et filtres dynamiques

- **Filtrer toutes les données** par région, groupe, période, valeur minimale/maximale.
- Requete full-text sur les noms de pays et d'indicateurs.
- Voir uniquement des pays "favoris" ou "sélectionnés".

---

## 8. Affichage et explication des notes/contextes  

- **Afficher les footnotes et country notes** pour un indicateur/année/pays (contextualisation des chiffres).
- Visualisation de la méthodologie/exception sur chaque série.

---

## 9. Export et partage

- **Télécharger/exporter** les tables et graphes (CSV, JSON, PNG pour le graphique).
- Copier un lien vers la visualisation courante pour partager avec d’autres utilisateurs.

---

## 10. Bonus UX/UI

- Sauvegarder des “favoris” ou des requêtes personnalisées.
- Passer en “dark mode” ou choisir le thème de visualisation.
- Interface responsive mobile/desktop.
- Système de pagination et navigation fluide dans la liste des pays/indicateurs.

---

## 11. Cas d'usage administrateur / enseignant (optionnel)

- Proposer des “quiz” ou des challenges (ex : identifie le pays ayant eu la plus forte croissance).
- Ajouter ou annoter des footnotes/mises à jour manuelles (si droits admin).

---

## 12. Explication des potentiels graphiques

- Indicateur temps (courbe/line chart)
- Classement (bar chart vertical top N)
- Carte (if wanted : color heatmap par région)
- Pie chart (par groupe de revenu)
- Tableau téléchargeable (vue tabulaire paginée)

---
