# currency-arbitrage
an app that analyzes currency rates of certain exchange to find arbitrage opportunities

To solve the problem we request tickers of some exchange (exmo), then we construct the graph. Each vertex is currency of ticker. Each ticker represents two edges with buy price (from base to quote) and 1 devided by sell price (from quote to base). 
In other words an edge weight is coefficient we have to multiply price to migrage to another currency in case of market order placing (for small volume)

Next to deal with comission we multiply each edge weight by 1 minus comission and transform each weight to negative natural logarifm of price to use all there edges in any APSP algoritm.

Currently implemented only Floyd-Warshall algorithm that immediately stops when negative cycle occurred

## Deploy

### With docker

Here is a default docker multistage build
`docker build . [-t IMAGE_NAME]`

`docker run (HASH|IMAGE_NAME] -p PORT:12346`
    
### Local

`go build && ./currency_arbitrage PORT`

### Demo
You see demo [here](https://immense-dawn-98427.herokuapp.com/)