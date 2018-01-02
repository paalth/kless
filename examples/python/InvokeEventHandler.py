#!/usr/bin/env python
"""
Simple invocation of event handler

Usage::
    ./InvokeEventHandler.py [<port>]

"""
from http.server import BaseHTTPRequestHandler, HTTPServer
from KlessContext import Context
from KlessRequest import Request
from KlessResponse import Response
from EventHandler1 import *

class S(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_HEAD(self):
        self._set_headers()
        
    def do_GET(self):
        self._set_headers()
        context = Context()
        request = Request()
        response = Response()
        handler = EventHandler1()
        handler.eventHandler(context, request, response)
        self.wfile.write("<html><body>OK</body></html>".encode("utf-8"))

    def do_POST(self):
        self._set_headers()
        context = Context()
        request = Request()
        response = Response()
        handler = EventHandler1()
        handler.eventHandler(context, request, response)
        content_length = int(self.headers['Content-Length']) 
        post_data = self.rfile.read(content_length) 
        self.wfile.write("<html><body>OK</body></html>".encode("utf-8"))
        
def run(server_class=HTTPServer, handler_class=S, port=8080):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print('Starting httpd...')
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()