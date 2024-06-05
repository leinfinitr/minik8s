import http.server
import os
import subprocess

class CustomHTTPHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write('This is the ifconfig output: '.encode() + os.environ['IFCONFIG_RESULT'].encode())

if __name__ == '__main__':
    server_address = ('', 7080)
    httpd = http.server.HTTPServer(server_address, CustomHTTPHandler)
    httpd.serve_forever()