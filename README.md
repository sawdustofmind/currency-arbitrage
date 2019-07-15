# currency-arbitrage
an app that analyzes currency rates of certain exchange to find arbitrage opportunities

To solve the problem we request tickers of some exchange (exmo), then we construct the graph. Each vertex is currency of ticker. Each ticker represents two edges with buy price (from base to quote) and 1 devided by sell price (from quote to base). 
In other words an edge weight is coefficient we have to multiply price to migrage to another currency in case of market order placing (for small volume)

Next to deal with commission we multiply each edge weight by 1 minus commission and transform each weight to negative natural logarifm of price to use all there edges in any APSP algoritm.

Currently implemented only Floyd-Warshall algorithm that immediately stops when negative cycle occurred

## Demo
You see demo [here](http://193.187.174.47/)
![image](https://user-images.githubusercontent.com/29863444/61198947-bae78480-a6e4-11e9-957d-c366d59aad03.png)

## Deploy

### With docker

Here is a default docker multistage build
`docker build . [-t IMAGE_NAME]`

`docker run (HASH|IMAGE_NAME] -p PORT:12346`
    
### Local

`go build && ./currency_arbitrage PORT`

