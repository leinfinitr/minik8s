import http.server
import os

class CustomHTTPHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write('This is the Pod IP: '.encode() + os.environ['POD_IP'].encode())

if __name__ == '__main__':
    server_address = ('', 7080)
    httpd = http.server.HTTPServer(server_address, CustomHTTPHandler)
    httpd.serve_forever()