# 实时公交
# Document
## 附近线路
### PATH
GET /api/v3/station-buses/bylocation/:location

### RESPONSE
```json
[{
  "sn": "公交站",
  "lc": 5,
  "slc": 1,
  "l":"23.001,113.002",
  "lines": [{
      "no": "25",
      "dir": "0",
      "adir": "1",
      "ssn": "起始站",
      "esn":"终点站",
      "nsn":"下一站",
      "o": 12,
      "buses":[{
        "o": 10,
        "s": 0.5,
        "l": "23.001,113.002",
        "d": 34
      },{
        "o": 10,
        "s": 1,
        "l": "23.001,113.002",
        "d": 34
      },{
        "o": 11,
        "s": 0.5,
        "l": "23.001,113.002",
        "d": 34
      }]
  }]
}]
```

## 公交站线路
### PATH
GET /api/v3/station-buses/byline/:city/:sn

### RESPONSE
同`附近线路`返回结果


## 线路信息
## PATH
GET /api/v3/lines/:city/:line/:dir


### 响应结果
```json
{
  "no": "25",
  "id": "xxxxxx",
  "adir": "1",
  "dir": "0",
  "ssn": "起始站",
  "esn": "终点站",
  "price": "1元",
  "ftime": "06:20",
  "ltime": "22:30",
  "ss":[{
    "o": 1,
    "sn":"起始站",
    "l": "36.20105,120.5255"
  },{
    "o": 2,
    "sn":"站点1",
    "l": "36.20105,120.5255"
  }],
  "buses": [{
    "o": 10,
    "s": 0.5,
    "l": "23.001,113.002",
    "d": 34
 }]
}
```

## 线路搜索
## PATH
GET /api/v3/search/:city/:keyword


### 响应结果
```json
[{
  "no": "25",
  "id": "xxxxxx",
  "adir": "1",
  "dir": "0",
  "ssn": "起始站",
  "esn": "终点站",
  "price": "1元",
  "ftime": "06:20",
  "ltime": "22:30",
  "ss":[{
    "o": 1,
    "sn":"起始站",
    "l": "36.20105,120.5255"
  },{
    "o": 2,
    "sn":"站点1",
    "l": "36.20105,120.5255"
  }],
  "buses": [{
    "o": 10,
    "s": 0.5,
    "l": "23.001,113.002",
    "d": 34
 }]
}]
```


## 线路运行公交

### PATH
GET /api/v3/lines/:city/:line/:dir/buses

### Parameter
| 字段 | 类型 | 说明 |
| --- | --- | --- |
|sn|string|公交站|


### RESPONSE
```json
[{
  "o": 10,
  "s": 0.5,
  "l": "23.001,113.002",
  "d": 34
},{
  "o": 10,
  "s": 0.5,
  "l": "23.001,113.002",
  "d": 34
}]
```
