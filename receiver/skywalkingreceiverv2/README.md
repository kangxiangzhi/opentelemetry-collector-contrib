# skywalkingv2 Receiver

Receives trace data in [Skywalking](https://skywalking.apache.org/) format.

Supported pipeline types: traces

## ⚠️ Warning

Note: This component is currently work in progress, and traces receiver is not yet fully functional.

## Getting Started

By default, the Skywalking receiver will not serve any protocol. A protocol must be
named under the `protocols` object for the Skywalking receiver to start. The
below protocols are supported, each supports an optional `endpoint`
object configuration parameter.

- `grpc` (default `endpoint` = 0.0.0.0:11800)

Examples:

```yaml
receivers:
  skywalkingv2:
    protocols:
      grpc:
    components:
      - name: Unknown
        id: 0
      - name: Tomcat
        id: 1
      - name: HttpClient
        id: 2
      - name: Dubbo
        id: 3
      - name: H2
        id: 4
      - name: Mysql
        id: 5
      - name: ORACLE
        id: 6
      - name: Redis
        id: 7
      - name: Motan
        id: 8
      - name: MongoDB
        id: 9
      - name: Resin
        id: 10
      - name: Feign
        id: 11
      - name: OKHttp
        id: 12
      - name: SpringRestTemplate
        id: 13
      - name: SpringMVC
        id: 14
      - name: Struts2
        id: 15
      - name: NutzMVC
        id: 16
      - name: NutzHttp
        id: 17
      - name: JettyClient
        id: 18
      - name: JettyServer
        id: 19
      - name: Memcached
        id: 20
      - name: ShardingJDBC
        id: 21
      - name: PostgreSQL
        id: 22
      - name: GRPC
        id: 23
      - name: ElasticJob
        id: 24
      - name: RocketMQ
        id: 25
      - name: httpasyncclient
        id: 26
      - name: Kafka
        id: 27
      - name: ServiceComb
        id: 28
      - name: Hystrix
        id: 29
      - name: Jedis
        id: 30
      - name: SQLite
        id: 31
      - name: H2
        id: 32
      - name: mysql
        id: 33
      - name: ojdbc
        id: 34
      - name: Memcached
        id: 35
      - name: Memcached
        id: 36
      - name: PostgreSQL
        id: 37
      - name: RocketMQ
        id: 38
      - name: RocketMQ
        id: 39
      - name: kafka
        id: 40
      - name: kafka
        id: 41
      - name: MongoDB
        id: 42
      - name: SOFARPC
        id: 43
      - name: ActiveMQ
        id: 44
      - name: ActiveMQ
        id: 45
      - name: ActiveMQ
        id: 46
      - name: Elasticsearch
        id: 47
      - name: Elasticsearch
        id: 48
      - name: http
        id: 49
      - name: rpc
        id: 50
      - name: RabbitMQ
        id: 51
      - name: RabbitMQ
        id: 52
      - name: RabbitMQ
        id: 53
      - name: Canal
        id: 54
      - name: Gson
        id: 55
      - name: Redis
        id: 56
      - name: Redis
        id: 57
      - name: Zookeeper
        id: 58
      - name: Vertx
        id: 59
      - name: ShardingSphere
        id: 60
      - name: spring-cloud-gateway
        id: 61
      - name: RESTEasy
        id: 62
      - name: Solr
        id: 63
      - name: Solr
        id: 64
      - name: SpringAsync
        id: 65
      - name: JdkHttp
        id: 66
      - name: spring-webflux
        id: 67
      - name: Play
        id: 68
      - name: Cassandra
        id: 69
      - name: Cassandra
        id: 70
      - name: Light4J
        id: 71
      - name: Pulsar
        id: 72
      - name: Pulsar
        id: 73
      - name: Pulsar
        id: 74
      - name: Ehcache
        id: 75
      - name: SocketIO
        id: 76
      - name: Elasticsearch
        id: 77
      - name: spring-tx
        id: 78
      - name: Armeria
        id: 79
      - name: JdkThreading
        id: 80
      - name: KotlinCoroutine
        id: 81
      - name: AvroServer
        id: 82
      - name: AvroClient
        id: 83
      - name: Undertow
        id: 84
      - name: Finagle
        id: 85
      - name: Mariadb
        id: 86
      - name: Mariadb
        id: 87
      - name: quasar
        id: 88
      - name: InfluxDB
        id: 89
      - name: InfluxDB
        id: 90
      - name: brpc-java
        id: 91
      - name: GraphQL
        id: 92
      - name: spring-annotation
        id: 93
      - name: HBase
        id: 94
      - name: kafka-consumer
        id: 95
      - name: SpringScheduled
        id: 96
      - name: quartz-scheduler
        id: 97
      - name: xxl-job
        id: 98
      - name: spring-webflux-webclient
        id: 99
      - name: thrift-server
        id: 100
      - name: thrift-client
        id: 101
      - name: AsyncHttpClient
        id: 102
      - name: dbcp
        id: 103
      - name: SqlServer
        id: 104
      - name: Apache-CXF
        id: 105
      - name: dolphinscheduler
        id: 106
      - name: JsonRpc
        id: 107
      - name: seata
        id: 108
      - name: MyBatis
        id: 109
      - name: tcp
        id: 110
      - name: AzureHttpTrigger
        id: 111
      - name: Neo4j
        id: 112
      - name: Sentinel
        id: 113
      - name: GuavaCache
        id: 114
      - name: AspNetCore
        id: 3001
      - name: EntityFrameworkCore
        id: 3002
      - name: SqlClient
        id: 3003
      - name: CAP
        id: 3004
      - name: StackExchange.Redis
        id: 3005
      - name: SqlServer
        id: 3006
      - name: Npgsql
        id: 3007
      - name: MySqlConnector
        id: 3008
      - name: EntityFrameworkCore.InMemory
        id: 3009
      - name: EntityFrameworkCore.SqlServer
        id: 3010
      - name: EntityFrameworkCore.Sqlite
        id: 3011
      - name: Pomelo.EntityFrameworkCore.MySql
        id: 3012
      - name: Npgsql.EntityFrameworkCore.PostgreSQL
        id: 3013
      - name: InMemoryDatabase
        id: 3014
      - name: AspNet
        id: 3015
      - name: SmartSql
        id: 3016
      - name: HttpServer
        id: 4001
      - name: Express
        id: 4002
      - name: Egg
        id: 4003
      - name: Koa
        id: 4004
      - name: Axios
        id: 4005
      - name: Mongoose
        id: 4006
      - name: ServiceCombMesher
        id: 5001
      - name: ServiceCombServiceCenter
        id: 5002
      - name: MOSN
        id: 5003
      - name: GoHttpServer
        id: 5004
      - name: GoHttpClient
        id: 5005
      - name: Gin
        id: 5006
      - name: Gear
        id: 5007
      - name: GoMicroClient
        id: 5008
      - name: GoMicroServer
        id: 5009
      - name: GoKratosServer
        id: 5010
      - name: GoKratosClient
        id: 5011
      - name: mysql
        id: 5012
      - name: Nginx
        id: 6000
      - name: Kong
        id: 6001
      - name: APISIX
        id: 6002
      - name: Python
        id: 7000
      - name: Flask
        id: 7001
      - name: Requests
        id: 7002
      - name: mysql
        id: 7003
      - name: Django
        id: 7004
      - name: Tornado
        id: 7005
      - name: Urllib3
        id: 7006
      - name: Sanic
        id: 7007
      - name: AioHttp
        id: 7008
      - name: Pyramid
        id: 7009
      - name: Psychopg
        id: 7010
      - name: Celery
        id: 7011
      - name: PHP
        id: 8001
      - name: cURL
        id: 8002
      - name: mysql
        id: 8003
      - name: mysql
        id: 8004
      - name: Yar
        id: 8005
      - name: redis
        id: 8006
      - name: EnvoyProxy
        id: 9000
      - name: JavaScript
        id: 10000
      - name: ajax
        id: 10001

service:
  pipelines:
    traces:   #链路
      receivers: [skywalking]
      
    metrics:  #JVM指标
      receivers: [skywalking]
```
