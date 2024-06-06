import http.server
import subprocess

class CustomHTTPHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        # 执行计算密集型任务
        result_of_computation = self.intensive_computation()

        # 执行ifconfig命令并获取结果
        ifconfig_result = subprocess.check_output(["ifconfig"], universal_newlines=True)

        # 发送HTTP响应
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()

        # 发送计算结果和ifconfig结果
        response = f"This is the computation result: {result_of_computation}\n\nThis is the ifconfig output:\n{ifconfig_result}"
        self.wfile.write(response.encode())

    def intensive_computation(self):
        # 这里加入计算密集型任务的代码
        # 例如，计算斐波那契数列的第n项
        def fib(n):
            if n <= 1:
                return n
            else:
                return fib(n-1) + fib(n-2)
        
        # 假设我们计算斐波那契数列的第30项
        result = fib(30)
        return result

if __name__ == '__main__':
    server_address = ('', 7080)
    httpd = http.server.HTTPServer(server_address, CustomHTTPHandler)
    httpd.serve_forever()