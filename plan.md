## 1. Base pays / indicateurs (métadonnées)
- [ ] **Query `countries`**  
  - Visualisation : table simple / select de pays  
  - But : alimenter filtres de pays  
  - Data retournée :  
    - `countrycode`, `shortname`, `region`, `incomegroup`  
  - Variables :  
    - `search: String`, `region: String`, `incomegroup: String`, `limit: Int`, `offset: Int`  

- [ ] **Query `indicatorsMeta`**  
  - Visualisation : select d’indicateurs, recherche  
  - But : choisir facilement les séries WDI  
  - Data :  
    - `indicatorcode`, `indicatorname`, `topic`, `unitofmeasure`  
  - Variables :  
    - `search: String`, `topic: String`, `limit: Int`, `offset: Int`  

***

## 2. Séries temporelles (line / area charts)

- [ ] **Query `indicatorTimeSeries`**  
  - Visualisation : **line chart** (1 pays, 1 indicateur)  
  - But : évolution d’un indicateur pour un pays (vue “GDP per capita – France 1990–2020”)  
  - Data :  
    - `year: Int`, `value: Float`  
  - Variables :  
    - `countrycode: String!`  
    - `indicatorcode: String!`  
    - `startyear: Int!`  
    - `endyear: Int!`  

- [ ] **Query `compareIndicatorTwoCountries`**  
  - Visualisation : **line chart 2 séries** (comparaison)  
  - But : comparer deux pays sur un indicateur (déjà commencé avec `compareIndicator`)  
  - Data (par point) :  
    - `year`, `indicatorcode`, `countrycode1`, `value1`, `countrycode2`, `value2`  
  - Variables :  
    - `indicatorcode: String!`  
    - `countrycode1: String!`  
    - `countrycode2: String!`  
    - `startyear: Int!`  
    - `endyear: Int!`  

- [ ] **Query `indicatorRegionTimeSeries`**  
  - Visualisation : **multi-line chart** (régions ou income groups)  
  - But : suivre la moyenne régionale / par groupe de revenu dans le temps  
  - Data (par point) :  
    - `year`, `group: String` (region ou incomegroup), `avgvalue`  
  - Variables :  
    - `indicatorcode: String!`  
    - `groupby: IndicatorGroupBy!` (enum `REGION | INCOME_GROUP`)  
    - `startyear: Int!`  
    - `endyear: Int!`  

***

## 3. Cartes (Leaflet, choroplèthe / bubbles)

- [ ] **Query `worldMapIndicator`**  
  - Visualisation : **choroplèthe** Leaflet (ou cercles)  
  - But : valeur d’un indicateur pour tous les pays sur une année (FDI map, CO₂ map, etc.)  
  - Data (par pays) :  
    - `countrycode`, `shortname`, `region`, `value`  
  - Variables :  
    - `indicatorcode: String!`  
    - `year: Int!`  

- [ ] **Query `countryIndicatorSnapshot`**  
  - Visualisation : popup de carte / fiche pays courte  
  - But : afficher dans le popup de la carte plusieurs indicateurs clés pour un pays/année  
  - Data :  
    - `countrycode`, `year`, `[ { indicatorcode, indicatorname, value } ]`  
  - Variables :  
    - `countrycode: String!`  
    - `year: Int!`  
    - `indicatorcodes: [String!]!`  

***

## 4. Rankings & distributions (bar / column / histogram)

- [ ] **Query `topCountriesByIndicator`**  
  - Visualisation : **bar chart horizontal** (Top 10 / Bottom 10)  
  - But : afficher les pays avec les valeurs max/min pour un indicateur donné  
  - Data :  
    - `countrycode`, `shortname`, `value`, `rank`  
  - Variables :  
    - `indicatorcode: String!`  
    - `year: Int!`  
    - `limit: Int!` (ex. 10)  
    - `direction: SortDirection!` (enum `DESC` pour top, `ASC` pour bottom)  

