## 1. Liste paginée + recherche + filtres
```graphql
type Query {
  # Liste paginée, avec recherche et filtres simples
  countries(
    search: String
    region: String
    incomeGroup: String
    limit: Int = 10
    offset: Int = 0
  ): [Country!]!
}
```

Back :  
- `search` → `COALESCE(ShortName, ''), COALESCE(LongName, ''), COALESCE(TableName, '') LIKE %search%`
- `region` → `WHERE Region = ?`
- `incomeGroup` → `WHERE IncomeGroup = ?`
- `limit/offset` → `.Limit(limit).Offset(offset).Order("ShortName ASC")`

***

## 2. Détail d’un pays (par code WDI ou Alpha2)

```graphql
type Query {
  country(code: String!): Country
  countryByAlpha2(code2: String!): Country
}
```

Back :
- `country` → `WHERE CountryCode = ? LIMIT 1`
- `countryByAlpha2` → `WHERE Alpha2Code = ? LIMIT 1`

***

## 3. Liste des régions disponibles

```graphql
type Query {
  regions: [String!]!
}
```

Back :
- `SELECT DISTINCT Region FROM Country WHERE Region IS NOT NULL ORDER BY Region`

***

## 4. Liste des groupes de revenu disponibles

```graphql
type Query {
  incomeGroups: [String!]!
}
```

Back :
- `SELECT DISTINCT IncomeGroup FROM Country WHERE IncomeGroup IS NOT NULL ORDER BY IncomeGroup`

***

## 5. Compteurs / Statistiques simples

### a) Nombre total de pays

```graphql
type Query {
  countryCount: Int!
}
```

Back :
- `SELECT COUNT(*) FROM Country`

### b) Nombre de pays par région

```graphql
type RegionCount {
  region: String!
  count: Int!
}

type Query {
  regionCounts: [RegionCount!]!
}
```

Back (SQL brut ou GORM.Raw) :
- `SELECT Region, COUNT(*) AS Count FROM Country WHERE Region IS NOT NULL GROUP BY Region`

***

## 6. Liste de pays par liste de codes

Utile plus tard pour la comparaison multi-pays côté frontend.

```graphql
type Query {
  countriesByCodes(codes: [String!]!): [Country!]!
}
```

Back :
- `WHERE CountryCode IN (?)`

***

## 7. Autocomplete / suggestion (pour les champs de recherche)

```graphql
type Query {
  autocompleteCountries(prefix: String!, limit: Int = 10): [Country!]!
}
```

Back :
- `WHERE ShortName LIKE prefix% OR LongName LIKE prefix% OR TableName LIKE prefix% LIMIT ?`

***

En résumé, pour `Country` tu peux donc viser ces queries côté schéma :

```graphql
type RegionCount {
  region: String!
  count: Int!
}

type Query {
  hello: String!

  countries(
    search: String
    region: String
    incomeGroup: String
    limit: Int = 10
    offset: Int = 0
  ): [Country!]!

  country(code: String!): Country
  countryByAlpha2(code2: String!): Country

  regions: [String!]!
  incomeGroups: [String!]!
  countryCount: Int!
  regionCounts: [RegionCount!]!

  countriesByCodes(codes: [String!]!): [Country!]!
  autocompleteCountries(prefix: String!, limit: Int = 10): [Country!]!
}
```

On a déjà `countries` et une partie de la recherche/pagination opérationnelle.  
Tu peux maintenant décider quelle query ajouter en premier (par ex. `country`, `regionCounts`, `regions`) et on écrit ensemble `usecase + repository + resolver` pour celle que tu choisis.
