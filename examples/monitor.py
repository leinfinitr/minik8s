from prometheus_client import start_http_server, Gauge, Counter, Histogram, Summary
import random
import time

# 创建自定义监控指标
REQUEST_TIME = Summary('request_processing_seconds', 'Time spent processing request')
REQUEST_COUNT = Counter('request_count', 'Total request count')
TEMPERATURE_GAUGE = Gauge('temperature', 'Temperature in Celsius')
REQUEST_LATENCY = Histogram('request_latency_seconds', 'Request latency in seconds')

def process_request():
    """模拟处理请求并记录各种指标"""
    start_time = time.time()
    
    # 增加请求计数
    REQUEST_COUNT.inc()

    # 模拟处理时间和温度变化
    process_time = random.uniform(0.1, 2.0)
    time.sleep(process_time)
    temperature = random.uniform(20.0, 30.0)
    TEMPERATURE_GAUGE.set(temperature)

    # 记录处理时间
    REQUEST_TIME.observe(process_time)
    REQUEST_LATENCY.observe(time.time() - start_time)

if __name__ == '__main__':
    # 启动一个 HTTP 服务器，监听 9290 端口
    start_http_server(9290)
    print("Prometheus metrics server started on port 9290")

    # 不断处理请求，更新指标
    while True:
        process_request()
