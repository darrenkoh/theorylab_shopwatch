create (p:Person {name: 'darren', deviceid: 'A0001'})-[:ASK {price: 200, eq: '<'}]->(n:Product {name: 'Jackery SolarSaga 60W Solar Panel for Explorer 160/240/500 as Portable Solar Generator, Portable Foldable Solar Charger for Summer Camping Van RV(Can\'t Charge Explorer 440/ PowerPro)'})-[:AT {url: 'https://www.amazon.com/Jackery-SolarSaga-Portable-Explorer-Foldable/dp/B07Q71LX84', price: 299.99}]->(m:Merchant {name: 'Amazon'})

match (p:Person {name: 'darren', deviceid: 'A0001'}), (n:Product {name: 'Jackery SolarSaga 60W Solar Panel for Explorer 160/240/500 as Portable Solar Generator, Portable Foldable Solar Charger for Summer Camping Van RV(Can\'t Charge Explorer 440/ PowerPro)'})
merge (m:Merchant {name: 'Walmart'})
merge (n)-[:AT {url: 'https://www.walmart.com/ip/Jackery-SolarSaga-60W-Solar-Panel-Explorer-160-240-500-Portable-Generator-Foldable-Charger-Summer-Camping-Van-RV-Can-t-Charge-440-PowerPro/939780115', price: 290}]->(m)