- [ ] **Query `indicatorDistribution`**  
  - Visualisation : **histogram/bar chart** (distribution par classes)  
  - But : distribution d’un indicateur entre pays pour une année (0–2k, 2k–5k, etc.)  
  - Data (par classe) :  
    - `bucketLabel: String`, `min`, `max`, `count`  
  - Variables :  
    - `indicatorcode: String!`  
    - `year: Int!`  
    - `buckets: Int!`  

***

## 5. Profils pays & KPI (cards, radar, donut)

- [ ] **Query `countryProfile`**  
  - Visualisation :  
    - cards KPI + mini sparklines, éventuellement **radar chart**  
  - But : page “Country detail” avec plusieurs indicateurs clés  
  - Data :  
    - `country` : `countrycode`, `shortname`, `region`, `incomegroup`  
    - `kpis`: `[ {indicatorcode, indicatorname, latestYear, latestValue, series: [ {year, value} ] } ]`  
  - Variables :  
    - `countrycode: String!`  
    - `indicatorkeys: [String!]!` (ex. `["NY.GDP.PCAP.CD", "SP.DYN.LE00.IN", "EN.ATM.CO2E.PC"]`)  

- [ ] **Query `regionElectricityAccess`**  
  - Visualisation : **donut chart** (access / no access / not electrified)  
  - But : part de population avec accès à l’électricité pour une région donnée  
  - Data :  
    - `region`, `year`, `accessPercent`, `noAccessPercent`, `notElectrifiedPercent`  
  - Variables :  
    - `region: String!`  
    - `year: Int!`  

***

## 6. Tables d’exploration (data grid)

- [ ] **Query `indicatorTable`**  
  - Visualisation : **table paginée + triable**  
  - But : explorer les données brutes pour un indicateur / une plage d’années  
  - Data (par ligne) :  
    - `countrycode`, `shortname`, `year`, `value`  
  - Variables :  
    - `indicatorcode: String!`  
    - `startyear: Int!`  
    - `endyear: Int!`  
    - `limit: Int!`, `offset: Int!`  
    - `orderBy: IndicatorTableOrderBy!` (enum `YEAR`, `VALUE`, `COUNTRY`)  
    - `direction: SortDirection!`  

- [ ] **Query `latestIndicatorValues`**  
  - Visualisation : **table “Latest data – CO₂”**  
  - But : tableau des derniers chiffres disponibles par pays pour un indicateur  
  - Data :  
    - `countrycode`, `shortname`, `year`, `value`  
  - Variables :  
    - `indicatorcode: String!`  
    - `limit: Int!`, `offset: Int!`  

***

## 7. Checklist par visualisation (liens vers queries communes)

- [ ] **Line chart : évolution d’un indicateur pour un pays**  
  - Query utilisée : `indicatorTimeSeries`  
  - Variables typiques : `countrycode="FRA"`, `indicatorcode="NY.GDP.PCAP.CD"`, `startyear=1990`, `endyear=2020`  

- [ ] **Line chart comparaison 2 pays**  
  - Query utilisée : `compareIndicatorTwoCountries`  

- [ ] **Multi-line chart régions / revenus**  
  - Query : `indicatorRegionTimeSeries`  

- [ ] **Carte Leaflet choroplèthe**  
  - Query : `worldMapIndicator` + éventuellement `countryIndicatorSnapshot` pour les popups  

- [ ] **Bar chart Top 10 / Bottom 10**  
  - Query : `topCountriesByIndicator`  

- [ ] **Histogram / distribution**  
  - Query : `indicatorDistribution`  

- [ ] **Page Country Profile (KPI + mini-charts)**  
  - Query : `countryProfile`  

- [ ] **Donut électricité / autres parts**  
  - Query : `regionElectricityAccess`  

- [ ] **Table d’exploration indicateur**  
  - Query : `indicatorTable`  

- [ ] **Table “Latest data – CO₂”**  
  - Query : `latestIndicatorValues`  
