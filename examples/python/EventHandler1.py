import sys

class EventHandler1:
    def eventHandler(self, context, request, response):
        sys.stdout.write("Inside event handler example 1\n")
        sys.stdout.flush()