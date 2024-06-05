import http.server
import os
import subprocess

class CustomHTTPHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        try:
            ifconfig_output = subprocess.check_output(['ifconfig'], text=True)
            self.wfile.write(f'This is the ifconfig output:\n{ifconfig_output}'.encode())
        except subprocess.CalledProcessError as e:
            self.wfile.write(f'Error executing ifconfig: {e}'.encode())

if __name__ == '__main__':
    server_address = ('', 7080)
    httpd = http.server.HTTPServer(server_address, CustomHTTPHandler)
    httpd.serve_forever()